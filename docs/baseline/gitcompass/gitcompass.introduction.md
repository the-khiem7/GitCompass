---
baseline_schema: "2.0"
pack: "gitcompass"
document: "introduction"
status: "active"
updated: "2026-09-02"
code_ref: "unknown"
---

# GitCompass baseline

## Scope

GitCompass is a planned Windows-first, local-first desktop utility for multi-account Git identity, authentication context, and repository routing. It maps reusable Profiles to repositories with deterministic rules, materializes intent through native Git configuration, and checks the result through Identity Guard. The product configures Git rather than replacing it, so ordinary `git commit`, `git pull`, and `git push` must keep working while the app is closed.

## Current truth

This workspace contains no implementation, repository metadata, test suite, executable configuration, or canonical documentation other than the absorbed proposal. Every product, architecture, UX, stack, and roadmap statement in this pack is planned or decided design, not verified runtime behavior. The intended implementation is Go, Wails 3, Vue 3, TypeScript, and Vite; system Git is the source of truth. Windows is Tier 1; Linux and macOS are future work.

The lossless source artifact is [PROPOSAL.md](sources/PROPOSAL.md). It is retained inside this pack to make the pack safe to use after the root-level proposal is deleted. Its SHA-256 is `A455E5A76E9AF495FBC018817BC3CC36BF927395493A91A69C19D8F65C4D686C` and its raw bytes matched the root source on 2026-09-02.

## Product outcome

The user should create distinct Git Profiles, assign them through exact-repository, remote, folder, and default rules, then inspect and safely repair mismatches between the expected Profile and Git's effective context. Commit identity and authentication context are separate: a healthy repository needs both the expected identity and a coherent HTTPS or SSH route. HTTPS through Git Credential Manager or native Git credential helpers is primary; SSH is a supported alternative; GitCompass must not own raw credentials or private-key content.

## Constraints and boundaries

- Native Git configuration and secure credential stores remain authoritative; app persistence stores intent and metadata, not secrets.
- GitCompass owns only isolated generated fragments and its minimal global include. Existing user configuration and repository-local overrides remain user-owned unless a user explicitly approves a minimal Fix.
- Routing precedence is Exact Repository, Remote, Folder, then Default. Resolution and health must be explainable.
- Multi-remote resolution follows Git's effective fetch and push targets rather than assuming `origin`; incompatible remote-to-Profile matches without an Exact Repository rule block materialization.
- Conditional includes are the default materialization mechanism. Repository-local writes are permitted only as minimal, reversible, explicitly approved exceptions.
- The product is a focused identity/context manager, not a hosting platform, full Git GUI, password manager, secret vault, custom Git implementation, proxy daemon, authentication server, or mandatory public CLI.
- The app must prioritize Windows realities: drive letters, case-insensitive paths, Git for Windows, Git Credential Manager, Windows Credential Manager, OpenSSH, junctions, long/Unicode paths, file locks, and path normalization.
- From the first Go code change, use the repository's `use-modern-go` skill and its Modern Go Guidelines CLI. From the first Wails v3 code, configuration, binding, runtime, lifecycle, CLI, build, or packaging change, use both the repository's `wails3` and `wails:wails` skills; keep all Wails APIs and tooling strictly v3 and pin the selected v3 alpha version.

## Source coverage

The archived source contains the full background, problem cases, target user, auth model, profile data model, routing examples, CLI decision, configuration strategy, guard and repair examples, health states, UI concepts, architecture, persistence, security principles, platform considerations, testing, failure modes, roadmap, success criteria, non-goals, and future possibilities. The active documents below own the reusable current state; the source artifact is the complete historical-detail fallback.

| Document | Owns |
|---|---|
| `gitcompass.roadmap.md` | Delivery phases, verification gates, risks, next action |
| `gitcompass.hallucination.md` | Decisions, open questions, rejected approaches, reasoning |
| `gitcompass.sourcecode.md` | Planned architecture, models, flows, config topology |
| `gitcompass.useguide.md` | Planned user and operator contracts |
