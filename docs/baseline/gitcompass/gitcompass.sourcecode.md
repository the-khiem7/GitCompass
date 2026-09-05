---
baseline_schema: "2.0"
pack: "gitcompass"
document: "sourcecode"
status: "active"
updated: "2026-09-06"
code_ref: "ee2d97d9429f0e24169fc707b002b1e7d4d496a8"
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

## Implemented P1 foundation

The repository uses Go module `github.com/TheKhiem7/GitCompass`, Wails `v3.0.0-beta.8`, Vue 3, TypeScript, and Vite. `internal/git` owns safe system-Git discovery and read-only repository inspection. It invokes Git without a shell, rejects non-worktrees, and returns remotes plus effective configuration origin/scope, identity, and credential-helper information. The adapter does not write Git configuration or read credential payloads.

Business logic must not live in Vue components. Suggested services are `ProfileService`, `RoutingService`, `RepositoryService`, `MaterializationService`, `IdentityGuardService`, `SettingsService`, and `DiagnosticsService`. The basic Guard flow is UI -> `IdentityGuardService.CheckRepository()` -> Routing Engine + Git Engine -> health result.

## Implementation governance

Before editing a Go file, use the repository `use-modern-go` skill and run its Modern Go Guidelines CLI against that file. Before creating or changing any Wails v3 application, service, binding, runtime call, lifecycle hook, CLI/build task, or packaging configuration, use both repository skills `wails3` and `wails:wails`. This project is v3-only: pin the chosen Wails v3 alpha version in `go.mod`, use the v3 `wails3` commands and service/binding model, and never introduce v2 imports, `wailsjs` bindings, `wails` CLI commands, or v2 runtime APIs.

## Domain model

A Profile represents identity, HTTPS authentication context, optional SSH configuration, signing configuration, validation state, and metadata. Routing Rules associate a Profile using an exact repository path, a remote host/pattern, a folder prefix, or a default. Repository state records path, remotes, resolved Profile, effective identity, authentication method, health, config origin, and discovery metadata. Persistence uses an application-owned SQLite database for Profiles, rules, known repositories, discovery roots, managed-config metadata, settings, schema state, and migrations; health cache is rebuildable. The database stores intent and metadata only, not raw credentials.

## Persistence boundary

`internal/persistence` is the sole writer to a versioned relational SQLite schema. Core rows use stable identifiers; repository records retain display and normalized Windows paths. Embedded sequential migrations run in transactions with a busy timeout for transient locks, create a consistent backup before upgrade, roll back on failure, and support forward migration only. An older application version must refuse to write a newer schema.

P2 implements the initial transactional schema migration for `schema_migrations` and `profiles`, a five-second SQLite busy timeout, UUID Profile identifiers, and timestamps. `internal/profile` owns validation and service behavior while `internal/persistence` is its storage adapter; no raw credential, password, token, or private-key field exists in the Profile model.

## Resolution and health

Resolution evaluates normalized repository identity against the precedence `Exact Repository > Remote > Folder > Default`, detects ambiguous/conflicting rules, and returns the winning Profile with an explanation. Remote analysis resolves the current branch's effective fetch/pull and default push targets from Git configuration and evaluates every remote URL. Different remote-to-Profile matches without an Exact Repository rule are a Routing Conflict with `Unknown` health and block materialization. An Exact Repository rule fixes the expected Profile, but Identity Guard still classifies an incompatible effective fetch or push target as `Mismatch` and incompatible secondary remotes as `Warning`. Guard also compares `user.name`, `user.email`, credential mechanism, SSH route if used, signing configuration, config origin, and unexpected local overrides. Health is `Healthy`, `Warning`, `Mismatch`, `Broken`, or `Unknown`; each state must include the observed facts and cause.

`internal/routing` implements the initial pure resolution engine. It normalizes drive/case path forms, supports HTTPS and SSH remote URL matching, returns the winning rule and explanation, and blocks materialization on different Remote-rule Profiles. The Git adapter does not yet supply effective fetch/pull or push targets to this engine.

## Native configuration model

The planned layout is:

```text
User-owned Git config
  + [include] path = C:/Users/User/.gitcompass/gitconfig
  -> GitCompass-owned root fragment
  -> GitCompass-owned profile fragments
  -> effective Git configuration
```

The managed directory is expected to contain `C:\\Users\\User\\.gitcompass\\gitconfig` and `profiles\\personal.gitconfig`, `profiles\\company.gitconfig`, and equivalent future fragments. Its include order is Default, broad Folder, narrow Folder, Remote, then Exact Repository. Folder and normal Exact rules use `gitdir/i` against normalized Windows paths; Remote rules use `hasconfig:remote.*.url` only after a capability check and never include fragments that define remote URLs. A moved Folder rule re-evaluates by location, a Remote rule continues only while its URL matches, and a moved Exact rule becomes stale until explicitly rebound. Materialization blocks rather than silently falling back when a condition cannot be represented. Repository-local writes are minimal, reversible, ownership-tracked exceptions requiring approval of the exact file, key or managed include, current value, proposed value, and reason.

`internal/materialize` currently creates managed profile fragments and a root fragment, then appends one minimal global include. `Manager.Apply` writes all Profile fragments, writes the root, and calls `ensureInclude`; no production service calls it yet. Profile values are interpolated directly, each file uses a fixed `.tmp` replacement path, the global include is detected by substring, and no lock, cleanup, manifest, or removal path exists. The root groups Default, Folder, Remote, and Exact rules but preserves caller order within each kind. Controlled Git 2.53 probes confirmed that current application patterns cannot be emitted unchanged: Folder needs recursive Git syntax, Exact needs the absolute Git directory, and Remote needs explicit Git URL glob families. The selected replacement architecture and reasoning live in D12 (desired-state materialization lifecycle) in `gitcompass.hallucination.md`.

## Selected P4 continuation

The decided materialization flow is `Inspect -> Plan -> Apply -> Verify`. It is not implemented behavior yet.

```text
User global config
  -> exact include: ~/.gitcompass/gitconfig
Stable GitCompass root
  -> include: generations/<desired-hash>/gitconfig
Immutable generation root
  -> immutable Profile fragments
```

The compiler validates one Default, duplicate and conflicting rules, and Profile values before file creation. It sorts broad Folder rules before narrow Folder rules with stable tie-breakers, compiles Exact from `git rev-parse --absolute-git-dir`, and expands each logical Remote rule into explicit HTTPS and SCP-style SSH Git globs. Git executable/version-specific semantic probes gate every condition family. Native `git config` performs serialization and exact fixed-value mutation.

SQLite should add only the metadata required to prove lifecycle ownership: managed installation and active generation state, a hash-bearing artifact manifest, and approval-fingerprinted local exceptions. Apply builds and validates a new immutable generation before switching the stable root. Remove first detaches the exact global include, verifies inactive origins, then deletes only manifest-listed artifacts whose ownership remains valid. A stale or changed artifact is preserved and reported.

Exact rules retain the worktree display path and observed absolute Git directory. Missing, inaccessible, moved, and explicitly rebound states remain distinct; no automatic rebind is permitted. Linked worktrees receive a dedicated semantic acceptance path. A P4 local exception is limited to one approved local `include.path` targeting an owned per-repository fragment, with state re-read immediately before mutation; direct user-owned key repair remains P5 scope.

## Authentication evidence model

The Git adapter must expose the ordered credential-helper chain and origin/scope of the effective credential settings, not collapse repeated helpers into one value. A credential context contains protocol, exact host, optional username, and path only when effective `credential.useHttpPath` is true. Identity Guard combines core inspection with optional capability-probed provider adapters and returns `Configured`, `Bound`, `Operation verified`, or `Unknown` evidence without requesting credential payloads. Automatic checks never invoke `git credential fill` or a helper `get`; only an explicit user action may run a non-mutating network operation, and its success is stored as a timestamped operation observation rather than proof of a human account. D13 (evidence-labeled authentication verification) owns the decision reasoning.

## Repository discovery lifecycle

`internal/repository` will own explicit registration, candidate validation, user-selected roots, scan budgets, cancellation, and inventory state. `internal/persistence` will store roots, exclusions, scan checkpoints, last-complete timestamps, repository canonical/display paths, `missing_since`, and state. A scan has at most two concurrent root workers and a per-root pass budget of 30 seconds or 50,000 visited directories; exhausted work persists a resumable frontier. It detects a `.git` directory or file without entering `.git`, then delegates validation to the Git adapter. Reparse points and offline, network, or removable targets require explicit scope. No v1 watcher, service, OS-recents integration, or IDE-history integration exists. Inventory transitions are Active to Unavailable for inaccessible scope, Active to Missing only after a complete accessible scan, Missing to Archived after 30 days, and any non-deleted state to Active on rediscovery. D14 (opt-in bounded repository discovery) owns the decision reasoning.

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
