package blueprint

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	AuditSchema         = "ao.blueprint.sufficiency-audit.v0.1"
	AuthorizationSchema = "ao.blueprint.build-authorization.v0.1"
	LintSchema          = "ao.blueprint.lint-report.v0.1"
	InspectionSchema    = "ao.blueprint.pack-inspection.v0.1"
)

type Diagnostic struct {
	Code       string `json:"code"`
	Severity   string `json:"severity"`
	Message    string `json:"message"`
	Path       string `json:"path,omitempty"`
	NextAction string `json:"next_action,omitempty"`
}

type CategoryScore struct {
	ID      string   `json:"id"`
	Points  int      `json:"points"`
	Status  string   `json:"status"`
	Reasons []string `json:"reasons,omitempty"`
}

type SufficiencyAudit struct {
	Schema                           string          `json:"schema"`
	ProjectID                        string          `json:"project_id"`
	Status                           string          `json:"status"`
	Score                            int             `json:"score"`
	ApprovedByUser                   bool            `json:"approved_by_user"`
	BlockingAssumptions              []string        `json:"blocking_assumptions"`
	ProductionReadinessExitCondition string          `json:"production_readiness_exit_condition"`
	NextAllowedAction                string          `json:"next_allowed_action"`
	Categories                       []CategoryScore `json:"categories"`
	Blockers                         []Diagnostic    `json:"blockers"`
	Checks                           []string        `json:"checks"`
}

type BuildAuthorization struct {
	Schema                           string       `json:"schema"`
	ProjectID                        string       `json:"project_id"`
	Status                           string       `json:"status"`
	Score                            int          `json:"score"`
	ApprovedByUser                   bool         `json:"approved_by_user"`
	BlockingAssumptions              []string     `json:"blocking_assumptions"`
	BlueprintPackDigest              string       `json:"blueprint_pack_digest,omitempty"`
	RequirementsDigest               string       `json:"requirements_digest,omitempty"`
	TraceabilityDigest               string       `json:"traceability_digest,omitempty"`
	SDDPlanDigest                    string       `json:"sdd_plan_digest,omitempty"`
	SDDPlanPath                      string       `json:"sdd_plan_path,omitempty"`
	AOForgeHandoffPath               string       `json:"ao_forge_handoff_path,omitempty"`
	AOFoundryTaskPath                string       `json:"ao_foundry_task_path,omitempty"`
	ProductionReadinessExitCondition string       `json:"production_readiness_exit_condition"`
	NextAllowedAction                string       `json:"next_allowed_action"`
	Blockers                         []Diagnostic `json:"blockers,omitempty"`
}

type LintFinding struct {
	Path    string `json:"path"`
	Line    int    `json:"line"`
	Kind    string `json:"kind"`
	Message string `json:"message"`
}

type LintReport struct {
	Schema       string        `json:"schema"`
	Status       string        `json:"status"`
	FindingCount int           `json:"finding_count"`
	Findings     []LintFinding `json:"findings"`
}

type PackInspection struct {
	Schema        string       `json:"schema"`
	Status        string       `json:"status"`
	ProjectID     string       `json:"project_id"`
	ArtifactCount int          `json:"artifact_count"`
	Artifacts     []string     `json:"artifacts"`
	Blockers      []Diagnostic `json:"blockers,omitempty"`
}

type categorySpec struct {
	id     string
	points int
	files  []string
}

var categorySpecs = []categorySpec{
	{id: "objective", points: 10, files: []string{"project-brief.md", "prd.md"}},
	{id: "scope", points: 10, files: []string{"non-goals.md"}},
	{id: "domain_workflows", points: 12, files: []string{"domain-model.md", "workflow-map.md"}},
	{id: "interfaces_contracts", points: 10, files: []string{"contracts.md", "requirements.json"}},
	{id: "data_integrations", points: 8, files: []string{"requirements.json"}},
	{id: "security_privacy", points: 10, files: []string{"security-privacy.md"}},
	{id: "tests_evaluation", points: 12, files: []string{"test-evaluation-plan.md", "implementation-slices.md"}},
	{id: "operations_release", points: 8, files: []string{"operations-runbook.md"}},
	{id: "traceability", points: 10, files: []string{"traceability-matrix.json"}},
	{id: "approval_handoff", points: 10, files: []string{"sufficiency-audit.json", "sdd-plan.json", "ao-forge-handoff.md", "ao-foundry-task.json"}},
}

var jsonArtifacts = []string{
	"requirements.json",
	"risk-register.json",
	"traceability-matrix.json",
	"sufficiency-audit.json",
	"sdd-plan.json",
	"ao-foundry-task.json",
	"build-authorization.json",
}

var localPathPattern = "(" + "/" + "Users/|" + "/" + "Volumes/|" + "/" + "private/|" + "C:" + `\\Users` + ")"

var secretPatterns = []struct {
	kind string
	re   *regexp.Regexp
}{
	{kind: "local_path", re: regexp.MustCompile(localPathPattern)},
	{kind: "bearer_token", re: regexp.MustCompile(`Authorization:\s*Bearer\s+[A-Za-z0-9._-]+`)},
	{kind: "private_key", re: regexp.MustCompile(`BEGIN (RSA |OPENSSH |PRIVATE )?PRIVATE KEY`)},
	{kind: "openai_key", re: regexp.MustCompile(`sk-[A-Za-z0-9]{20,}`)},
	{kind: "github_token", re: regexp.MustCompile(`gh[pousr]_[A-Za-z0-9]{20,}`)},
	{kind: "aws_key", re: regexp.MustCompile(`AKIA[0-9A-Z]{16}`)},
	{kind: "secret_assignment", re: regexp.MustCompile(`(?i)\b(password|passwd|token|cookie)\s*[:=]\s*[^[:space:]]+`)},
}

func AuditPack(pack string) (SufficiencyAudit, error) {
	pack = filepath.Clean(pack)
	declared, declaredErr := readDeclaredAudit(pack)
	audit := SufficiencyAudit{
		Schema:                           AuditSchema,
		ProjectID:                        declared.ProjectID,
		ApprovedByUser:                   declared.ApprovedByUser,
		BlockingAssumptions:              append([]string{}, declared.BlockingAssumptions...),
		ProductionReadinessExitCondition: declared.ProductionReadinessExitCondition,
		NextAllowedAction:                declared.NextAllowedAction,
		Checks:                           []string{"required_artifacts", "json_parse", "public_safety", "approval", "handoff"},
	}
	if audit.ProjectID == "" {
		audit.ProjectID = filepath.Base(pack)
	}
	if declaredErr != nil {
		audit.Blockers = append(audit.Blockers, blocker("audit_read", "cannot read sufficiency audit: "+declaredErr.Error(), "sufficiency-audit.json"))
	}

	if info, err := os.Stat(pack); err != nil || !info.IsDir() {
		audit.Blockers = append(audit.Blockers, blocker("pack_missing", "blueprint pack directory is missing", pack))
		return finalizeAudit(audit), fmt.Errorf("blueprint pack not found: %s", pack)
	}

	missingByFile := map[string]bool{}
	for _, spec := range categorySpecs {
		category := CategoryScore{ID: spec.id, Points: spec.points, Status: "passed"}
		for _, name := range spec.files {
			path := filepath.Join(pack, name)
			if err := requireNonEmptyFile(path); err != nil {
				rel := filepath.ToSlash(name)
				missingByFile[rel] = true
				category.Status = "failed"
				category.Reasons = append(category.Reasons, fmt.Sprintf("%s: %v", rel, err))
			}
		}
		if category.Status == "passed" {
			audit.Score += category.Points
		} else {
			audit.Blockers = append(audit.Blockers, blocker("artifact_missing", "required blueprint artifact missing or empty", strings.Join(category.Reasons, "; ")))
		}
		audit.Categories = append(audit.Categories, category)
	}

	for _, name := range jsonArtifacts {
		if missingByFile[name] {
			continue
		}
		path := filepath.Join(pack, name)
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err := validateJSONFile(path); err != nil {
			audit.Blockers = append(audit.Blockers, blocker("invalid_json", err.Error(), name))
		}
	}

	lint, lintErr := LintPath(pack)
	if lintErr != nil {
		for _, finding := range lint.Findings {
			audit.Blockers = append(audit.Blockers, Diagnostic{
				Code:       "unsafe_public_artifact",
				Severity:   "critical",
				Message:    finding.Message,
				Path:       finding.Path,
				NextAction: "remove or redact unsafe public artifact content",
			})
		}
	}

	if !audit.ApprovedByUser {
		audit.Blockers = append(audit.Blockers, blocker("approval_missing", "user approval is required before build authorization", "sufficiency-audit.json"))
	}
	if len(audit.BlockingAssumptions) > 0 {
		audit.Blockers = append(audit.Blockers, blocker("blocking_assumptions", "blocking assumptions remain", "sufficiency-audit.json"))
	}
	if strings.TrimSpace(audit.ProductionReadinessExitCondition) == "" {
		audit.Blockers = append(audit.Blockers, blocker("exit_condition_missing", "production-readiness exit condition is required", "sufficiency-audit.json"))
	}
	if audit.NextAllowedAction != "ao-foundry" && audit.NextAllowedAction != "ao-forge" {
		audit.Blockers = append(audit.Blockers, blocker("handoff_target_invalid", "next allowed action must be ao-foundry or ao-forge", "sufficiency-audit.json"))
	}
	if declared.Score > 0 && declared.Score != audit.Score {
		audit.Blockers = append(audit.Blockers, blocker("score_mismatch", fmt.Sprintf("declared score %d does not match computed score %d", declared.Score, audit.Score), "sufficiency-audit.json"))
	}

	return finalizeAudit(audit), nil
}

func AuthorizePack(pack string) (BuildAuthorization, error) {
	audit, err := AuditPack(pack)
	auth := BuildAuthorization{
		Schema:                           AuthorizationSchema,
		ProjectID:                        audit.ProjectID,
		Status:                           "blocked",
		Score:                            audit.Score,
		ApprovedByUser:                   audit.ApprovedByUser,
		BlockingAssumptions:              audit.BlockingAssumptions,
		ProductionReadinessExitCondition: audit.ProductionReadinessExitCondition,
		NextAllowedAction:                audit.NextAllowedAction,
		Blockers:                         audit.Blockers,
		SDDPlanPath:                      "sdd-plan.json",
		AOForgeHandoffPath:               "ao-forge-handoff.md",
		AOFoundryTaskPath:                "ao-foundry-task.json",
	}
	if err != nil {
		return auth, err
	}
	if audit.Status != "ready" {
		return auth, fmt.Errorf("authorization blocked: readiness status=%s score=%d blockers=%d", audit.Status, audit.Score, len(audit.Blockers))
	}

	var digestErr error
	auth.BlueprintPackDigest, digestErr = digestDir(pack)
	if digestErr != nil {
		return auth, digestErr
	}
	auth.RequirementsDigest, digestErr = digestFile(filepath.Join(pack, "requirements.json"))
	if digestErr != nil {
		return auth, digestErr
	}
	auth.TraceabilityDigest, digestErr = digestFile(filepath.Join(pack, "traceability-matrix.json"))
	if digestErr != nil {
		return auth, digestErr
	}
	auth.SDDPlanDigest, digestErr = digestFile(filepath.Join(pack, "sdd-plan.json"))
	if digestErr != nil {
		return auth, digestErr
	}
	auth.Status = "ready"
	auth.Blockers = nil
	return auth, nil
}

func EmitSDDPlan(pack string, out string) error {
	source := filepath.Join(filepath.Clean(pack), "sdd-plan.json")
	if err := validateJSONFile(source); err != nil {
		return err
	}
	body, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	var doc map[string]any
	if err := json.Unmarshal(body, &doc); err != nil {
		return err
	}
	if doc["schema"] != "ao2.sdd-plan.v1" {
		return fmt.Errorf("sdd plan schema = %v, want ao2.sdd-plan.v1", doc["schema"])
	}
	return writeFile(out, body)
}

func InspectPack(pack string) (PackInspection, error) {
	audit, err := AuditPack(pack)
	inspection := PackInspection{
		Schema:    InspectionSchema,
		Status:    audit.Status,
		ProjectID: audit.ProjectID,
		Blockers:  audit.Blockers,
	}
	seen := map[string]bool{}
	for _, spec := range categorySpecs {
		for _, name := range spec.files {
			if seen[name] {
				continue
			}
			seen[name] = true
			if _, statErr := os.Stat(filepath.Join(pack, name)); statErr == nil {
				inspection.Artifacts = append(inspection.Artifacts, filepath.ToSlash(name))
			}
		}
	}
	sort.Strings(inspection.Artifacts)
	inspection.ArtifactCount = len(inspection.Artifacts)
	return inspection, err
}

func LintPath(path string) (LintReport, error) {
	report := LintReport{Schema: LintSchema, Status: "passed"}
	err := filepath.WalkDir(filepath.Clean(path), func(current string, entry fs.DirEntry, err error) error {
		if err != nil {
			report.Findings = append(report.Findings, LintFinding{Path: current, Kind: "walk_error", Message: err.Error()})
			return nil
		}
		if entry.IsDir() {
			if shouldSkipDir(entry.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if !isTextCandidate(current) {
			return nil
		}
		body, readErr := os.ReadFile(current)
		if readErr != nil {
			report.Findings = append(report.Findings, LintFinding{Path: current, Kind: "read_error", Message: readErr.Error()})
			return nil
		}
		if strings.HasSuffix(current, ".json") && !json.Valid(body) {
			report.Findings = append(report.Findings, LintFinding{Path: filepath.ToSlash(current), Kind: "invalid_json", Message: "JSON file does not parse"})
		}
		for i, line := range strings.Split(string(body), "\n") {
			for _, pattern := range secretPatterns {
				if pattern.re.MatchString(line) {
					report.Findings = append(report.Findings, LintFinding{
						Path:    filepath.ToSlash(current),
						Line:    i + 1,
						Kind:    pattern.kind,
						Message: "unsafe public artifact content detected",
					})
				}
			}
		}
		return nil
	})
	if err != nil {
		return report, err
	}
	report.FindingCount = len(report.Findings)
	if report.FindingCount > 0 {
		report.Status = "failed"
		return report, fmt.Errorf("lint failed with %d findings", report.FindingCount)
	}
	return report, nil
}

func WriteJSON(path string, value any) error {
	body, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	body = append(body, '\n')
	return writeFile(path, body)
}

func readDeclaredAudit(pack string) (SufficiencyAudit, error) {
	path := filepath.Join(pack, "sufficiency-audit.json")
	var audit SufficiencyAudit
	body, err := os.ReadFile(path)
	if err != nil {
		return audit, err
	}
	if err := json.Unmarshal(body, &audit); err != nil {
		return audit, err
	}
	return audit, nil
}

func finalizeAudit(audit SufficiencyAudit) SufficiencyAudit {
	if audit.Score == 100 && len(audit.Blockers) == 0 {
		audit.Status = "ready"
	} else {
		audit.Status = "blocked"
	}
	if audit.BlockingAssumptions == nil {
		audit.BlockingAssumptions = []string{}
	}
	if audit.Blockers == nil {
		audit.Blockers = []Diagnostic{}
	}
	return audit
}

func blocker(code string, message string, path string) Diagnostic {
	return Diagnostic{
		Code:       code,
		Severity:   "critical",
		Message:    message,
		Path:       filepath.ToSlash(path),
		NextAction: "fix blueprint pack and rerun readiness audit",
	}
}

func requireNonEmptyFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("is a directory")
	}
	if info.Size() == 0 {
		return fmt.Errorf("empty file")
	}
	return nil
}

func validateJSONFile(path string) error {
	body, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if !json.Valid(body) {
		return fmt.Errorf("invalid JSON: %s", path)
	}
	return nil
}

func writeFile(path string, body []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, body, 0o644)
}

func digestFile(path string) (string, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(body)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func digestDir(root string) (string, error) {
	hash := sha256.New()
	err := filepath.WalkDir(filepath.Clean(root), func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if shouldSkipDir(entry.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		hash.Write([]byte(filepath.ToSlash(rel)))
		hash.Write([]byte{0})
		hash.Write(body)
		hash.Write([]byte{0})
		return nil
	})
	if err != nil {
		return "", err
	}
	return "sha256:" + hex.EncodeToString(hash.Sum(nil)), nil
}

func shouldSkipDir(name string) bool {
	switch name {
	case ".git", "tmp", "target", ".idea", ".vscode", "__pycache__":
		return true
	default:
		return false
	}
}

func isTextCandidate(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go", ".md", ".json", ".yml", ".yaml", ".txt", ".sh", ".ps1", ".mod", ".sum":
		return true
	default:
		return false
	}
}
