#!/usr/bin/env python3

import argparse
import hashlib
import json
import os
from pathlib import Path


EXPECTED_TARGETS = {
    "linux-x86_64": {"goos": "linux", "goarch": "amd64"},
    "macos-aarch64": {"goos": "darwin", "goarch": "arm64"},
    "windows-x86_64": {"goos": "windows", "goarch": "amd64"},
}
PUBLICATION_FIELDS = (
    "tag_creation_attempted",
    "release_creation_attempted",
    "public_upload_attempted",
    "public_release_attempted",
)


def fail(message):
    raise SystemExit(message)


def parse_args():
    parser = argparse.ArgumentParser()
    parser.add_argument("--candidates", required=True, type=Path)
    parser.add_argument("--plan", required=True, type=Path)
    parser.add_argument("--checksum", required=True, type=Path)
    return parser.parse_args()


def expected_bindings():
    return {
        "version": os.environ["VERSION"],
        "tag": os.environ["TAG"],
        "source_commit": os.environ["SOURCE_COMMIT"],
        "approved_manifest_digest": os.environ["APPROVED_MANIFEST_DIGEST"],
    }


def read_json(path):
    try:
        return json.loads(path.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError) as error:
        fail(f"invalid candidate JSON: {path}: {error}")


def verify_candidate(candidate_path, expected):
    candidate = read_json(candidate_path)
    target = candidate.get("target")
    if any(candidate.get(field) != value for field, value in expected.items()):
        fail(f"stale candidate binding: {candidate_path}")
    expected_platform = EXPECTED_TARGETS.get(target)
    if expected_platform is None:
        fail(f"unexpected candidate target: {candidate_path}: {target}")
    if any(candidate.get(field) != value for field, value in expected_platform.items()):
        fail(f"target platform substitution: {candidate_path}")
    if candidate.get("dry_run") is not True:
        fail(f"stale candidate status: {candidate_path}")

    publication = candidate.get("publication", {})
    if any(publication.get(field) is not False for field in PUBLICATION_FIELDS):
        fail(f"unexpected publication attempt: {candidate_path}")

    smoke = candidate.get("smoke", {})
    if any(smoke.get(name, {}).get("status") != "passed" for name in ("help", "version", "provider_free")):
        fail(f"stale candidate smoke status: {candidate_path}")

    directory = candidate_path.parent
    binary = candidate.get("binary")
    binary_path = directory / str(binary)
    sums_path = directory / "SHA256SUMS"
    if not binary or not binary_path.is_file() or not sums_path.is_file():
        fail(f"missing candidate inventory: {candidate_path}")

    checksum_fields = sums_path.read_text(encoding="utf-8").strip().split()
    if len(checksum_fields) != 2 or checksum_fields[1] != binary or Path(checksum_fields[1]).name != checksum_fields[1]:
        fail(f"non-portable SHA256SUMS entry: {candidate_path}")
    actual_sha256 = hashlib.sha256(binary_path.read_bytes()).hexdigest()
    if candidate.get("binary_sha256") != actual_sha256 or checksum_fields[0] != actual_sha256:
        fail(f"substituted candidate binary: {candidate_path}")

    help_smoke = directory / "help-smoke.txt"
    version_smoke_path = directory / "version-smoke.json"
    provider_smoke_path = directory / "provider-free-smoke.json"
    if not help_smoke.read_text(encoding="utf-8").strip():
        fail(f"missing help smoke: {candidate_path}")
    version_smoke = read_json(version_smoke_path)
    if (
        version_smoke.get("binary") != binary
        or version_smoke.get("command") != "--version"
        or version_smoke.get("release_version_readback") != expected["version"]
        or version_smoke.get("source_commit") != expected["source_commit"]
        or version_smoke.get("module_path") != "github.com/uesugitorachiyo/ao-blueprint"
        or version_smoke.get("status") != "passed"
    ):
        fail(f"stale version smoke: {candidate_path}")
    candidate_version_smoke = smoke.get("version", {})
    if (
        candidate_version_smoke.get("binary") != binary
        or candidate_version_smoke.get("command") != "--version"
        or candidate_version_smoke.get("release_version_readback") != expected["version"]
        or candidate_version_smoke.get("source_commit") != expected["source_commit"]
    ):
        fail(f"stale candidate version smoke: {candidate_path}")
    provider_smoke = read_json(provider_smoke_path)
    if (
        provider_smoke.get("schema") != "ao.blueprint.pack-inspection.v0.1"
        or provider_smoke.get("status") != "ready"
    ):
        fail(f"stale provider-free smoke: {candidate_path}")

    inventory = sorted(path.name for path in directory.iterdir() if path.is_file())
    expected_inventory = sorted(
        [
            binary,
            "SHA256SUMS",
            "LICENSE",
            "NOTICE",
            "candidate.json",
            "help-smoke.txt",
            "provider-free-smoke.json",
            "version-smoke.json",
        ]
    )
    if inventory != expected_inventory:
        fail(f"unexpected candidate inventory: {candidate_path}: {inventory}")

    return {
        "artifact": directory.name,
        "binary": binary,
        "binary_sha256": actual_sha256,
        "goarch": candidate["goarch"],
        "goos": candidate["goos"],
        "smoke": smoke,
        "target": target,
    }


def main():
    args = parse_args()
    expected = expected_bindings()
    candidate_paths = sorted(args.candidates.rglob("candidate.json"))
    targets = []
    candidates = []
    for candidate_path in candidate_paths:
        candidate = read_json(candidate_path)
        targets.append(candidate.get("target"))
        candidates.append(verify_candidate(candidate_path, expected))

    expected_targets = set(EXPECTED_TARGETS)
    found_targets = set(targets)
    if expected_targets - found_targets:
        fail(f"missing candidates: {sorted(expected_targets - found_targets)}")
    if len(targets) != len(found_targets):
        fail(f"duplicate candidates: {targets}")
    if found_targets - expected_targets:
        fail(f"unexpected candidates: {sorted(found_targets - expected_targets)}")

    plan = {
        "approved_manifest_digest": expected["approved_manifest_digest"],
        "candidates": sorted(candidates, key=lambda candidate: candidate["target"]),
        "dry_run": True,
        "immutable": True,
        "inventory_targets": sorted(found_targets),
        "publication": {field: False for field in PUBLICATION_FIELDS},
        "repository": "ao-blueprint",
        "schema_version": "ao.blueprint.specialist-promotion-plan.v0.1",
        "source_commit": expected["source_commit"],
        "tag": expected["tag"],
        "version": expected["version"],
        "workflow_run": {
            "attempt": os.environ.get("GITHUB_RUN_ATTEMPT", ""),
            "id": os.environ.get("GITHUB_RUN_ID", ""),
        },
    }
    plan_bytes = (json.dumps(plan, sort_keys=True, separators=(",", ":")) + "\n").encode()
    args.plan.write_bytes(plan_bytes)
    args.checksum.write_text(
        hashlib.sha256(plan_bytes).hexdigest() + "  " + args.plan.name + "\n",
        encoding="utf-8",
    )


if __name__ == "__main__":
    main()
