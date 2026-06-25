package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func rootDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime caller unavailable")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func runCLI(args ...string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := Run(args, &stdout, &stderr)
	return stdout.String(), stderr.String(), err
}

func TestHelpListsCoreCommands(t *testing.T) {
	stdout, _, err := runCLI("--help")
	if err != nil {
		t.Fatalf("help returned error: %v", err)
	}
	for _, want := range []string{"interview", "compile", "lint", "readiness", "sdd", "authorize", "pack"} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("help missing %q in %s", want, stdout)
		}
	}
}

func TestReadinessAuditCommandWritesReadyJSON(t *testing.T) {
	pack := filepath.Join(rootDir(t), "examples", "blueprints", "valid", "ao-blueprint-self")
	out := filepath.Join(t.TempDir(), "readiness.json")

	stdout, _, err := runCLI("readiness", "audit", "--pack", pack, "--out", out)
	if err != nil {
		t.Fatalf("readiness audit returned error: %v", err)
	}
	if !strings.Contains(stdout, "production readiness: 100/100 status=ready") {
		t.Fatalf("unexpected stdout: %s", stdout)
	}
	var body map[string]any
	readJSON(t, out, &body)
	if body["status"] != "ready" || int(body["score"].(float64)) != 100 {
		t.Fatalf("audit body = %#v, want ready score 100", body)
	}
}

func TestAuthorizeCommandBlocksInvalidPack(t *testing.T) {
	pack := filepath.Join(rootDir(t), "examples", "blueprints", "invalid", "missing-approval")
	out := filepath.Join(t.TempDir(), "authorization.json")

	_, stderr, err := runCLI("authorize", "--pack", pack, "--out", out)
	if err == nil {
		t.Fatal("authorize returned nil error for invalid pack")
	}
	if !strings.Contains(stderr, "authorization blocked") {
		t.Fatalf("stderr = %q, want authorization blocked", stderr)
	}
}

func TestInterviewCommandsAdvanceSession(t *testing.T) {
	session := filepath.Join(t.TempDir(), "session.json")

	if _, _, err := runCLI("interview", "start", "--idea", "Build a governed requirements gate", "--out", session); err != nil {
		t.Fatalf("interview start returned error: %v", err)
	}
	stdout, _, err := runCLI("interview", "next", "--session", session)
	if err != nil {
		t.Fatalf("interview next returned error: %v", err)
	}
	if !strings.Contains(stdout, "question_id=") {
		t.Fatalf("next stdout = %q, want question_id", stdout)
	}
	if _, _, err := runCLI("interview", "answer", "--session", session, "--question-id", "q-objective", "--answer", "Success is a 100 point authorization gate", "--out", session); err != nil {
		t.Fatalf("interview answer returned error: %v", err)
	}
	status, _, err := runCLI("interview", "status", "--session", session)
	if err != nil {
		t.Fatalf("interview status returned error: %v", err)
	}
	if !strings.Contains(status, "answered=1") {
		t.Fatalf("status stdout = %q, want answered=1", status)
	}
}

func TestCompileCommandWritesBlockedDraftPack(t *testing.T) {
	tmp := t.TempDir()
	session := filepath.Join(tmp, "session.json")
	outDir := filepath.Join(tmp, "pack")

	if _, _, err := runCLI("interview", "start", "--idea", "Build a governed requirements gate", "--out", session); err != nil {
		t.Fatalf("interview start returned error: %v", err)
	}
	if _, _, err := runCLI("interview", "answer", "--session", session, "--question-id", "q-objective", "--answer", "Success is a reviewed blueprint pack", "--out", session); err != nil {
		t.Fatalf("interview answer returned error: %v", err)
	}
	if _, _, err := runCLI("compile", "--session", session, "--out-dir", outDir); err != nil {
		t.Fatalf("compile returned error: %v", err)
	}

	var audit map[string]any
	readJSON(t, filepath.Join(outDir, "sufficiency-audit.json"), &audit)
	if audit["approved_by_user"] != false {
		t.Fatalf("compiled draft approval = %#v, want false", audit["approved_by_user"])
	}
	if audit["status"] != "blocked" {
		t.Fatalf("compiled draft status = %#v, want blocked", audit["status"])
	}
}

func readJSON(t *testing.T, path string, out any) {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if err := json.Unmarshal(raw, out); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
}
