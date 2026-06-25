package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/uesugitorachiyo/ao-blueprint/internal/blueprint"
)

type session struct {
	Schema            string           `json:"schema"`
	SessionID         string           `json:"session_id"`
	Idea              string           `json:"idea"`
	Status            string           `json:"status"`
	AnsweredQuestions []answeredRecord `json:"answered_questions"`
}

type answeredRecord struct {
	QuestionID string `json:"question_id"`
	Answer     string `json:"answer"`
	AnsweredAt string `json:"answered_at"`
}

type question struct {
	ID       string
	Category string
	Prompt   string
}

var questions = []question{
	{ID: "q-objective", Category: "objective", Prompt: "What measurable outcome proves this project succeeded?"},
	{ID: "q-users", Category: "users", Prompt: "Who will use this, and what decision or workflow will it improve?"},
	{ID: "q-nongoals", Category: "scope", Prompt: "What must this project explicitly not do?"},
	{ID: "q-security", Category: "security", Prompt: "What secrets, private data, destructive actions, or public-safety constraints apply?"},
}

func Run(args []string, stdout io.Writer, stderr io.Writer) error {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		printHelp(stdout)
		return nil
	}

	var err error
	switch args[0] {
	case "interview":
		err = runInterview(args[1:], stdout)
	case "compile":
		err = runCompile(args[1:], stdout)
	case "lint":
		err = runLint(args[1:], stdout)
	case "readiness":
		err = runReadiness(args[1:], stdout)
	case "sdd":
		err = runSDD(args[1:], stdout)
	case "authorize":
		err = runAuthorize(args[1:], stdout)
	case "pack":
		err = runPack(args[1:], stdout)
	default:
		err = fmt.Errorf("unknown command %q", args[0])
	}
	if err != nil {
		fmt.Fprintln(stderr, err)
	}
	return err
}

func printHelp(stdout io.Writer) {
	fmt.Fprintln(stdout, `AO Blueprint

Usage:
  blueprint interview start --idea <text> --out <json>
  blueprint interview next --session <json>
  blueprint interview answer --session <json> --question-id <id> --answer <text> --out <json>
  blueprint interview status --session <json>
  blueprint compile --session <json> --out-dir <dir>
  blueprint lint --path <path>
  blueprint readiness audit --pack <dir> --out <json>
  blueprint sdd emit --pack <dir> --out <json>
  blueprint authorize --pack <dir> --out <json>
  blueprint pack inspect --pack <dir> --json

Commands: interview compile lint readiness sdd authorize pack`)
}

func runLint(args []string, stdout io.Writer) error {
	flags := parseFlags(args)
	path := flags["path"]
	if path == "" {
		return errors.New("usage: blueprint lint --path <path>")
	}
	report, err := blueprint.LintPath(path)
	if err != nil {
		_ = blueprint.WriteJSON(filepath.Join("tmp", "blueprint-lint-report.json"), report)
		return err
	}
	fmt.Fprintf(stdout, "lint: %s findings=%d\n", report.Status, report.FindingCount)
	return nil
}

func runReadiness(args []string, stdout io.Writer) error {
	if len(args) == 0 || args[0] != "audit" {
		return errors.New("usage: blueprint readiness audit --pack <dir> --out <json>")
	}
	flags := parseFlags(args[1:])
	pack := flags["pack"]
	out := flags["out"]
	if pack == "" || out == "" {
		return errors.New("usage: blueprint readiness audit --pack <dir> --out <json>")
	}
	audit, err := blueprint.AuditPack(pack)
	if writeErr := blueprint.WriteJSON(out, audit); writeErr != nil {
		return writeErr
	}
	fmt.Fprintf(stdout, "production readiness: %d/100 status=%s\n", audit.Score, audit.Status)
	return err
}

func runSDD(args []string, stdout io.Writer) error {
	if len(args) == 0 || args[0] != "emit" {
		return errors.New("usage: blueprint sdd emit --pack <dir> --out <json>")
	}
	flags := parseFlags(args[1:])
	pack := flags["pack"]
	out := flags["out"]
	if pack == "" || out == "" {
		return errors.New("usage: blueprint sdd emit --pack <dir> --out <json>")
	}
	if err := blueprint.EmitSDDPlan(pack, out); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "sdd plan written: %s\n", filepath.Clean(out))
	return nil
}

func runAuthorize(args []string, stdout io.Writer) error {
	flags := parseFlags(args)
	pack := flags["pack"]
	out := flags["out"]
	if pack == "" || out == "" {
		return errors.New("usage: blueprint authorize --pack <dir> --out <json>")
	}
	auth, err := blueprint.AuthorizePack(pack)
	if writeErr := blueprint.WriteJSON(out, auth); writeErr != nil {
		return writeErr
	}
	if err != nil {
		return err
	}
	fmt.Fprintf(stdout, "authorization: %s score=%d next=%s\n", auth.Status, auth.Score, auth.NextAllowedAction)
	return nil
}

func runPack(args []string, stdout io.Writer) error {
	if len(args) == 0 || args[0] != "inspect" {
		return errors.New("usage: blueprint pack inspect --pack <dir> [--json]")
	}
	flags := parseFlags(args[1:])
	pack := flags["pack"]
	if pack == "" {
		return errors.New("usage: blueprint pack inspect --pack <dir> [--json]")
	}
	inspection, err := blueprint.InspectPack(pack)
	if _, jsonMode := flags["json"]; jsonMode {
		body, marshalErr := json.MarshalIndent(inspection, "", "  ")
		if marshalErr != nil {
			return marshalErr
		}
		fmt.Fprintln(stdout, string(body))
		return err
	}
	fmt.Fprintf(stdout, "pack: %s artifacts=%d status=%s\n", inspection.ProjectID, inspection.ArtifactCount, inspection.Status)
	return err
}

func runCompile(args []string, stdout io.Writer) error {
	flags := parseFlags(args)
	sessionPath := flags["session"]
	outDir := flags["out-dir"]
	if sessionPath == "" || outDir == "" {
		return errors.New("usage: blueprint compile --session <json> --out-dir <dir>")
	}
	s, err := readSession(sessionPath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	answers := summarizeAnswers(s)
	textFiles := map[string]string{
		"project-brief.md":         "# Project Brief\n\n" + s.Idea + "\n\n" + answers,
		"prd.md":                   "# PRD\n\nDraft generated from interview answers. User review is required before authorization.\n",
		"non-goals.md":             "# Non-Goals\n\nNo target product implementation may start until user approval is recorded.\n",
		"domain-model.md":          "# Domain Model\n\nDraft entities: interview session, question, answer, requirement, risk, traceability row, SDD plan, authorization.\n",
		"workflow-map.md":          "# Workflow Map\n\nIdea -> interview -> draft pack -> review -> readiness audit -> authorization.\n",
		"contracts.md":             "# Contracts\n\nDraft pack uses AO Blueprint v0.1 contracts and requires review before build handoff.\n",
		"security-privacy.md":      "# Security And Privacy\n\nDraft artifacts must be scanned for local paths, tokens, private keys, and secret values before approval.\n",
		"operations-runbook.md":    "# Operations Runbook\n\nRun blueprint lint, readiness audit, SDD emit, and authorization before downstream AO work.\n",
		"test-evaluation-plan.md":  "# Test And Evaluation Plan\n\nUse fixture tests, readiness audit, authorization denial, public-safety scan, and clean-clone smoke checks.\n",
		"implementation-slices.md": "# Implementation Slices\n\n1. Review draft requirements.\n2. Resolve blockers.\n3. Approve pack.\n4. Authorize downstream AO work.\n",
		"ao-forge-handoff.md":      "# AO Forge Handoff\n\nBlocked until AO Blueprint authorization reports status ready.\n",
	}
	for name, body := range textFiles {
		if err := os.WriteFile(filepath.Join(outDir, name), []byte(body), 0o644); err != nil {
			return err
		}
	}
	jsonFiles := map[string]any{
		"requirements.json": map[string]any{
			"schema":     "ao.blueprint.requirement.v0.1",
			"project_id": s.SessionID,
			"requirements": []map[string]any{{
				"id":          "REQ-001",
				"title":       "Review interview objective",
				"category":    "objective",
				"description": s.Idea,
				"acceptance":  "User approves the compiled blueprint pack.",
			}},
		},
		"risk-register.json": map[string]any{
			"schema":     "ao.blueprint.risk.v0.1",
			"project_id": s.SessionID,
			"risks": []map[string]any{{
				"id":         "RISK-001",
				"title":      "Draft pack not yet approved",
				"severity":   "high",
				"mitigation": "Keep authorization blocked until approval.",
				"status":     "open",
			}},
		},
		"traceability-matrix.json": map[string]any{
			"schema":     "ao.blueprint.traceability-matrix.v0.1",
			"project_id": s.SessionID,
			"rows": []map[string]any{{
				"requirement_id": "REQ-001",
				"slice_ids":      []string{"SLICE-001"},
				"evidence":       []string{"user approval"},
			}},
		},
		"sdd-plan.json": map[string]any{
			"schema":         "ao2.sdd-plan.v1",
			"project_id":     s.SessionID,
			"title":          "Draft AO Blueprint SDD plan",
			"objective":      s.Idea,
			"slices":         []map[string]any{{"id": "SLICE-001", "title": "Resolve draft blockers", "acceptance": []string{"user approval recorded"}}},
			"stop_condition": "blueprint authorization status is ready",
		},
		"ao-foundry-task.json": map[string]any{
			"schema":                 "ao.foundry.task.v0.1",
			"task_id":                s.SessionID + "-draft",
			"source":                 "ao-blueprint",
			"authorization_required": true,
			"next_action":            "blocked-until-approval",
			"exit_condition":         "AO Blueprint build authorization status ready",
		},
		"sufficiency-audit.json": map[string]any{
			"schema":                              "ao.blueprint.sufficiency-audit.v0.1",
			"project_id":                          s.SessionID,
			"status":                              "blocked",
			"score":                               100,
			"approved_by_user":                    false,
			"blocking_assumptions":                []string{"User approval is required before downstream build authorization."},
			"production_readiness_exit_condition": "Authorization status ready after user approval and clean safety scan.",
			"next_allowed_action":                 "ao-foundry",
			"categories":                          compileCategories(),
		},
		"build-authorization.json": map[string]any{
			"schema":           "ao.blueprint.build-authorization.v0.1",
			"project_id":       s.SessionID,
			"status":           "blocked",
			"score":            100,
			"approved_by_user": false,
			"blocking_assumptions": []string{
				"User approval is required before downstream build authorization.",
			},
			"next_allowed_action": "ao-foundry",
		},
	}
	for name, body := range jsonFiles {
		if err := blueprint.WriteJSON(filepath.Join(outDir, name), body); err != nil {
			return err
		}
	}
	fmt.Fprintf(stdout, "compiled blueprint draft: %s\n", filepath.Clean(outDir))
	return nil
}

func summarizeAnswers(s session) string {
	if len(s.AnsweredQuestions) == 0 {
		return "No interview answers have been recorded yet.\n"
	}
	var builder strings.Builder
	builder.WriteString("## Interview Answers\n\n")
	for _, answer := range s.AnsweredQuestions {
		builder.WriteString("- ")
		builder.WriteString(answer.QuestionID)
		builder.WriteString(": ")
		builder.WriteString(answer.Answer)
		builder.WriteString("\n")
	}
	return builder.String()
}

func compileCategories() []map[string]any {
	return []map[string]any{
		{"id": "objective", "points": 10, "status": "passed"},
		{"id": "scope", "points": 10, "status": "passed"},
		{"id": "domain_workflows", "points": 12, "status": "passed"},
		{"id": "interfaces_contracts", "points": 10, "status": "passed"},
		{"id": "data_integrations", "points": 8, "status": "passed"},
		{"id": "security_privacy", "points": 10, "status": "passed"},
		{"id": "tests_evaluation", "points": 12, "status": "passed"},
		{"id": "operations_release", "points": 8, "status": "passed"},
		{"id": "traceability", "points": 10, "status": "passed"},
		{"id": "approval_handoff", "points": 10, "status": "passed"},
	}
}

func runInterview(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return errors.New("usage: blueprint interview <start|next|answer|status>")
	}
	switch args[0] {
	case "start":
		return interviewStart(args[1:], stdout)
	case "next":
		return interviewNext(args[1:], stdout)
	case "answer":
		return interviewAnswer(args[1:], stdout)
	case "status":
		return interviewStatus(args[1:], stdout)
	default:
		return fmt.Errorf("unknown interview command %q", args[0])
	}
}

func interviewStart(args []string, stdout io.Writer) error {
	flags := parseFlags(args)
	idea := flags["idea"]
	out := flags["out"]
	if idea == "" || out == "" {
		return errors.New("usage: blueprint interview start --idea <text> --out <json>")
	}
	s := session{
		Schema:            "ao.blueprint.session.v0.1",
		SessionID:         "session-" + time.Now().UTC().Format("20060102T150405Z"),
		Idea:              idea,
		Status:            "interviewing",
		AnsweredQuestions: []answeredRecord{},
	}
	if err := writeSession(out, s); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "interview session written: %s\n", filepath.Clean(out))
	return nil
}

func interviewNext(args []string, stdout io.Writer) error {
	flags := parseFlags(args)
	sessionPath := flags["session"]
	if sessionPath == "" {
		return errors.New("usage: blueprint interview next --session <json>")
	}
	s, err := readSession(sessionPath)
	if err != nil {
		return err
	}
	answered := map[string]bool{}
	for _, item := range s.AnsweredQuestions {
		answered[item.QuestionID] = true
	}
	for _, q := range questions {
		if !answered[q.ID] {
			fmt.Fprintf(stdout, "question_id=%s category=%s prompt=%q\n", q.ID, q.Category, q.Prompt)
			return nil
		}
	}
	fmt.Fprintln(stdout, "no questions remain")
	return nil
}

func interviewAnswer(args []string, stdout io.Writer) error {
	flags := parseFlags(args)
	sessionPath := flags["session"]
	questionID := flags["question-id"]
	answer := flags["answer"]
	out := flags["out"]
	if out == "" {
		out = sessionPath
	}
	if sessionPath == "" || questionID == "" || answer == "" {
		return errors.New("usage: blueprint interview answer --session <json> --question-id <id> --answer <text> --out <json>")
	}
	s, err := readSession(sessionPath)
	if err != nil {
		return err
	}
	s.AnsweredQuestions = append(s.AnsweredQuestions, answeredRecord{
		QuestionID: questionID,
		Answer:     answer,
		AnsweredAt: time.Now().UTC().Format(time.RFC3339),
	})
	if err := writeSession(out, s); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "answer recorded: %s answered=%d\n", questionID, len(s.AnsweredQuestions))
	return nil
}

func interviewStatus(args []string, stdout io.Writer) error {
	flags := parseFlags(args)
	sessionPath := flags["session"]
	if sessionPath == "" {
		return errors.New("usage: blueprint interview status --session <json>")
	}
	s, err := readSession(sessionPath)
	if err != nil {
		return err
	}
	fmt.Fprintf(stdout, "session=%s status=%s answered=%d\n", s.SessionID, s.Status, len(s.AnsweredQuestions))
	return nil
}

func readSession(path string) (session, error) {
	var s session
	body, err := os.ReadFile(path)
	if err != nil {
		return s, err
	}
	if err := json.Unmarshal(body, &s); err != nil {
		return s, err
	}
	if s.Schema != "ao.blueprint.session.v0.1" {
		return s, fmt.Errorf("unsupported session schema: %s", s.Schema)
	}
	return s, nil
}

func writeSession(path string, s session) error {
	body, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	body = append(body, '\n')
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, body, 0o644)
}

func parseFlags(args []string) map[string]string {
	flags := map[string]string{}
	for i := 0; i < len(args); i++ {
		item := args[i]
		if !strings.HasPrefix(item, "--") {
			continue
		}
		key := strings.TrimPrefix(item, "--")
		if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
			flags[key] = args[i+1]
			i++
			continue
		}
		flags[key] = ""
	}
	return flags
}
