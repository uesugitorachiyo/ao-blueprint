# Quality Profile

## Code Quality

AO Blueprint requires downstream implementation to use simple, readable,
typed code with clear names, bounded functions, repository-local evidence, and
no speculative abstractions.

## TDD And Evals

Each implementation slice must name the failing-first test, deterministic
fixture, scorecard, or reviewer gate that proves the expected behavior before
code changes begin.

## Verification Loop

Build, vet or type checks, lint, tests, schema validation, public-safety scan,
and production-readiness commands must be recorded before downstream AO work
claims completion.

## Security Review

Security-sensitive work must cover secret handling, input validation,
authorization, dependency posture, logs, errors, and public artifact safety.
Findings become AO Sentinel or AO Covenant packets instead of private notes.
