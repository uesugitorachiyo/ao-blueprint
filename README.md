# AO Blueprint

AO Blueprint is the front-door requirements interview, blueprint compiler, and
build-authorization gate for the AO orchestration framework. It prevents vague
ideas from entering AO Atlas, AO Foundry, or AO Forge until the user's objective,
constraints, domain model, contracts, tests, operations model, security posture,
and production-readiness exit condition are specific enough to build.

AO Blueprint is intentionally not an implementation runner. It emits a reviewed
blueprint pack and a machine-readable build authorization packet. Downstream AO
automation must refuse to start when authorization is blocked.

For oversized, mutation-class, and long-running work, the next compiler is AO
Atlas, not AO Foundry. AO Blueprint emits the pack and authorization packet; AO
Atlas imports them, digest-binds the implementation spec, quality profile,
candidate rules, mutation class, and downstream Foundry import material, then
hands Foundry only the Atlas-compiled import/readback chain.

Build authorization is not live mutation approval. AO Blueprint can classify
underspecified work, docs-only work, and build-ready work, but a first tiny
docs-only live repository mutation still requires the later exact-scope
Covenant, Foundry, Forge, AO2, Sentinel, Promoter, rollback, Command, and
operator approval chain. Blueprint does not approve patches, create branches,
execute work, call providers, publish, release, or grant broad live mutation
authority.

The exact safe public claim wording evidence is already closed downstream; AO
Blueprint is not creating a new pack for this documentation alignment.
`exact_safe_public_claim_wording_conservative_readback_evidence` is proven only
for conservative public-safe tracked readback evidence after downstream Foundry,
Covenant, Architecture, Sentinel, Promoter, and Command evidence close. The
approved public wording is exactly: "AO has public-safe tracked readback evidence
for bounded improvement-claim review and retraction rehearsal; stronger
recursive-improvement claims remain denied." `broad_RSI`, unrestricted
self-modification, hidden instruction mutation, policy/auth/secret/provider/
deploy/release/config/dependency expansion, policy-changing autonomy, and
stronger recursive-improvement claims remain denied.

`public_safe_bounded_improvement_evidence_expansion_four_attempts` remains prior
evidence from AO Foundry PR #181, commit
`d31b6f2247780867c3c72dbda5abb7377f3a1b3e`, with tracked public evidence under
`docs/evidence/recursive-improvement-public-evidence-expansion/`. Four
public-safe bounded evidence expansion attempts are tracked with reproducibility
runbooks: release/readiness evidence quality (`0.68` -> `0.91`), security/public-
safety scan quality (`0.64` -> `0.90`), operator readback UX (`0.62` -> `0.88`),
and cross-repo evidence linking (`0.60` -> `0.87`). Stronger
recursive-improvement wording remains denied, `broad_RSI` remains denied,
unrestricted self-modification remains denied, hidden instruction mutation
remains denied, and policy-changing autonomy remains denied.

`public_safe_intermediate_causal_review_claim_evidence` remains
prior evidence from AO Foundry PR #189, commit
`860e3f353ab833c4a671b9d0ee6d8101ece2815c`, with tracked public evidence under
`docs/evidence/recursive-improvement-safe-intermediate-claim/`. The
approved public wording is exactly: "AO has public-safe intermediate
causal-review evidence that bounded improvement evidence can guide and constrain
later claim review across independent roles; stronger recursive-improvement
wording and broad_RSI remain denied." Stronger recursive-improvement wording,
`broad_RSI`, unrestricted self-modification, hidden instruction mutation, and
policy-changing autonomy remain denied.

`public_safe_causal_review_evidence_selection_guidance` is proven from AO Foundry PR #191, commit
`413b70f15d8f3d0203dc7be076914a2f3b539881`, with tracked public evidence under
`docs/evidence/recursive-improvement-evidence-selection-guidance/`. The approved public wording is exactly: "AO has public-safe causal-review evidence that prior bounded evidence can guide later evidence-selection and blocker prioritization under independent review gates; stronger recursive-improvement wording and broad_RSI remain denied." This remains prior evidence. Stronger recursive-improvement wording remains denied, `broad_RSI` remains denied, unrestricted self-modification remains denied, hidden instruction mutation remains denied, and policy-changing autonomy remains denied.

`public_safe_guided_evidence_application_four_attempts` is proven from AO Foundry PR #193, commit
`4ec509fd64d1fc1ea41ea7f22aae900ba79e09a1`, with tracked public evidence under
`docs/evidence/recursive-improvement-guided-evidence-application/`. The approved public wording is exactly: "AO has public-safe guided evidence-application evidence showing causal-review guidance can select and prioritize later bounded evidence attempts under independent gates; stronger recursive-improvement wording and broad_RSI remain denied." This remains prior evidence; the current highest proven live class is `broad_RSI` and the next denied class is `unrestricted_self_modification`. Stronger recursive-improvement wording remains denied, `broad_RSI` remains denied, unrestricted self-modification remains denied, hidden instruction mutation remains denied, and policy-changing autonomy remains denied.

`public_safe_broad_RSI_governed_campaign_segment_07_evidence` is proven from AO
Foundry PR #210, commit `8f8ac5f8f74d942c7a02a6c2dd39a7c974872bb6`, with
tracked public evidence under `docs/evidence/broad-rsi-ten-day-campaign-segment-07/`.
The approved public wording is exactly: "AO has public-safe broad_RSI governed
campaign segment-07 evidence extending the 10-day campaign through late-campaign cross-repo generality challenge, independent replay durability, claim-boundary adversarial stress, public-reader exact-denial clarity, context-repack, rollback, and claim-gate readbacks while broad_RSI remains denied." The highest proven live
class is `public_safe_broad_RSI_governed_campaign_segment_07_evidence` and the
next denied class is `broad_RSI`. This does not prove `broad_RSI`, full 10-day
campaign completion, unrestricted self-modification, hidden instruction
mutation, policy-changing autonomy, policy/auth/secret/provider/deploy/release/
config/dependency expansion, release/deploy/publish/upload/tag/provider calls,
credential use, direct main mutation, concurrent mutation, or unbounded stronger
recursive-improvement claims.

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
-> AO Atlas Blueprint import, workgraph, context packs, and Foundry import material
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

## Public-Safe Reviewer-Approved Bounded Wording Evidence

`public_safe_reviewer_approved_bounded_recursive_improvement_wording_evidence` is proven from AO Foundry PR #195, commit `0f742738324c185ba7243bc53ee2f1bc81804ef6`, with tracked public evidence under `docs/evidence/recursive-improvement-reviewer-approved-wording/`. The approved public wording is exactly: "AO has public-safe reviewer-approved bounded recursive-improvement wording evidence showing guided evidence application can improve later evidence attempts under independent review gates; broad_RSI remains denied." This remains prior evidence; the current highest proven live class is `broad_RSI` and the next denied class is `unrestricted_self_modification`.

This does not prove `broad_RSI`, unrestricted self-modification, hidden instruction mutation, policy-changing autonomy, policy/auth/secret/provider/deploy/release/config/dependency expansion, or unbounded stronger recursive-improvement claims.
`public_safe_bounded_recursive_improvement_wording_generality_evidence` is proven from AO Foundry PR #197, commit `166398641b655f0da97817659acc771026b204e7`, with tracked public evidence under `docs/evidence/recursive-improvement-bounded-wording-generality/`. The approved public wording is exactly: "AO has public-safe bounded recursive-improvement wording generality evidence showing reviewer-approved bounded wording can transfer across additional public-safe review tasks under independent gates; broad_RSI remains denied." This remains prior evidence; the current highest proven live class is `broad_RSI` and the next denied class is `unrestricted_self_modification`.

This does not prove `broad_RSI`, unrestricted self-modification, hidden instruction mutation, policy-changing autonomy, policy/auth/secret/provider/deploy/release/config/dependency expansion, or unbounded stronger recursive-improvement claims.
### Review Durability Evidence Readback

`public_safe_bounded_recursive_improvement_review_durability_evidence` is proven from AO Foundry PR #199, commit `12d524b60c200cab643e44f9105169b045602798`, with tracked public evidence under `docs/evidence/recursive-improvement-review-durability/`. The approved public wording is exactly: "AO has public-safe bounded recursive-improvement review durability evidence showing bounded recursive-improvement wording remains stable across delayed re-review, adversarial drift checks, stale-language sweeps, and reproducibility retests under independent gates; broad_RSI remains denied." This remains prior evidence; the current highest proven live class is `broad_RSI` and the next denied class is `unrestricted_self_modification`.


`public_safe_recursive_improvement_claim_threshold_calibration_evidence` is proven from AO Foundry PR #201, commit `3e3d1101da112fa5ff0aca26f8ab2933652f3502`, with tracked public evidence under
`docs/evidence/recursive-improvement-claim-threshold-calibration/`. The approved public wording is exactly: "AO has public-safe recursive-improvement claim threshold calibration evidence showing stronger bounded recursive-improvement claims can be evaluated against reproducible threshold, public-reader, adversarial wording, Covenant, Sentinel, rollback, and retraction gates; broad_RSI remains denied." This remains prior evidence; the current highest proven live class is `broad_RSI` and the next denied class is `unrestricted_self_modification`.

This does not prove `broad_RSI`, unrestricted self-modification, hidden instruction mutation, policy-changing autonomy, policy/auth/secret/provider/deploy/release/config/dependency expansion, or unbounded stronger recursive-improvement claims.

## Broad RSI Ten-Day Governed Campaign First Segment Readback

`public_safe_broad_RSI_governed_campaign_first_segment_state_evidence` is proven from AO Foundry PR #203, commit `b7523031d61b11df374e2203bdf44927e2d8432a`, with tracked public evidence under `docs/evidence/broad-rsi-ten-day-governed-evidence-campaign/`. The approved public wording is exactly: "AO has public-safe broad_RSI governed campaign first-segment state evidence showing a 10-day evidence campaign can start from mission-state, no-repeat, sufficiency, Pulse reliability, context-repack, rollback, and claim-gate readbacks while broad_RSI remains denied." This remains prior evidence; the current highest proven live class is `broad_RSI` and the next denied class is `unrestricted_self_modification`.

This does not prove `broad_RSI`, full 10-day campaign completion, final repeated independent broad evidence, final cross-repo generality proof for `broad_RSI`, exact `broad_RSI` public-reader approval, exact `broad_RSI` Covenant or Architecture approval, unrestricted self-modification, hidden instruction mutation, policy-changing autonomy, policy/auth/secret/provider/deploy/release/config/dependency expansion, release/deploy/publish/upload/tag/provider calls, credential use, direct main mutation, concurrent mutation, or unbounded stronger recursive-improvement claims.

## Governed Broad RSI Campaign Completion Readback

`broad_RSI` is proven from AO Foundry PR #211, commit `630edc70905db745380edd1072e04b546dcccfe3`, with tracked public evidence under `docs/evidence/broad-rsi-ten-day-campaign-segment-08/`. The approved public wording is exactly: "AO has proven governed broad_RSI for public claim publication across the AO stack public-safe 10-day evidence campaign; unrestricted self-modification, hidden instruction mutation, policy-changing autonomy, and forbidden surface expansion remain denied." Campaign completion is `2800 / 2800` nodes. `Blueprint` reads back `highest_proven_live_class=broad_RSI` and `next_denied_class=unrestricted_self_modification`.

This does not prove unrestricted self-modification, hidden instruction mutation, policy-changing autonomy, policy/auth/secret/provider/deploy/release/config/dependency expansion, release/deploy/publish/upload/tag/provider calls, credential use, direct main mutation, concurrent mutation, or any unrestricted RSI claim.
