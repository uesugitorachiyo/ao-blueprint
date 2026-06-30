# AO Blueprint SDD Planner Prompt

Create an `ao2.sdd-plan.v1` plan for AO Blueprint.

AO Blueprint is the requirements interview, blueprint compiler, and
build-authorization gate before downstream AO work. For oversized,
mutation-class, or long-running work, AO Atlas is the required compiler between
Blueprint and Foundry. The plan must include
contracts, interview state, readiness scoring, safety scanning, SDD emission,
authorization, CLI commands, tests, CI, public-safety rules, and production
readiness gates.

The output must preserve these constraints:

- Go CLI first;
- Ubuntu, macOS, and Windows portability;
- fail-closed authorization;
- no target product implementation execution;
- public-safe durable artifacts;
- 100/100 readiness gate before downstream AO work.
