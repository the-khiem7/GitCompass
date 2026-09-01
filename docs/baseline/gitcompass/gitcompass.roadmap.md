---
baseline_schema: "2.0"
pack: "gitcompass"
document: "roadmap"
status: "active"
updated: "2026-09-02"
code_ref: "unknown"
---

# GitCompass delivery roadmap

## Status

No phase has implementation evidence. The stages below are planned and ordered by product dependency; passing builds or static checks later must not be recorded as proof that Git, credentials, native configuration, or Windows integrations work in practice.

| Phase | Deliverable | Depends on | Status |
|---|---|---|---|
| P1 | Windows Git foundation | - | complete |
| P2 | Profile engine | P1 | planned |
| P3 | Routing engine | P2 | planned |
| P4 | Managed configuration materialization | P1-P3 | planned |
| P5 | Identity Guard and approved repair | P3-P4 | planned |
| P6 | Desktop product | P2-P5 | planned |
| P7 | Windows hardening | P1-P6 | planned |
| P8 | Future-platform preparation | P7 | planned |

## P1 - Windows Git foundation

Deliver Git for Windows discovery, safe Git command execution, repository detection, Git version inspection, effective configuration with origin and scope, remotes, effective commit identity, and credential-helper inspection. Before every Go code edit from this phase onward, invoke the repository `use-modern-go` skill and run its Modern Go Guidelines CLI for the target file. Acceptance requires temporary repositories demonstrating those reads against a real Git for Windows installation, including unavailable Git and malformed configuration failures.

Completed with `internal/git`: Git discovery and repository inspection run against the installed Git for Windows 2.53.0.windows.1. Its focused tests create temporary repositories and verify identity, credential-helper, remote, origin/scope configuration, and non-repository behavior. `go test ./...`, `wails3 build`, and `wails3 doctor` passed; the Wails build produced `bin/gitcompass.exe`, while the desktop UI itself remains a later P6 acceptance gate.

## P2 - Profile engine

Deliver Profile CRUD, identity, HTTPS context and credential-helper references, optional SSH support, signing configuration, metadata, and validation. Acceptance requires validation tests for incomplete and invalid profile data without handling raw credentials.

## P3 - Routing engine

Deliver exact-repository, remote, folder, and default rules; deterministic precedence; conflict detection; and a human-readable resolution explanation. Multi-remote routing must resolve Git's effective fetch/pull and default push targets, evaluate every remote, and block materialization when remote rules select different Profiles without an Exact Repository rule. Acceptance requires unit tests covering compatible remotes, conflicting remotes with and without an Exact rule, different fetch and push targets, overlapping rules, moved repositories, changed remotes, normalized drive/case paths, and no-match behavior.

## P4 - Managed configuration materialization

Deliver a GitCompass-owned root, profile fragments, global include integration, ordered Default/Folder/Remote/Exact conditional materialization, capability checks, stale Exact-rule handling, approved local-write exceptions, ownership metadata, idempotent apply, and clean removal. Acceptance requires before/after real Git config inspection proving precedence through config origin and scope, a moved Folder rule re-evaluates, a moved Exact rule becomes stale, an unsupported condition blocks instead of writing locally, unrelated user configuration is unchanged, repeated apply is stable, and uninstall removes only GitCompass-owned configuration.

## P5 - Identity Guard and approved repair

Deliver expected-versus-effective comparison, Healthy/Warning/Mismatch/Broken/Unknown classification, origin explanation, HTTPS and SSH diagnostics, unexpected-local-override detection, and a minimal user-approved Fix flow. Acceptance requires integration tests for a wrong `user.email`, missing helper, malformed include, conflicting remotes, missing SSH key, and a repair that changes only the displayed approved key.

## P6 - Desktop product

Deliver Home, Profiles, Rules, Repositories, Settings, health views, repository detail, and Fix approval flows. Before every Wails v3 code, configuration, binding generation, runtime API, lifecycle, CLI, build, or packaging action from this phase onward, use both the repository `wails3` and `wails:wails` skills; detect and retain Wails v3-only APIs, pin the selected alpha version, and do not mix v2 artifacts or commands. Acceptance requires a desktop user to create two Profiles, route repositories, see why a Profile won, inspect config origin, and approve a mismatch repair without a public CLI.

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

Implement P2 Profile CRUD and validation using the decided SQLite persistence boundary. Preserve the no-secrets rule, and add tests for incomplete and invalid Profile data before beginning P3 routing.
