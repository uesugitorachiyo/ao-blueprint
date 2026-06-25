# AO Blueprint Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build AO Blueprint as a Go CLI that interviews, compiles, audits, emits SDD plans, and authorizes downstream AO builds only when readiness is 100/100.

**Architecture:** The CLI delegates to `internal/blueprint` for contract loading, safety scanning, sufficiency scoring, SDD emission, and build authorization. Durable fixtures and schemas live under `examples/` and `docs/contracts/`; generated outputs go under `tmp/`.

**Tech Stack:** Go 1.22, standard library only, JSON fixtures, JSON Schema documents, shell and PowerShell readiness scripts.

## Global Constraints

- Authorization fails closed unless score is 100, user approval is true, and no blockers exist.
- Durable public artifacts must not contain local machine paths, raw tokens, private keys, or secret values.
- The default downstream path is AO Blueprint -> AO Foundry -> AO Forge -> AO2.
- AO Blueprint must not execute target product implementation work.
- The product gate is `./scripts/production-readiness.sh`.

---

### Task 1: Contract Spine And Red Tests

**Files:**
- Create: `go.mod`
- Create: `internal/blueprint/blueprint_test.go`
- Create: `internal/cli/cli_test.go`
- Create: `examples/blueprints/valid/ao-blueprint-self/*`
- Create: `examples/blueprints/invalid/missing-approval/*`

**Interfaces:**
- Produces tests for `AuditPack`, `AuthorizePack`, `EmitSDDPlan`, `InspectPack`, `LintPath`, and `cli.Run`.

- [ ] Write tests that require valid readiness to score 100 and invalid authorization to block.
- [ ] Run `go test ./...` and verify undefined implementation failures.

### Task 2: Core Blueprint Package

**Files:**
- Create: `internal/blueprint/blueprint.go`

**Interfaces:**
- Produces `AuditPack(pack string) (SufficiencyAudit, error)`.
- Produces `AuthorizePack(pack string) (BuildAuthorization, error)`.
- Produces `EmitSDDPlan(pack string, out string) error`.
- Produces `InspectPack(pack string) (PackInspection, error)`.
- Produces `LintPath(path string) (LintReport, error)`.

- [ ] Implement public-safety scanning.
- [ ] Implement required file checks.
- [ ] Implement JSON parsing checks.
- [ ] Implement 100-point readiness scoring.
- [ ] Implement authorization with SHA-256 digests.
- [ ] Run `go test ./internal/blueprint`.

### Task 3: CLI Package

**Files:**
- Create: `cmd/blueprint/main.go`
- Create: `internal/cli/cli.go`

**Interfaces:**
- Produces `Run(args []string, stdout io.Writer, stderr io.Writer) error`.

- [ ] Implement help and command routing.
- [ ] Implement `lint`, `readiness audit`, `sdd emit`, `authorize`, `pack inspect`.
- [ ] Implement `interview start`, `interview next`, `interview answer`, and `interview status`.
- [ ] Run `go test ./internal/cli`.

### Task 4: Production Gate

**Files:**
- Create: `scripts/production-readiness.sh`
- Create: `scripts/production-readiness.ps1`
- Create: `.github/workflows/ci.yml`

**Interfaces:**
- Produces the executable local gate and hosted CI definition.

- [ ] Run tests and vet.
- [ ] Run lint and public-safety scan.
- [ ] Run readiness audit, SDD emit, authorization, and pack inspect.
- [ ] Verify invalid authorization blocks.

### Task 5: Final Verification

**Files:**
- Modify: `README.md`

**Interfaces:**
- Produces a public-facing quickstart and production-readiness evidence.

- [ ] Run `gofmt -w`.
- [ ] Run `go test ./...`.
- [ ] Run `go vet ./...`.
- [ ] Run `./scripts/production-readiness.sh`.
- [ ] Run a final public-safety scan.
- [ ] Commit when clean.
