package blueprint

import (
	"encoding/json"
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

func TestBoundedRSIControlSurfaceReadbackPackAuthorizesToFoundry(t *testing.T) {
	pack := filepath.Join(repoRoot(t), "examples", "blueprints", "valid", "bounded-governed-rsi-control-surface-readback")

	audit, err := AuditPack(pack)
	if err != nil {
		t.Fatalf("AuditPack returned error: %v", err)
	}
	if audit.ProjectID != "bounded-governed-rsi-control-surface-readback" {
		t.Fatalf("project_id = %q, want bounded-governed-rsi-control-surface-readback", audit.ProjectID)
	}
	if audit.Status != "ready" || audit.Score != 100 {
		t.Fatalf("audit status=%q score=%d blockers=%v, want ready score=100", audit.Status, audit.Score, audit.Blockers)
	}

	auth, err := AuthorizePack(pack)
	if err != nil {
		t.Fatalf("AuthorizePack returned error: %v", err)
	}
	if auth.Status != "ready" ||
		auth.Score != 100 ||
		!auth.ApprovedByUser ||
		auth.NextAllowedAction != "ao-foundry" ||
		auth.ProjectID != "bounded-governed-rsi-control-surface-readback" {
		t.Fatalf("authorization drifted: %+v", auth)
	}
	if strings.Contains(strings.Join(auth.BlockingAssumptions, " "), "self-authorized") {
		t.Fatalf("authorization must not allow Blueprint self-authorization: %+v", auth)
	}
}

func TestAtlasRequiredPackRoutesReadinessAndAuthorizationToAtlas(t *testing.T) {
	source := filepath.Join(repoRoot(t), "examples", "blueprints", "valid", "ao-blueprint-self")
	pack := filepath.Join(t.TempDir(), "pack")
	if err := copyDirForTest(source, pack); err != nil {
		t.Fatalf("copy valid pack: %v", err)
	}
	writeJSONForTest(t, filepath.Join(pack, "ao-foundry-task.json"), map[string]any{
		"schema":                 "ao.foundry.task.v0.1",
		"task_id":                "atlas-required-mutation-class",
		"source":                 "ao-blueprint",
		"authorization_required": true,
		"atlas_required":         true,
		"summary":                "Compile a mutation-class workgraph, context packs, candidate records, and Foundry import material.",
		"exit_condition":         "AO Atlas emits the first safe Foundry import.",
	})

	audit, err := AuditPack(pack)
	if err != nil {
		t.Fatalf("AuditPack returned error: %v", err)
	}
	if audit.Status != "ready" || audit.Score != 100 {
		t.Fatalf("audit status=%q score=%d blockers=%v, want ready score=100", audit.Status, audit.Score, audit.Blockers)
	}
	if audit.NextAllowedAction != "ao-atlas" {
		t.Fatalf("audit next_allowed_action = %q, want ao-atlas", audit.NextAllowedAction)
	}

	auth, err := AuthorizePack(pack)
	if err != nil {
		t.Fatalf("AuthorizePack returned error: %v", err)
	}
	if auth.NextAllowedAction != audit.NextAllowedAction {
		t.Fatalf("authorization route = %q, audit route = %q; want agreement", auth.NextAllowedAction, audit.NextAllowedAction)
	}
	if auth.NextAllowedAction != "ao-atlas" {
		t.Fatalf("authorization next_allowed_action = %q, want ao-atlas", auth.NextAllowedAction)
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

func TestMissingQualityProfileBlocksReadiness(t *testing.T) {
	source := filepath.Join(repoRoot(t), "examples", "blueprints", "valid", "ao-blueprint-self")
	pack := filepath.Join(t.TempDir(), "pack")
	if err := copyDirForTest(source, pack); err != nil {
		t.Fatalf("copy valid pack: %v", err)
	}
	if err := os.Remove(filepath.Join(pack, "quality-profile.md")); err != nil {
		t.Fatalf("remove quality profile: %v", err)
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
	if !diagnosticsContainPath(audit.Blockers, "quality-profile.md") {
		t.Fatalf("blockers = %#v, want quality-profile.md blocker", audit.Blockers)
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

func writeJSONForTest(t *testing.T, path string, value any) {
	t.Helper()
	body, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatalf("marshal %s: %v", path, err)
	}
	body = append(body, '\n')
	if err := os.WriteFile(path, body, 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
