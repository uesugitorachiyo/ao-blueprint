package blueprint

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime caller unavailable")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func TestValidPackAuditsToReadyScore100(t *testing.T) {
	pack := filepath.Join(repoRoot(t), "examples", "blueprints", "valid", "ao-blueprint-self")

	audit, err := AuditPack(pack)
	if err != nil {
		t.Fatalf("AuditPack returned error: %v", err)
	}
	if audit.Status != "ready" {
		t.Fatalf("status = %q, want ready; blockers=%v", audit.Status, audit.Blockers)
	}
	if audit.Score != 100 {
		t.Fatalf("score = %d, want 100", audit.Score)
	}
	if len(audit.Blockers) != 0 {
		t.Fatalf("blockers = %#v, want none", audit.Blockers)
	}
}

func TestMissingImplementationSpecBlocksReadiness(t *testing.T) {
	source := filepath.Join(repoRoot(t), "examples", "blueprints", "valid", "ao-blueprint-self")
	pack := filepath.Join(t.TempDir(), "pack")
	if err := copyDirForTest(source, pack); err != nil {
		t.Fatalf("copy valid pack: %v", err)
	}
	if err := os.Remove(filepath.Join(pack, "implementation-spec.md")); err != nil {
		t.Fatalf("remove implementation spec: %v", err)
	}

	audit, err := AuditPack(pack)
	if err != nil {
		t.Fatalf("AuditPack returned unexpected error: %v", err)
	}
	if audit.Status != "blocked" {
		t.Fatalf("status = %q, want blocked", audit.Status)
	}
	if audit.Score >= 100 {
		t.Fatalf("score = %d, want below 100", audit.Score)
	}
	if !diagnosticsContainPath(audit.Blockers, "implementation-spec.md") {
		t.Fatalf("blockers = %#v, want implementation-spec.md blocker", audit.Blockers)
	}
}

func TestMissingApprovalPackBlocksAuthorization(t *testing.T) {
	pack := filepath.Join(repoRoot(t), "examples", "blueprints", "invalid", "missing-approval")

	auth, err := AuthorizePack(pack)
	if err == nil {
		t.Fatal("AuthorizePack returned nil error for blocked pack")
	}
	if auth.Status != "blocked" {
		t.Fatalf("status = %q, want blocked", auth.Status)
	}
	if !strings.Contains(err.Error(), "authorization blocked") {
		t.Fatalf("error = %q, want authorization blocked", err.Error())
	}
}

func TestEmitSDDPlanCopiesValidAO2Plan(t *testing.T) {
	pack := filepath.Join(repoRoot(t), "examples", "blueprints", "valid", "ao-blueprint-self")
	out := filepath.Join(t.TempDir(), "sdd-plan.json")

	if err := EmitSDDPlan(pack, out); err != nil {
		t.Fatalf("EmitSDDPlan returned error: %v", err)
	}
	raw, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read emitted SDD: %v", err)
	}
	if !strings.Contains(string(raw), `"schema": "ao2.sdd-plan.v1"`) {
		t.Fatalf("emitted SDD missing ao2 schema: %s", raw)
	}
}

func TestLintFindsUnsafeLocalPath(t *testing.T) {
	dir := t.TempDir()
	unsafeFile := filepath.Join(dir, "unsafe.md")
	unsafePath := "leaks " + "/" + "Users/example/private.txt"
	if err := os.WriteFile(unsafeFile, []byte(unsafePath), 0o600); err != nil {
		t.Fatalf("write unsafe file: %v", err)
	}

	report, err := LintPath(dir)
	if err == nil {
		t.Fatal("LintPath returned nil error for unsafe path")
	}
	if report.Status != "failed" {
		t.Fatalf("status = %q, want failed", report.Status)
	}
	if report.FindingCount == 0 {
		t.Fatal("finding count = 0, want unsafe finding")
	}
}

func TestInspectPackReportsRequiredArtifacts(t *testing.T) {
	pack := filepath.Join(repoRoot(t), "examples", "blueprints", "valid", "ao-blueprint-self")

	inspection, err := InspectPack(pack)
	if err != nil {
		t.Fatalf("InspectPack returned error: %v", err)
	}
	if inspection.ArtifactCount < 10 {
		t.Fatalf("artifact count = %d, want at least 10", inspection.ArtifactCount)
	}
	if inspection.Status != "ready" {
		t.Fatalf("status = %q, want ready", inspection.Status)
	}
}

func diagnosticsContainPath(items []Diagnostic, want string) bool {
	for _, item := range items {
		if strings.Contains(item.Path, want) {
			return true
		}
	}
	return false
}

func copyDirForTest(source string, target string) error {
	return filepath.WalkDir(source, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		dst := filepath.Join(target, rel)
		if entry.IsDir() {
			return os.MkdirAll(dst, 0o755)
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(dst, body, 0o644)
	})
}
