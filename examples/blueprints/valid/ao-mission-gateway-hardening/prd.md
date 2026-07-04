# Product Requirements

AO Mission must leave downstream operators and tools with concrete next
commands and digest-bound readbacks for gateway continuation surfaces. Gateway
readbacks can be fresh, stale, or unknown; they must not be treated as always
fresh. A2A fixture-server readback must expose only local Agent Card and JSON-RPC
fixture paths. Timeline compaction must be bindable by Atlas, Foundry, and
Command as readback-only evidence.

Success means AO Atlas receives a complete workgraph/context/candidate plan
before AO Foundry sees a ready node. Foundry imports only the first safe node and
does not treat generated handoff files as completion.
