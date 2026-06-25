# AO Blueprint Interview

## State Machine

```text
draft -> interviewing -> compiling -> review -> ready
                               \-> blocked
                               \-> abandoned
```

## Question Categories

- product goal and success metric;
- target users and roles;
- domain entities and lifecycle states;
- workflows and operator paths;
- interfaces and output artifacts;
- integrations, data, and persistence;
- security, privacy, and public safety;
- tests, evals, and adversarial hardening;
- release, rollback, monitoring, and ownership;
- explicit non-goals and deferrals.

## Stop Conditions

The interview may stop only when every category is answered, declared not
applicable with a reason, or blocked by an explicit user-owned decision.

The first implementation uses deterministic category coverage. Later versions
can add model-assisted next-question selection behind the same contracts.
