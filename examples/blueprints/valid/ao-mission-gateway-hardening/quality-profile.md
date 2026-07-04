# Quality Profile

- Routing must be Atlas-first.
- Authorization and readiness audit must agree on `next_allowed_action=ao-atlas`.
- Every generated task must state safe-to-execute and no-authority flags.
- Public docs must avoid claims that gateway readbacks are always fresh.
- Tests must cover valid readbacks and invalid authority-widening fixtures.
- CI/readiness evidence must be fresh before merge.
