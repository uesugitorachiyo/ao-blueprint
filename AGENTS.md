# AO Blueprint Agent Instructions

## Status And Role

AO Blueprint is the active requirements-clarification and bounded build-authorization component. It owns interview state, requirements, assumptions, traceability, sufficiency scoring, SDD guidance, blueprint packs, and authorization packets.

Blueprint authorization means only that the declared requirements are sufficiently bounded for the named build scope. It does not schedule or execute work, mutate repositories, approve runtime side effects, satisfy downstream policy, or authorize a release or deployment.

## Sources Of Truth

- [docs/sdd/AO-BLUEPRINT-PRD.md](docs/sdd/AO-BLUEPRINT-PRD.md), [docs/sdd/AO-BLUEPRINT-ARCHITECTURE.md](docs/sdd/AO-BLUEPRINT-ARCHITECTURE.md), and [docs/sdd/AO-BLUEPRINT-CONTRACTS.md](docs/sdd/AO-BLUEPRINT-CONTRACTS.md) define scope, flow, and contracts.
- [docs/sdd/AO-BLUEPRINT-INTERVIEW.md](docs/sdd/AO-BLUEPRINT-INTERVIEW.md) and [docs/sdd/AO-BLUEPRINT-READINESS.md](docs/sdd/AO-BLUEPRINT-READINESS.md) own interview and authorization stop conditions.
- `docs/contracts/` owns JSON schemas; `internal/blueprint/`, `internal/cli/`, and their tests are authoritative for implemented behavior.
- [docs/design/PRODUCTION-READINESS.md](docs/design/PRODUCTION-READINESS.md), `scripts/production-readiness.sh`, and [`.github/workflows/ci.yml`](.github/workflows/ci.yml) define the broad verification gate.

## Ownership And Boundaries

- Preserve explicit objective, scope, non-goals, assumptions, acceptance criteria, risks, approvals, and traceability. Missing or contradictory requirements must block authorization.
- Keep authorization bound to the exact pack and declared build scope. Do not present readiness, SDD output, inspection, or an authorization packet as execution, policy, release, or deployment authority.
- Treat `examples/blueprints/valid/` and `examples/blueprints/invalid/` as contract fixtures. Change fixtures with consuming tests and never weaken an invalid case to force readiness.
- Keep generated material under ignored `tmp/` or `target/`; do not hand-edit generated output into durable evidence or rewrite historical examples to support a current claim.
- Do not add secrets, provider credentials, private paths, account identifiers, live-provider behavior, or implicit repository mutation.
- Release rehearsal, release, deployment, publication, credentialed operation, and direct-main changes require separate explicit authority.

## Working Method

- Change the smallest owned requirements or contract surface. Preserve fail-closed scoring, deterministic output, path safety, exact provenance, and downstream non-authority flags.
- Add negative tests for missing approvals, incomplete traceability, unsafe paths, malformed contracts, or over-broad authorization.
- Update this file in the same pull request when durable commands, architecture, ownership, or authority boundaries change.

## Verification

- Requirements and readiness logic: `go test ./internal/blueprint -count=1`.
- CLI behavior: `go test ./internal/cli -count=1`.
- Format relevant Go source with `gofmt -d cmd internal`; run `go test ./... -count=1`, `go vet ./...`, and `go build -o tmp/blueprint ./cmd/blueprint`.
- Run `./scripts/production-readiness.sh` as the full local gate. Release-rehearsal workflows are conditional on separate release authority and are not part of an instruction-only change.
- For instruction changes run `python3 ../ao-architecture/scripts/verify_agent_instruction_layout.py --workspace-root .. --repository ao-blueprint`. Always run `git diff --check`.

## Evidence And Completion

- Record the source head, command exits, and input/output digests where a contract binds them. Report skipped, unavailable, or failed checks explicitly.
- A readiness score or build authorization is bounded evidence, not permission to implement outside its exact scope or to release, deploy, publish, or mutate another repository.
- Completion requires focused and broad gates, green pull-request CI, clean synchronized `main`, and task-branch cleanup.
