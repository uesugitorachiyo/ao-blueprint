# AO Blueprint

AO Blueprint turns an objective and interview answers into a build-ready requirements package. It records assumptions and constraints, checks whether the work is specific enough to proceed, compiles traceability and implementation guidance, and emits bounded build authorization. Use it before implementation when the objective, scope, acceptance criteria, or risk boundaries still need to be made explicit.

## How it fits in AO

- **Primary responsibility:** Requirements clarification and bounded build authorization.
- **Inputs:** Operator objectives, interview answers, constraints, and existing project documentation.
- **Outputs:** Requirements, assumptions, traceability, implementation and quality guidance, blueprint packs, and authorization packets.
- **Upstream:** AO Mission or an operator starting directly from Blueprint.
- **Downstream:** AO Atlas for workgraph and context compilation, or AO Foundry for already bounded work.

See the
[AO Architecture guide](https://github.com/uesugitorachiyo/ao-architecture)
and the
[AO Blueprint component page](https://github.com/uesugitorachiyo/ao-architecture/blob/main/components/ao-blueprint.md)
for the cross-repository flow.

## Commands

```bash
go run ./cmd/blueprint --help
go run ./cmd/blueprint lint --path .
go run ./cmd/blueprint readiness audit --pack examples/blueprints/valid/ao-blueprint-self --out tmp/readiness.json
go run ./cmd/blueprint sdd emit --pack examples/blueprints/valid/ao-blueprint-self --out tmp/sdd-plan.json
go run ./cmd/blueprint authorize --pack examples/blueprints/valid/ao-blueprint-self --out tmp/build-authorization.json
go run ./cmd/blueprint pack inspect --pack examples/blueprints/valid/ao-blueprint-self --json
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
