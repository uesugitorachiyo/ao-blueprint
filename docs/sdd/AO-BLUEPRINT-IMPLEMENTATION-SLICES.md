# AO Blueprint Implementation Slices

## Slice 1: Contract Spine

- Go module, CLI entrypoint, help text, license, notice, README, and gitignore.
- Contract schemas for session, question, answer, requirement, assumption,
  decision, risk, traceability, sufficiency audit, pack, SDD plan, and build
  authorization.
- Valid and invalid fixture packs.

Acceptance gate: `go test ./...` reaches the expected red phase before
implementation and later passes.

## Slice 2: Readiness And Safety

- Public-safety scanner.
- JSON lint for durable examples and schemas.
- Blueprint pack inspection.
- Sufficiency audit with 100-point scoring and hard blockers.
- Required implementation spec and quality profile artifacts for downstream
  AO build readiness.

Acceptance gate: valid fixture scores 100/100; invalid fixture is blocked.

## Slice 3: Authorization And SDD Handoff

- SDD plan emit.
- Build authorization with content digests.
- Downstream handoff paths, with AO Atlas required before AO Foundry for
  oversized, mutation-class, or long-running work.
- CLI commands for readiness, SDD, authorize, and inspect.

Acceptance gate: authorization is ready for the valid pack and non-zero for the
invalid pack.

## Slice 4: Interview State

- Start, next, answer, and status commands.
- Deterministic question category coverage.
- Session fixture coverage.

Acceptance gate: session starts from an idea, asks the next missing category,
records an answer, and reports status.

## Slice 5: Production Gate

- Shell and PowerShell readiness scripts.
- Cross-platform CI.
- Clean-clone public-readiness docs.

Acceptance gate: `./scripts/production-readiness.sh` passes locally.
