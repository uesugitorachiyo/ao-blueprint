# AO Blueprint Readiness

## Score Matrix

| Category | Points |
| --- | ---: |
| Objective and success metrics | 10 |
| Scope and non-goals | 10 |
| Domain and workflows | 12 |
| Interfaces and contracts | 10 |
| Data and integrations | 8 |
| Security, privacy, and public safety | 10 |
| Tests and evaluation | 12 |
| Operations and release | 8 |
| Traceability | 10 |
| User approval and build handoff | 10 |

## Blockers

Any of these blocks authorization regardless of score:

- missing approval;
- score below 100;
- blocking assumptions;
- unsafe public artifact finding;
- invalid JSON artifact;
- missing SDD plan;
- missing implementation spec or quality profile;
- missing traceability matrix;
- missing production-readiness exit condition;
- next action outside the authorized AO path, including AO Atlas before AO
  Foundry for oversized, mutation-class, or long-running work.

Build authorization must also remain scoped to requirement readiness. It must
not be interpreted as approval for live repository mutation, provider calls,
branch creation, patch application, release, publication, or fully
unsupervised complex mutation. A docs-only or low-risk live mutation candidate
still needs Atlas import/readback plus the downstream exact-scope approval and
gate chain before any PR rehearsal can execute.
