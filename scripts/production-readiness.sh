#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT="${BLUEPRINT_READINESS_DIR:-$ROOT/tmp/production-readiness}"
VALID_PACK="$ROOT/examples/blueprints/valid/ao-blueprint-self"
INVALID_PACK="$ROOT/examples/blueprints/invalid/missing-approval"

rm -rf "$OUT"
mkdir -p "$OUT"

cd "$ROOT"

go test ./...
go vet ./...
go build -o "$OUT/blueprint" ./cmd/blueprint

go run ./cmd/blueprint --help > "$OUT/help.txt"
go run ./cmd/blueprint lint --path .
go run ./cmd/blueprint readiness audit --pack "$VALID_PACK" --out "$OUT/readiness.json"
go run ./cmd/blueprint sdd emit --pack "$VALID_PACK" --out "$OUT/sdd-plan.json"
go run ./cmd/blueprint authorize --pack "$VALID_PACK" --out "$OUT/build-authorization.json"
go run ./cmd/blueprint pack inspect --pack "$VALID_PACK" --json > "$OUT/pack-inspect.json"
go run ./cmd/blueprint interview start --idea "Build a governed requirements gate" --out "$OUT/session.json"
go run ./cmd/blueprint interview answer --session "$OUT/session.json" --question-id q-objective --answer "Success is an approved 100 point blueprint authorization gate" --out "$OUT/session.json"
go run ./cmd/blueprint compile --session "$OUT/session.json" --out-dir "$OUT/draft-pack"

if go run ./cmd/blueprint authorize --pack "$INVALID_PACK" --out "$OUT/blocked-authorization.json" > "$OUT/blocked-stdout.txt" 2> "$OUT/blocked-stderr.txt"; then
  echo "invalid fixture unexpectedly authorized" >&2
  exit 1
fi

if go run ./cmd/blueprint authorize --pack "$OUT/draft-pack" --out "$OUT/draft-authorization.json" > "$OUT/draft-stdout.txt" 2> "$OUT/draft-stderr.txt"; then
  echo "draft pack unexpectedly authorized without approval" >&2
  exit 1
fi

python3 - "$OUT" "$ROOT" <<'PY'
import json
import pathlib
import sys

out = pathlib.Path(sys.argv[1])
root = pathlib.Path(sys.argv[2])

for path in list((root / "docs" / "contracts").glob("*.json")) + list((root / "examples").glob("**/*.json")) + list(out.glob("*.json")):
    with path.open("r", encoding="utf-8") as handle:
        json.load(handle)

readiness = json.loads((out / "readiness.json").read_text(encoding="utf-8"))
if readiness.get("score") != 100 or readiness.get("status") != "ready":
    raise SystemExit(f"readiness gate failed: {readiness}")

auth = json.loads((out / "build-authorization.json").read_text(encoding="utf-8"))
if auth.get("status") != "ready" or auth.get("score") != 100:
    raise SystemExit(f"authorization gate failed: {auth}")

blocked = json.loads((out / "blocked-authorization.json").read_text(encoding="utf-8"))
if blocked.get("status") != "blocked":
    raise SystemExit(f"invalid fixture did not block: {blocked}")

draft = json.loads((out / "draft-authorization.json").read_text(encoding="utf-8"))
if draft.get("status") != "blocked" or draft.get("approved_by_user") is not False:
    raise SystemExit(f"draft pack did not remain blocked before approval: {draft}")

print("AO Blueprint production readiness: 100/100 status=ready")
PY
