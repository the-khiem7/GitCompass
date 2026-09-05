---
baseline_schema: "2.0"
pack: "gitcompass"
document: "roadmap"
status: "active"
updated: "2026-09-06"
code_ref: "ee2d97d9429f0e24169fc707b002b1e7d4d496a8"
---

# GitCompass delivery roadmap

## Status

P1-P3 have focused implementation evidence; P4 is in progress. All known design questions required through P6 are decided, so no recorded Open requires an implementation pause before the final phase; phase transitions remain gated by implementation and acceptance evidence. The stages remain ordered by product dependency; passing builds or static checks do not prove Git, credentials, native configuration, or Windows integrations work in practice.

| Phase | Deliverable | Depends on | Status |
|---|---|---|---|
| P1 | Windows Git foundation | - | complete |
| P2 | Profile engine | P1 | complete |
| P3 | Routing engine | P2 | complete |
| P4 | Managed configuration materialization | P1-P3 | in progress |
| P5 | Identity Guard and approved repair | P3-P4 | planned |
| P6 | Desktop product | P2-P5 | planned |
| P7 | Windows hardening | P1-P6 | planned |
| P8 | Future-platform preparation | P7 | planned |

## P1 - Windows Git foundation

Deliver Git for Windows discovery, safe Git command execution, repository detection, Git version inspection, effective configuration with origin and scope, remotes, effective commit identity, and credential-helper inspection. Before every Go code edit from this phase onward, invoke the repository `use-modern-go` skill and run its Modern Go Guidelines CLI for the target file. Acceptance requires temporary repositories demonstrating those reads against a real Git for Windows installation, including unavailable Git and malformed configuration failures.

Completed with `internal/git`: Git discovery and repository inspection run against the installed Git for Windows 2.53.0.windows.1. Its focused tests create temporary repositories and verify identity, credential-helper, remote, origin/scope configuration, and non-repository behavior. `go test ./...`, `wails3 build`, and `wails3 doctor` passed; the Wails build produced `bin/gitcompass.exe`, while the desktop UI itself remains a later P6 acceptance gate.

## P2 - Profile engine

Deliver Profile CRUD, identity, HTTPS context and credential-helper references, optional SSH support, signing configuration, metadata, and validation. Acceptance requires validation tests for incomplete and invalid profile data without handling raw credentials.

Completed with `internal/profile` and `internal/persistence`: Profile CRUD persists to SQLite through the sole writer, uses stable UUIDs and UTC timestamps, and stores helper/key references but no credential payloads. Focused tests cover incomplete and invalid Profiles plus SQLite create, update, list, get, and delete. `go test ./...` and `wails3 build` passed.

## P3 - Routing engine

Deliver exact-repository, remote, folder, and default rules; deterministic precedence; conflict detection; and a human-readable resolution explanation. Multi-remote routing must resolve Git's effective fetch/pull and default push targets, evaluate every remote, and block materialization when remote rules select different Profiles without an Exact Repository rule. Acceptance requires unit tests covering compatible remotes, conflicting remotes with and without an Exact rule, different fetch and push targets, overlapping rules, moved repositories, changed remotes, normalized drive/case paths, and no-match behavior.

Completed with `internal/routing`: exact repository, remote, folder, and Default matching follows the required precedence with case-insensitive Windows-path normalization. Remote matching accepts HTTPS and SSH forms, reports conflicting Profile matches as non-materializable, and returns a winning rule plus explanation. Unit tests cover precedence, path normalization, conflicting multi-remote selection, and Default fallback. `go test ./...` passed; effective fetch/pull and push-target extraction remains to be wired from the Git adapter before P5.

## P4 - Managed configuration materialization

Deliver a GitCompass-owned root, profile fragments, global include integration, ordered Default/Folder/Remote/Exact conditional materialization, capability checks, stale Exact-rule handling, approved local-write exceptions, ownership metadata, idempotent apply, and clean removal. Acceptance requires before/after real Git config inspection proving precedence through config origin and scope, a moved Folder rule re-evaluates, a moved Exact rule becomes stale, an unsupported condition blocks instead of writing locally, unrelated user configuration is unchanged, repeated apply is stable, and uninstall removes only GitCompass-owned configuration.

Initial implementation in `internal/materialize` generates owned Profile and root fragments in Default/Folder/Remote/Exact group order and inserts one idempotent global include without overwriting an unrelated setting in its focused test. The current implementation directly interpolates Profile values, rewrites files through a fixed `.tmp` path, checks the global include with a substring search, has no production caller, and does not sort rules by specificity within a kind. Focused package tests for materialization, routing, and Git inspection pass, but they do not prove effective Git conditional-include behavior.

Controlled temporary-repository probes against Git for Windows 2.53.0 established the current P4 checkpoint:

- `hasconfig:remote.*.url` is supported, but a substring such as `github.com/acme` does not match; explicit HTTPS and SCP-style SSH glob families do match.
- A Folder `gitdir/i` pattern without a trailing slash does not match descendants; the trailing slash form does.
- An Exact `gitdir/i` condition compiled from the worktree path does not match; the absolute Git directory returned by `git rev-parse --absolute-git-dir` does.
- `worktree/i` did not match in the tested Git version and cannot be used as an assumed fallback.
- Exact fixed-value removal can delete duplicate managed `include.path` values while preserving unrelated includes.
- Correct Default, broad Folder, narrow Folder, Remote, and Exact file order yields Exact as the effective last value, with origin and scope visible through Git inspection.

P4 remains blocked from acceptance by six implementation areas: correct deterministic rule compilation and escaping; semantic capability probes; durable ownership metadata and a generation-based apply lifecycle; clean removal and recovery; stale Exact and linked-worktree handling; and approved local-exception planning. The selected architecture is recorded in D12 (desired-state materialization lifecycle) in `gitcompass.hallucination.md`; no P4 design Open remains.

## P5 - Identity Guard and approved repair

Deliver expected-versus-effective comparison, Healthy/Warning/Mismatch/Broken/Unknown classification, origin explanation, HTTPS and SSH diagnostics, unexpected-local-override detection, and a minimal user-approved Fix flow. Acceptance requires integration tests for a wrong `user.email`, missing helper, malformed include, conflicting remotes, missing SSH key, and a repair that changes only the displayed approved key.

Authentication assurance uses the evidence levels selected in D13 (evidence-labeled authentication verification) in `gitcompass.hallucination.md`: Configured, Bound, Operation verified, and Unknown. Background checks never retrieve credentials or prompt. P5 has no remaining design Open, but still depends on completed P4 materialization and must implement provider capability probes and the evidence-to-health mapping.

## P6 - Desktop product

Deliver Home, Profiles, Rules, Repositories, Settings, health views, repository detail, and Fix approval flows. Before every Wails v3 code, configuration, binding generation, runtime API, lifecycle, CLI, build, or packaging action from this phase onward, use both the repository `wails3` and `wails:wails` skills; detect and retain Wails v3-only APIs, pin the selected alpha version, and do not mix v2 artifacts or commands. Acceptance requires a desktop user to create two Profiles, route repositories, see why a Profile won, inspect config origin, and approve a mismatch repair without a public CLI.

Repository inventory follows D14 (opt-in bounded repository discovery) in `gitcompass.hallucination.md`: explicit registration and user-selected roots only, cancellable in-process scans with fixed default budgets, no v1 watcher or service, local-only metadata, and non-destructive aging. P6 has no remaining design Open, but implementation must prove cancellation, budget, inaccessible-root, missing, archive, and rediscovery behavior.

## P7 - Windows hardening

Deliver installer, migrations, redacted diagnostic export, recovery from malformed config, Git Credential Manager edge-case handling, Windows path edge cases, security review, and file-lock recovery. Acceptance requires installed-Windows testing with drive letters, case-insensitive paths, junctions, long/Unicode paths, locked config, Windows Credential Manager, and Git Credential Manager behavior.

## P8 - Future-platform preparation

Research Linux and macOS only after Windows quality is satisfactory. Validate abstractions and packaging/credential-store mappings for Linux credential helpers and keyrings, and macOS Keychain, filesystem behavior, SSH, app signing, and webview differences. This phase must not dilute the Windows-first implementation.

## Cross-phase verification

- Unit tests: routing precedence and conflicts, Profile validation, materialization planning, health classification, path normalization, and ownership rules.
- Integration tests: temporary repositories with global config, managed include, conditional includes, repository-local overrides, multiple remotes, credential helpers, optional SSH, Guard, and repair.
- Safety tests: no unrelated config deletion, no secret exposure, no silent user-owned rewrite, idempotent managed configuration, and clean managed-config removal.
- Product success: two Profiles, native HTTPS credentials, optional SSH, deterministic routing, normal Git with the app closed, explainable resolution and health, config-origin inspection, and safe mismatch repair.

## Risks

- Git's effective behavior is composed from system, global, conditional, and local sources, so a model-only implementation can misdiagnose the real source.
- Credential helper behavior is provider and Windows-installation dependent; presence of configuration is not proof that authentication succeeds.
- Provider-specific credential behavior can still make an operation's account context unknown even after remote routing is deterministic.
- File locks, junctions, path casing, and malformed include files can make apparently simple writes unsafe.

## Next action

Finish P4 in this order:

1. Correct the compiler for Folder, Exact, and Remote conditions, Profile-value escaping, conflict validation, and deterministic broad-to-narrow ordering.
2. Add isolated semantic capability probes keyed by Git executable and version; block unsupported behavior without a silent local fallback.
3. Add minimal ownership metadata and generation-based apply with desired-state hashing and native `git config` mutations.
4. Add exact fixed-value removal, manifest-limited cleanup, interrupted-apply recovery, and post-change verification.
5. Add stale Exact detection, explicit rebind, linked-worktree classification, and an approval-fingerprinted local managed-include exception.
6. Run the real Git acceptance matrix for precedence, origins, moves, repeated apply, removal, malformed or locked config, and Windows path variants before starting P5.
