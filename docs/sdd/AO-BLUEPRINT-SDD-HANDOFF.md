# AO Blueprint SDD Handoff

Use this handoff when asking AO Forge, AO Foundry, AO2, or Codex to implement
AO Blueprint.

```text
Build AO Blueprint at production-readiness 100/100.

Follow docs/sdd/AO-BLUEPRINT-PRD.md, AO-BLUEPRINT-ARCHITECTURE.md,
AO-BLUEPRINT-CONTRACTS.md, AO-BLUEPRINT-INTERVIEW.md,
AO-BLUEPRINT-READINESS.md, AO-BLUEPRINT-IMPLEMENTATION-SLICES.md, and
AO-BLUEPRINT-ACCEPTANCE-GATES.md.

Implement slice by slice. Preserve the rule that AO Blueprint interviews,
audits, compiles, and authorizes, but does not execute target product builds.

Stop only when ./scripts/production-readiness.sh passes, the valid fixture
audits to score=100 status=ready, the invalid fixture blocks authorization, and
no public-safety findings exist.
```
