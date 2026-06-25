# AO Blueprint Production Readiness

AO Blueprint defines production readiness as an executable gate, not an
informal estimate.

## Required Evidence

- `go test ./...` passes.
- `go vet ./...` passes.
- `blueprint lint --path .` reports zero findings.
- The valid blueprint fixture audits to `score=100` and `status=ready`.
- The valid blueprint fixture authorizes with `status=ready`.
- The invalid fixture fails authorization with `status=blocked`.
- Contract schemas, examples, and generated gate artifacts parse as JSON.
- CI runs the gate on Ubuntu, macOS, and Windows.

## Local Command

```bash
./scripts/production-readiness.sh
```

## Windows Command

```powershell
.\scripts\production-readiness.ps1
```
