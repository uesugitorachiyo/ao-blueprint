$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent $PSScriptRoot
$Out = if ($env:BLUEPRINT_READINESS_DIR) { $env:BLUEPRINT_READINESS_DIR } else { Join-Path $Root "tmp/production-readiness" }
$ValidPack = Join-Path $Root "examples/blueprints/valid/ao-blueprint-self"
$InvalidPack = Join-Path $Root "examples/blueprints/invalid/missing-approval"

Remove-Item -Recurse -Force $Out -ErrorAction SilentlyContinue
New-Item -ItemType Directory -Force -Path $Out | Out-Null

Push-Location $Root
try {
  go test ./...
  go vet ./...
  go build -o (Join-Path $Out "blueprint.exe") ./cmd/blueprint

  go run ./cmd/blueprint --help | Out-File -Encoding utf8 (Join-Path $Out "help.txt")
  go run ./cmd/blueprint lint --path .
  go run ./cmd/blueprint readiness audit --pack $ValidPack --out (Join-Path $Out "readiness.json")
  go run ./cmd/blueprint sdd emit --pack $ValidPack --out (Join-Path $Out "sdd-plan.json")
  go run ./cmd/blueprint authorize --pack $ValidPack --out (Join-Path $Out "build-authorization.json")
  go run ./cmd/blueprint pack inspect --pack $ValidPack --json | Out-File -Encoding utf8 (Join-Path $Out "pack-inspect.json")
  go run ./cmd/blueprint interview start --idea "Build a governed requirements gate" --out (Join-Path $Out "session.json")
  go run ./cmd/blueprint interview answer --session (Join-Path $Out "session.json") --question-id q-objective --answer "Success is an approved 100 point blueprint authorization gate" --out (Join-Path $Out "session.json")
  go run ./cmd/blueprint compile --session (Join-Path $Out "session.json") --out-dir (Join-Path $Out "draft-pack")

  $blockedOut = Join-Path $Out "blocked-authorization.json"
  go run ./cmd/blueprint authorize --pack $InvalidPack --out $blockedOut
  $blockedExit = $LASTEXITCODE
  if ($blockedExit -eq 0) {
    throw "invalid fixture unexpectedly authorized"
  }

  $draftOut = Join-Path $Out "draft-authorization.json"
  go run ./cmd/blueprint authorize --pack (Join-Path $Out "draft-pack") --out $draftOut
  $draftExit = $LASTEXITCODE
  if ($draftExit -eq 0) {
    throw "draft pack unexpectedly authorized without approval"
  }

  Get-ChildItem -Recurse -File -Include *.json docs,examples,$Out | ForEach-Object {
    Get-Content -Raw $_.FullName | ConvertFrom-Json | Out-Null
  }

  $readiness = Get-Content -Raw (Join-Path $Out "readiness.json") | ConvertFrom-Json
  if ($readiness.score -ne 100 -or $readiness.status -ne "ready") {
    throw "readiness gate failed"
  }

  $auth = Get-Content -Raw (Join-Path $Out "build-authorization.json") | ConvertFrom-Json
  if ($auth.score -ne 100 -or $auth.status -ne "ready") {
    throw "authorization gate failed"
  }

  $blocked = Get-Content -Raw $blockedOut | ConvertFrom-Json
  if ($blocked.status -ne "blocked") {
    throw "invalid fixture did not block"
  }

  $draft = Get-Content -Raw $draftOut | ConvertFrom-Json
  if ($draft.status -ne "blocked" -or $draft.approved_by_user -ne $false) {
    throw "draft pack did not remain blocked before approval"
  }

  Write-Output "AO Blueprint production readiness: 100/100 status=ready"
} finally {
  Pop-Location
}
