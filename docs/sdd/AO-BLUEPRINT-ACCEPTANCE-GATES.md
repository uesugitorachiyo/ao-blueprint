# AO Blueprint Acceptance Gates

## Local Gate

```bash
go test ./...
go vet ./...
go run ./cmd/blueprint --help
go run ./cmd/blueprint lint --path .
go run ./cmd/blueprint readiness audit --pack examples/blueprints/valid/ao-blueprint-self --out tmp/readiness.json
go run ./cmd/blueprint sdd emit --pack examples/blueprints/valid/ao-blueprint-self --out tmp/sdd-plan.json
go run ./cmd/blueprint authorize --pack examples/blueprints/valid/ao-blueprint-self --out tmp/build-authorization.json
go run ./cmd/blueprint pack inspect --pack examples/blueprints/valid/ao-blueprint-self --json
```

## Product Gate

```bash
./scripts/production-readiness.sh
```

Expected result:

- tests pass;
- vet passes;
- lint passes;
- valid fixture readiness is `score=100` and `status=ready`;
- invalid fixture authorization is blocked;
- public-safety scan has zero findings;
- generated JSON artifacts parse.
