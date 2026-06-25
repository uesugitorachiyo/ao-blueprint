# AO Blueprint Architecture

## Boundaries

AO Blueprint owns interview state, blueprint compilation, sufficiency scoring,
SDD plan emission, public-safety scanning, and build authorization. It does not
own implementation execution, policy authority, benchmark scoring, adversarial
hardening, monitoring, or promotion.

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

## Failure Model

All mutating or downstream-enabling commands fail closed. Missing files, invalid
JSON, unsafe content, unresolved blockers, score below 100, missing user
approval, or digest mismatch keeps authorization blocked.
