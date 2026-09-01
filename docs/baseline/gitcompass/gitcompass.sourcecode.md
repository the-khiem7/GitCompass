---
baseline_schema: "2.0"
pack: "gitcompass"
document: "sourcecode"
status: "active"
updated: "2026-09-02"
code_ref: "unknown"
---

# Planned architecture and configuration topology

## System shape

```text
Vue desktop UI
  -> Wails bindings
  -> application services
  -> Profile / Routing / Materialization / Guard engines
  -> Git Engine and local persistence
  -> system Git, Git config files, credential helpers, secure stores
```

Business logic must not live in Vue components. Suggested services are `ProfileService`, `RoutingService`, `RepositoryService`, `MaterializationService`, `IdentityGuardService`, `SettingsService`, and `DiagnosticsService`. The basic Guard flow is UI -> `IdentityGuardService.CheckRepository()` -> Routing Engine + Git Engine -> health result.

## Domain model

A Profile represents identity, HTTPS authentication context, optional SSH configuration, signing configuration, validation state, and metadata. Routing Rules associate a Profile using an exact repository path, a remote host/pattern, a folder prefix, or a default. Repository state records path, remotes, resolved Profile, effective identity, authentication method, health, config origin, and discovery metadata. Persistence also needs managed-config metadata, health cache, UI preferences, and migrations; it must not hold raw credentials.

## Resolution and health

Resolution evaluates normalized repository identity against the precedence `Exact Repository > Remote > Folder > Default`, detects ambiguous/conflicting rules, and returns the winning Profile with an explanation. Identity Guard compares the expected Profile to actual `user.name`, `user.email`, remote context, credential mechanism, SSH route if used, signing configuration, config origin, and unexpected local overrides. Health is `Healthy`, `Warning`, `Mismatch`, `Broken`, or `Unknown`; each state must include the observed facts and cause.

## Native configuration model

The planned layout is:

```text
User-owned Git config
  + [include] path = C:/Users/User/.gitcompass/gitconfig
  -> GitCompass-owned root fragment
  -> GitCompass-owned profile fragments
  -> effective Git configuration
```

The managed directory is expected to contain `C:\\Users\\User\\.gitcompass\\gitconfig` and `profiles\\personal.gitconfig`, `profiles\\company.gitconfig`, and equivalent future fragments. Materialization translates Profile plus Routing intent into those fragments, conditional includes, or narrowly necessary repository-local values. It must be idempotent, minimal, reversible, ownership-aware, and capable of clean removal without deleting unrelated configuration.

## Planned modules

```text
internal/git          Git discovery, commands, config/origin, remotes, identity, helpers
internal/profile      Profile CRUD and validation
internal/routing      Matching, precedence, conflict detection, explanations
internal/materialize  Managed fragments, apply/remove planning, ownership metadata
internal/guard        Expected/effective comparison, health, repair planning
internal/repository   Inspection, inventory, discovery registration
internal/persistence  Local state and migrations
internal/config       Paths, settings, generated-config handling
apps/desktop/backend  Wails application bindings
apps/desktop/frontend Vue 3 user interface
```

## Security boundaries

Git commands must be executed safely and inputs must be validated for Windows paths, long/Unicode names, drive letters, case-insensitivity, junctions, and locks. Diagnostics can include Git version/path, Profile names, remote hosts, routing decisions, config origin, and health, but must redact tokens, passwords, private-key contents, and credential payloads. SSH handling retains only key references/paths and validates missing paths without reading private contents.
