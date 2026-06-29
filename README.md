# AO Blueprint

AO Blueprint is the front-door requirements interview, blueprint compiler, and
build-authorization gate for the AO orchestration framework. It prevents vague
ideas from entering AO Foundry or AO Forge until the user's objective,
constraints, domain model, contracts, tests, operations model, security posture,
and production-readiness exit condition are specific enough to build.

AO Blueprint is intentionally not an implementation runner. It emits a reviewed
blueprint pack and a machine-readable build authorization packet. Downstream AO
automation must refuse to start when authorization is blocked.

Build authorization is not live mutation approval. AO Blueprint can classify
underspecified work, docs-only work, and build-ready work, but a first tiny
docs-only live repository mutation still requires the later exact-scope
Covenant, Foundry, Forge, AO2, Sentinel, Promoter, rollback, Command, and
operator approval chain. Blueprint does not approve patches, create branches,
execute work, call providers, publish, release, or grant broad live mutation
authority.

Every ready blueprint pack must include `implementation-spec.md`, a concrete
pre-SDD build contract with outcome, scope, stack, constraints, and verification
sections. It must also include `quality-profile.md`, which records the
AO-tailored code quality, TDD/eval, verification-loop, and security-review bar
for downstream implementation. This keeps AO Foundry and AO Forge from starting
implementation from a vague interview transcript alone.

## Role In The AO Stack

```text
raw idea
-> AO Blueprint interview and blueprint pack
-> AO Blueprint build authorization packet
-> AO Foundry portfolio scheduling
-> AO Forge governed factory run
-> AO Covenant policy and side-effect gates
-> AO2 bounded local execution
-> AO Arena benchmark comparison
-> AO Crucible adversarial hardening
-> AO Sentinel safety and regression monitoring
-> AO Promoter gated activation
```

## Commands

```bash
go run ./cmd/blueprint --help
go run ./cmd/blueprint lint --path .
go run ./cmd/blueprint readiness audit --pack examples/blueprints/valid/ao-blueprint-self --out tmp/readiness.json
go run ./cmd/blueprint sdd emit --pack examples/blueprints/valid/ao-blueprint-self --out tmp/sdd-plan.json
go run ./cmd/blueprint authorize --pack examples/blueprints/valid/ao-blueprint-self --out tmp/build-authorization.json
go run ./cmd/blueprint pack inspect --pack examples/blueprints/valid/ao-blueprint-self --json
go run ./cmd/blueprint authorize --pack examples/blueprints/valid/bounded-governed-rsi-control-surface-readback --out tmp/bounded-rsi-build-authorization.json
```

## Production-Readiness Gate

```bash
./scripts/production-readiness.sh
```

The gate runs tests, vet, lint, public-safety scan, readiness audit, SDD emit,
authorization, pack inspection, and JSON parsing over durable examples.

## SDD Files

| File | Purpose |
| --- | --- |
| `docs/sdd/AO-BLUEPRINT-PRD.md` | Product scope, users, goals, non-goals, and readiness definition. |
| `docs/sdd/AO-BLUEPRINT-ARCHITECTURE.md` | CLI, packages, data flow, contracts, and AO stack boundaries. |
| `docs/sdd/AO-BLUEPRINT-CONTRACTS.md` | Contract families, required fields, and validation semantics. |
| `docs/sdd/AO-BLUEPRINT-INTERVIEW.md` | Interview state machine, question categories, and stop conditions. |
| `docs/sdd/AO-BLUEPRINT-READINESS.md` | 100/100 sufficiency scoring and build authorization blockers. |
| `docs/sdd/AO-BLUEPRINT-IMPLEMENTATION-SLICES.md` | Implementation slices in dependency order. |
| `docs/sdd/AO-BLUEPRINT-ACCEPTANCE-GATES.md` | Product and public-readiness verification commands. |
| `docs/sdd/AO-BLUEPRINT-SDD-HANDOFF.md` | Handoff prompt for AO Forge, AO Foundry, or Codex. |

## License

AO Blueprint is licensed under `Apache-2.0`. See `LICENSE`.
