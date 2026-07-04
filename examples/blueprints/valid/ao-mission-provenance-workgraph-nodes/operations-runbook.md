# Operations Runbook

```sh
go run ./cmd/blueprint lint --path examples/blueprints/valid/ao-mission-provenance-workgraph-nodes
go run ./cmd/blueprint readiness audit --pack examples/blueprints/valid/ao-mission-provenance-workgraph-nodes --out tmp/ao-mission-provenance-workgraph-nodes-readiness.json
go run ./cmd/blueprint sdd emit --pack examples/blueprints/valid/ao-mission-provenance-workgraph-nodes --out tmp/ao-mission-provenance-workgraph-nodes-sdd.json
go run ./cmd/blueprint authorize --pack examples/blueprints/valid/ao-mission-provenance-workgraph-nodes --out tmp/ao-mission-provenance-workgraph-nodes-authorization.json
go run ./cmd/blueprint pack inspect --pack examples/blueprints/valid/ao-mission-provenance-workgraph-nodes --json
```

Next step after authorization is AO Atlas first. AO Foundry waits for Atlas
workgraph/context/candidate compilation and imports only one safe node.
