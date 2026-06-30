# AO Blueprint Architecture

## Boundaries

AO Blueprint owns interview state, blueprint compilation, sufficiency scoring,
SDD plan emission, public-safety scanning, and build authorization. It does not
own implementation execution, policy authority, benchmark scoring, adversarial
hardening, monitoring, promotion, branch creation, patch application, provider
calls, publishing, releases, or live repository mutation approval.

## Packages

| Package | Responsibility |
| --- | --- |
| `cmd/blueprint` | Console entrypoint. |
| `internal/cli` | Argument parsing, command routing, and user-facing output. |
| `internal/blueprint` | Contract loading, safety scan, readiness audit, SDD emit, and authorization logic. |

## Data Flow

```text
idea -> interview session -> answers -> blueprint pack
-> lint -> readiness audit -> user approval
-> SDD plan -> build authorization -> downstream AO handoff
```

For oversized, mutation-class, and long-running work, the downstream handoff is
AO Blueprint -> AO Atlas -> AO Foundry. Blueprint emits the pack and build
authorization; Atlas is the mandatory compiler that imports the pack, verifies
authorization scope/freshness, builds the workgraph and context packs, records
candidate selection, and emits Foundry import material. Foundry must not accept
a direct Blueprint handoff for these classes.

For docs-only live mutation preparation, AO Blueprint may mark the requirement
as build-ready or blocked for clarification. That build authorization only lets
Atlas compile the work for downstream consideration; exact live mutation
permission remains with the later Covenant ticket, Foundry approval gate, Forge
guard, AO2 patch packet, Sentinel/Promoter boundaries, rollback rehearsal,
Command readback, and operator approval.

## Failure Model

All mutating or downstream-enabling commands fail closed. Missing files, invalid
JSON, unsafe content, unresolved blockers, score below 100, missing user
approval, or digest mismatch keeps authorization blocked.
