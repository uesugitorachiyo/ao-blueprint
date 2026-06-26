# Operations Runbook

Run:

```bash
go run ./cmd/blueprint readiness audit --pack examples/blueprints/valid/bounded-governed-rsi-control-surface-readback --out tmp/bounded-rsi-readiness.json
go run ./cmd/blueprint authorize --pack examples/blueprints/valid/bounded-governed-rsi-control-surface-readback --out tmp/bounded-rsi-build-authorization.json
```

Downstream operators must reject work that changes scope from bounded readback
improvement to full autonomous RSI publication.
