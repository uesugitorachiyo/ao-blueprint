# AO Blueprint Contracts

## Contract Families

- `ao.blueprint.session.v0.1`
- `ao.blueprint.question.v0.1`
- `ao.blueprint.answer.v0.1`
- `ao.blueprint.requirement.v0.1`
- `ao.blueprint.assumption.v0.1`
- `ao.blueprint.decision.v0.1`
- `ao.blueprint.risk.v0.1`
- `ao.blueprint.traceability-matrix.v0.1`
- `ao.blueprint.sufficiency-audit.v0.1`
- `ao.blueprint.pack.v0.1`
- `ao.blueprint.sdd-plan.v0.1`
- `ao.blueprint.build-authorization.v0.1`

## Required Semantics

Every contract includes a schema version, stable ID, public-safety class, and
human-readable validation errors. Pack-level contracts must be validated from a
repository-relative path and must not depend on local machine state.

Every ready blueprint pack must include `implementation-spec.md`. The
implementation spec is the pre-SDD build contract and must cover outcome,
scope, stack, constraints, and verification. Readiness blocks when it is
missing because downstream AO automation must not infer build detail from a
vague interview transcript.

Every ready blueprint pack must also include `quality-profile.md`. The quality
profile adapts Claude Code-style code-quality, TDD, eval, verification-loop,
and security-review patterns into AO-owned gates without copying prompt or
tooling-specific instructions into the public artifact.

## Build Authorization

Authorization requires:

- `status=ready`;
- `score=100`;
- `approved_by_user=true`;
- no blocking assumptions;
- matching digests for requirements, traceability, and SDD plan;
- production-readiness exit condition present;
- next allowed action targets AO Atlas, AO Foundry, or AO Forge.

The authorization contract is a build-readiness contract only. It can point a
ready bounded requirement toward downstream AO implementation, but for
oversized, mutation-class, or long-running work it must be consumed by AO Atlas
before AO Foundry. Atlas must emit the ready Blueprint import, workgraph,
context packs, candidate-selection record, and Foundry import material; missing,
blocked, stale, or out-of-scope authorization remains a blocked Blueprint
request, not a ready Foundry workgraph. Authorization must not claim patch
approval, live execution permission, provider authority, release authority, or
broad live mutation authority.
