package cli

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
)

const (
	rehearsalVersion  = "v1.2.3"
	rehearsalSource   = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	rehearsalManifest = "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
)

func TestReleaseRehearsalWorkflowContract(t *testing.T) {
	workflow := readReleaseRehearsalFile(t, ".github", "workflows", "release-rehearsal.yml")

	triggers := directYAMLMapping(t, topLevelYAMLSection(t, workflow, "on"), 2)
	if len(triggers) != 1 || triggers["workflow_dispatch"] != "" {
		t.Fatalf("workflow triggers = %#v, want workflow_dispatch only", triggers)
	}
	dispatch := nestedYAMLSection(t, topLevelYAMLSection(t, workflow, "on"), "workflow_dispatch", 2)
	inputs := directYAMLMapping(t, nestedYAMLSection(t, dispatch, "inputs", 4), 6)
	for _, input := range []string{"version", "tag", "source_commit", "approved_manifest_digest"} {
		if _, ok := inputs[input]; !ok {
			t.Fatalf("workflow_dispatch inputs missing %q", input)
		}
	}

	permissions := directYAMLMapping(t, topLevelYAMLSection(t, workflow, "permissions"), 2)
	if len(permissions) != 1 || permissions["contents"] != "read" {
		t.Fatalf("workflow permissions = %#v, want contents: read only", permissions)
	}

	jobsSection := topLevelYAMLSection(t, workflow, "jobs")
	jobs := directYAMLMapping(t, jobsSection, 2)
	expectedJobs := map[string]bool{
		"bind-release-inputs":     true,
		"build-native-candidates": true,
		"assemble-promotion-plan": true,
	}
	if len(jobs) != len(expectedJobs) {
		t.Fatalf("workflow jobs = %#v, want exactly the three rehearsal jobs", jobs)
	}
	allowedActions := map[string]bool{
		"actions/checkout@v4":          true,
		"actions/setup-go@v5":          true,
		"actions/upload-artifact@v4":   true,
		"actions/download-artifact@v4": true,
	}
	for jobID := range jobs {
		if !expectedJobs[jobID] {
			t.Fatalf("workflow contains unexpected job %q", jobID)
		}
		lowerJobID := strings.ToLower(jobID)
		for _, forbidden := range []string{"publish", "publisher", "publication", "deploy", "deployment", "upload", "create-tag"} {
			if strings.Contains(lowerJobID, forbidden) {
				t.Fatalf("workflow contains forbidden job %q", jobID)
			}
		}
		job := nestedYAMLSection(t, jobsSection, jobID, 2)
		jobPermissions := optionalNestedYAMLMapping(t, job, "permissions", 4, 6)
		if len(jobPermissions) > 0 && (len(jobPermissions) != 1 || jobPermissions["contents"] != "read") {
			t.Fatalf("job %s permissions = %#v, want inherited or contents: read only", jobID, jobPermissions)
		}
		for _, action := range yamlUsesValues(job) {
			if !allowedActions[action] {
				t.Fatalf("job %s contains non-rehearsal action %q", jobID, action)
			}
		}
	}

	bindJob := nestedYAMLSection(t, jobsSection, "bind-release-inputs", 2)
	shaBinding := `test "$SOURCE_COMMIT" = "$WORKFLOW_SHA"`
	for _, want := range []string{
		"WORKFLOW_SHA: ${{ github.sha }}",
		shaBinding,
	} {
		if !strings.Contains(bindJob, want) {
			t.Fatalf("input binding job missing %q", want)
		}
	}
	if strings.Index(bindJob, shaBinding) > strings.Index(bindJob, "actions/checkout") {
		t.Fatal("source commit must be bound to github.sha before checkout")
	}

	buildJob := nestedYAMLSection(t, jobsSection, "build-native-candidates", 2)
	for _, want := range []string{
		"needs: bind-release-inputs",
		"- os: ubuntu-latest\n            target_label: linux-x86_64\n            binary_suffix: \"\"\n            expected_goos: linux\n            expected_goarch: amd64",
		"- os: macos-15\n            target_label: macos-aarch64\n            binary_suffix: \"\"\n            expected_goos: darwin\n            expected_goarch: arm64",
		"- os: windows-latest\n            target_label: windows-x86_64\n            binary_suffix: \".exe\"\n            expected_goos: windows\n            expected_goarch: amd64",
		`-X github.com/uesugitorachiyo/ao-blueprint/internal/cli.version=${VERSION}`,
		`version_readback=$("$artifact_dir/$binary" --version)`,
		`test "$version_readback" = "$VERSION"`,
		"pack inspect --pack examples/blueprints/valid/ao-blueprint-self --json",
		"help-smoke.txt",
		"version-smoke.json",
		"provider-free-smoke.json",
		`"help":{"command":"--help","status":"passed"`,
		`"version":{"binary":"%s","command":"--version","release_version_readback":"%s","source_commit":"%s","status":"passed"`,
		`"provider_free":{"command":"pack inspect --pack examples/blueprints/valid/ao-blueprint-self --json","status":"passed"`,
		`cd "$artifact_dir"`,
		`sha256sum "$binary" > SHA256SUMS`,
	} {
		if !strings.Contains(buildJob, want) {
			t.Fatalf("native candidate job missing %q", want)
		}
	}

	planJob := nestedYAMLSection(t, jobsSection, "assemble-promotion-plan", 2)
	for _, want := range []string{
		"actions/download-artifact",
		".github/scripts/verify-release-candidates.py",
		"promotion-plan.json",
		"promotion-plan.sha256",
	} {
		if !strings.Contains(planJob, want) {
			t.Fatalf("promotion plan job missing %q", want)
		}
	}

	for _, forbidden := range []string{
		"contents: write",
		"packages: write",
		"id-to" + "ken: write",
		"actions: write",
		"gh release",
		"actions/create-release",
		"softprops/action-gh-release",
		"git tag ",
		"npm publish",
		"cargo publish",
		"docker push",
		"goreleaser",
		"kubectl apply",
		"helm upgrade",
		"deployment:",
		"environment:",
	} {
		if strings.Contains(workflow, forbidden) {
			t.Fatalf("release rehearsal workflow must not include %q", forbidden)
		}
	}
}

func TestReleaseCandidateVersionFlagReadsInjectedVersion(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "ao-blueprint")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	build := exec.Command(
		"go",
		"build",
		"-trimpath",
		"-ldflags",
		"-X github.com/uesugitorachiyo/ao-blueprint/internal/cli.version="+rehearsalVersion,
		"-o",
		binary,
		"./cmd/blueprint",
	)
	build.Dir = rootDir(t)
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build versioned candidate: %v\n%s", err, output)
	}
	output, err := exec.Command(binary, "--version").CombinedOutput()
	if err != nil {
		t.Fatalf("run candidate --version: %v\n%s", err, output)
	}
	if strings.TrimSpace(string(output)) != rehearsalVersion {
		t.Fatalf("candidate version = %q, want %q", output, rehearsalVersion)
	}
}

func TestReleaseCandidateVerifierAcceptsExactInventory(t *testing.T) {
	candidates := writeReleaseCandidateFixture(t)
	plan, checksum, output, err := runReleaseCandidateVerifier(t, candidates)
	if err != nil {
		t.Fatalf("verifier rejected valid inventory: %v\n%s", err, output)
	}
	if _, err := os.Stat(plan); err != nil {
		t.Fatalf("promotion plan missing: %v", err)
	}
	if _, err := os.Stat(checksum); err != nil {
		t.Fatalf("promotion plan checksum missing: %v", err)
	}
}

func TestReleaseCandidateVerifierRejectsInvalidInventory(t *testing.T) {
	tests := []struct {
		name       string
		mutate     func(t *testing.T, root string)
		wantOutput string
	}{
		{
			name: "missing",
			mutate: func(t *testing.T, root string) {
				t.Helper()
				if err := os.Remove(filepath.Join(root, "windows-x86_64", "candidate.json")); err != nil {
					t.Fatal(err)
				}
			},
			wantOutput: "missing candidates",
		},
		{
			name: "duplicate",
			mutate: func(t *testing.T, root string) {
				t.Helper()
				source := filepath.Join(root, "linux-x86_64")
				target := filepath.Join(root, "duplicate-linux")
				copyReleaseCandidateDir(t, source, target)
			},
			wantOutput: "duplicate candidates",
		},
		{
			name: "target platform substitution",
			mutate: func(t *testing.T, root string) {
				t.Helper()
				path := filepath.Join(root, "macos-aarch64", "candidate.json")
				candidate := readReleaseCandidateJSON(t, path)
				candidate["goos"] = "linux"
				candidate["goarch"] = "amd64"
				writeReleaseCandidateJSON(t, path, candidate)
			},
			wantOutput: "target platform substitution",
		},
		{
			name: "stale",
			mutate: func(t *testing.T, root string) {
				t.Helper()
				path := filepath.Join(root, "macos-aarch64", "candidate.json")
				candidate := readReleaseCandidateJSON(t, path)
				candidate["source_commit"] = strings.Repeat("c", 40)
				writeReleaseCandidateJSON(t, path, candidate)
			},
			wantOutput: "stale candidate binding",
		},
		{
			name: "altered manifest digest",
			mutate: func(t *testing.T, root string) {
				t.Helper()
				path := filepath.Join(root, "linux-x86_64", "candidate.json")
				candidate := readReleaseCandidateJSON(t, path)
				candidate["approved_manifest_digest"] = "sha256:" + strings.Repeat("d", 64)
				writeReleaseCandidateJSON(t, path, candidate)
			},
			wantOutput: "stale candidate binding",
		},
		{
			name: "wrong version",
			mutate: func(t *testing.T, root string) {
				t.Helper()
				path := filepath.Join(root, "windows-x86_64", "candidate.json")
				candidate := readReleaseCandidateJSON(t, path)
				candidate["version"] = "v9.9.9"
				writeReleaseCandidateJSON(t, path, candidate)
			},
			wantOutput: "stale candidate binding",
		},
		{
			name: "malformed JSON",
			mutate: func(t *testing.T, root string) {
				t.Helper()
				path := filepath.Join(root, "macos-aarch64", "candidate.json")
				if err := os.WriteFile(path, []byte("{not-json\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			wantOutput: "invalid candidate JSON",
		},
		{
			name: "substituted",
			mutate: func(t *testing.T, root string) {
				t.Helper()
				path := filepath.Join(root, "windows-x86_64", "ao-blueprint.exe")
				if err := os.WriteFile(path, []byte("substituted"), 0o755); err != nil {
					t.Fatal(err)
				}
			},
			wantOutput: "substituted candidate binary",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			candidates := writeReleaseCandidateFixture(t)
			tt.mutate(t, candidates)
			_, _, output, err := runReleaseCandidateVerifier(t, candidates)
			if err == nil {
				t.Fatalf("verifier accepted %s inventory", tt.name)
			}
			if !strings.Contains(output, tt.wantOutput) {
				t.Fatalf("verifier output = %q, want %q", output, tt.wantOutput)
			}
		})
	}
}

func readReleaseRehearsalFile(t *testing.T, elements ...string) string {
	t.Helper()
	path := append([]string{rootDir(t)}, elements...)
	body, err := os.ReadFile(filepath.Join(path...))
	if err != nil {
		t.Fatal(err)
	}
	return strings.ReplaceAll(string(body), "\r\n", "\n")
}

func TestReleaseRehearsalYAMLHelpersAcceptCRLF(t *testing.T) {
	workflow := "on:\r\n  workflow_dispatch:\r\njobs:\r\n  test:\r\n    runs-on: windows-latest\r\n"
	triggers := directYAMLMapping(t, topLevelYAMLSection(t, workflow, "on"), 2)
	if _, ok := triggers["workflow_dispatch"]; !ok {
		t.Fatalf("CRLF workflow triggers = %#v", triggers)
	}
}

func topLevelYAMLSection(t *testing.T, document string, key string) string {
	t.Helper()
	return yamlSection(t, document, key, 0)
}

func nestedYAMLSection(t *testing.T, document string, key string, indent int) string {
	t.Helper()
	return yamlSection(t, document, key, indent)
}

func yamlSection(t *testing.T, document string, key string, indent int) string {
	t.Helper()
	document = strings.ReplaceAll(document, "\r\n", "\n")
	lines := strings.Split(document, "\n")
	prefix := strings.Repeat(" ", indent) + key + ":"
	start := -1
	for index, line := range lines {
		if line == prefix || strings.HasPrefix(line, prefix+" ") {
			start = index
			break
		}
	}
	if start < 0 {
		t.Fatalf("YAML section %q not found", key)
	}
	end := len(lines)
	for index := start + 1; index < len(lines); index++ {
		line := lines[index]
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		if leadingSpaces(line) <= indent {
			end = index
			break
		}
	}
	return strings.Join(lines[start:end], "\n")
}

func optionalNestedYAMLMapping(t *testing.T, document string, key string, sectionIndent int, mappingIndent int) map[string]string {
	t.Helper()
	prefix := strings.Repeat(" ", sectionIndent) + key + ":"
	for _, line := range strings.Split(document, "\n") {
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		inline := strings.TrimSpace(strings.TrimPrefix(line, prefix))
		if inline != "" {
			return map[string]string{"$inline": inline}
		}
		return directYAMLMapping(t, yamlSection(t, document, key, sectionIndent), mappingIndent)
	}
	return nil
}

func directYAMLMapping(t *testing.T, section string, indent int) map[string]string {
	t.Helper()
	result := map[string]string{}
	for _, line := range strings.Split(section, "\n") {
		if leadingSpaces(line) != indent || strings.TrimSpace(line) == "" {
			continue
		}
		key, value, found := strings.Cut(strings.TrimSpace(line), ":")
		if !found {
			t.Fatalf("invalid YAML mapping line %q", line)
		}
		result[key] = strings.TrimSpace(value)
	}
	return result
}

func yamlUsesValues(section string) []string {
	var values []string
	for _, line := range strings.Split(section, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "uses:") {
			values = append(values, strings.TrimSpace(strings.TrimPrefix(trimmed, "uses:")))
		}
	}
	return values
}

func leadingSpaces(value string) int {
	return len(value) - len(strings.TrimLeft(value, " "))
}

func writeReleaseCandidateFixture(t *testing.T) string {
	t.Helper()
	root := filepath.Join(t.TempDir(), "candidates")
	platforms := map[string]struct {
		goos   string
		goarch string
	}{
		"linux-x86_64":   {goos: "linux", goarch: "amd64"},
		"macos-aarch64":  {goos: "darwin", goarch: "arm64"},
		"windows-x86_64": {goos: "windows", goarch: "amd64"},
	}
	for _, target := range []string{"linux-x86_64", "macos-aarch64", "windows-x86_64"} {
		dir := filepath.Join(root, target)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		binary := "ao-blueprint"
		if target == "windows-x86_64" {
			binary += ".exe"
		}
		binaryBody := []byte("candidate-" + target)
		digest := fmt.Sprintf("%x", sha256.Sum256(binaryBody))
		checksumSeparator := "  "
		if target == "windows-x86_64" {
			// GNU sha256sum uses the binary-mode marker on Windows.
			checksumSeparator = " *"
		}
		files := map[string][]byte{
			binary:                     binaryBody,
			"SHA256SUMS":               []byte(digest + checksumSeparator + binary + "\n"),
			"LICENSE":                  []byte("license\n"),
			"NOTICE":                   []byte("notice\n"),
			"help-smoke.txt":           []byte("AO Blueprint\n"),
			"provider-free-smoke.json": []byte("{\"schema\":\"ao.blueprint.pack-inspection.v0.1\",\"status\":\"ready\"}\n"),
			"version-smoke.json":       []byte("{\"binary\":\"" + binary + "\",\"command\":\"--version\",\"module_path\":\"github.com/uesugitorachiyo/ao-blueprint\",\"release_version_readback\":\"" + rehearsalVersion + "\",\"source_commit\":\"" + rehearsalSource + "\",\"status\":\"passed\"}\n"),
		}
		for name, body := range files {
			mode := os.FileMode(0o644)
			if name == binary {
				mode = 0o755
			}
			if err := os.WriteFile(filepath.Join(dir, name), body, mode); err != nil {
				t.Fatal(err)
			}
		}
		candidate := map[string]any{
			"approved_manifest_digest": rehearsalManifest,
			"binary":                   binary,
			"binary_sha256":            digest,
			"dry_run":                  true,
			"goarch":                   platforms[target].goarch,
			"goos":                     platforms[target].goos,
			"publication": map[string]any{
				"public_release_attempted":   false,
				"public_upload_attempted":    false,
				"release_creation_attempted": false,
				"tag_creation_attempted":     false,
			},
			"repository":     "ao-blueprint",
			"schema_version": "ao.blueprint.release-rehearsal-candidate.v0.1",
			"smoke": map[string]any{
				"help":          map[string]any{"status": "passed"},
				"provider_free": map[string]any{"status": "passed"},
				"version": map[string]any{
					"binary":                   binary,
					"command":                  "--version",
					"release_version_readback": rehearsalVersion,
					"source_commit":            rehearsalSource,
					"status":                   "passed",
				},
			},
			"source_commit": rehearsalSource,
			"tag":           rehearsalVersion,
			"target":        target,
			"version":       rehearsalVersion,
		}
		writeReleaseCandidateJSON(t, filepath.Join(dir, "candidate.json"), candidate)
	}
	return root
}

func runReleaseCandidateVerifier(t *testing.T, candidates string) (string, string, string, error) {
	t.Helper()
	outDir := t.TempDir()
	plan := filepath.Join(outDir, "promotion-plan.json")
	checksum := filepath.Join(outDir, "promotion-plan.sha256")
	script := filepath.Join(rootDir(t), ".github", "scripts", "verify-release-candidates.py")
	cmd := exec.Command("python3", script, "--candidates", candidates, "--plan", plan, "--checksum", checksum)
	cmd.Env = append(os.Environ(),
		"VERSION="+rehearsalVersion,
		"TAG="+rehearsalVersion,
		"SOURCE_COMMIT="+rehearsalSource,
		"APPROVED_MANIFEST_DIGEST="+rehearsalManifest,
		"GITHUB_RUN_ATTEMPT=1",
		"GITHUB_RUN_ID=123",
	)
	output, err := cmd.CombinedOutput()
	return plan, checksum, string(output), err
}

func readReleaseCandidateJSON(t *testing.T, path string) map[string]any {
	t.Helper()
	var candidate map[string]any
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(body, &candidate); err != nil {
		t.Fatal(err)
	}
	return candidate
}

func writeReleaseCandidateJSON(t *testing.T, path string, candidate map[string]any) {
	t.Helper()
	body, err := json.Marshal(candidate)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(body, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

func copyReleaseCandidateDir(t *testing.T, source string, target string) {
	t.Helper()
	entries, err := os.ReadDir(source)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	for _, entry := range entries {
		body, err := os.ReadFile(filepath.Join(source, entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(target, entry.Name()), body, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}
