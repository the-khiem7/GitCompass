---
baseline_schema: "2.0"
pack: "gitcompass"
document: "hallucination"
status: "active"
updated: "2026-09-02"
code_ref: "unknown"
---

# GitCompass decisions and open questions

## D1 - Native Git remains authoritative - decided

GitCompass configures and verifies the installed Git rather than intercepting ordinary Git operations or implementing Git itself. This preserves familiar `git commit`, `git pull`, and `git push` behavior when the desktop app is closed and avoids a mandatory command-line workflow. A Git proxy, daemon, or custom Git implementation was rejected because it would make the product a second source of truth and introduce a fragile workflow dependency.

## D2 - Windows-first local desktop product - decided

Windows 10/11, Git for Windows, Git Credential Manager, Windows Credential Manager, and the Wails desktop runtime are the first supported environment. Linux and macOS influence domain boundaries but are not v1 delivery scope. Broad cross-platform delivery now was rejected because it would dilute validation of the Windows-specific paths, credentials, locks, and packaging on which the v1 promise depends.

## D3 - HTTPS first, SSH supported - decided

HTTPS using Git Credential Manager and native credential helpers is the primary authentication strategy; SSH is a first-class alternative with key references or paths only. GitCompass must never own raw tokens, passwords, credential payloads, or private-key contents. Building a vault or password manager was rejected because platform credential stores already own that security boundary.

## D4 - Profiles and deterministic routing - decided

The primary abstraction is a reusable Profile containing identity, auth context, optional SSH, signing, and metadata. A repository resolves through Exact Repository, Remote, Folder, then Default precedence. This prevents a broad folder rule from hiding a repository-specific intent and makes selection explainable. Manual per-repository setup as the core workflow was rejected because it recreates the memory burden the product addresses.

## D5 - Managed fragments, minimal ownership - decided

GitCompass adds only a minimal include to the user's global Git config and owns generated files below `C:\\Users\\User\\.gitcompass\\`, including a root `gitconfig` and profile fragments. User config and repository-local manual overrides remain user-owned until an explicit approved Fix. Rewriting the global file, entire repository configuration, or credential stores was rejected because it makes cleanup destructive and obscures ownership.

## D6 - Explain before repair - decided

Identity Guard must detect, explain, recommend, and require user approval before applying the smallest safe change. Each potentially destructive write must show target file, key, current value, proposed value, and reason. Silent repair was rejected because a local override may represent deliberate user intent and trust depends on transparent changes.

## D7 - No public CLI in current scope - decided

The public product surface is a desktop utility. Business logic still belongs in independently testable Go services behind Wails bindings, so a future CLI can be added without moving domain logic out of the GUI. A public CLI is deferred until real demand justifies its support and compatibility burden.

## D8 - Mandatory implementation skills - decided

Every Go implementation change must use the repository `use-modern-go` skill, including its Modern Go Guidelines CLI before editing the target Go file. Every Wails v3 implementation action must use both repository skills: `wails3` for current v3 guidance and `wails:wails` to detect the version and prevent API mixing. This protects the codebase against outdated Go idioms and the incompatible Wails v2/v3 API, binding, runtime, CLI, and build models. Ignoring it risks compiling against the wrong interface, generating mismatched frontend bindings, or treating unstable v3 behavior as stable. Treating skill use as optional or relying on generic Wails examples was rejected because it cannot consistently enforce current, v3-only implementation practice; Wails v2 is rejected for this project unless a future explicit decision changes the stack.

## Q1 - Authentication-context verification - open

The proposal distinguishes expected HTTPS context from actual credential selection, but it does not establish a safe, provider-neutral way to prove which account Git Credential Manager or another helper will use without exposing credentials or triggering unwanted authentication. Before P5, define what can be reported as configured, inferred, verified, or unknown for each helper.

## D9 - SQLite local persistence and forward-only migrations - decided

This closes Q4. SQLite is the single local persistence store for Profiles, routing rules, known repositories, discovery roots, managed-config metadata, settings, schema state, and migrations; health cache is rebuildable. It stores intent and metadata only, never raw credentials, tokens, passwords, private-key contents, or credential payloads; system Git and platform credential stores remain authoritative. Core entities use a versioned relational schema with stable identifiers, and repository paths retain display and normalized forms for Windows matching. The backend persistence service is the sole writer, uses short transactions and a busy timeout for transient locks, and applies embedded, sequential migrations transactionally. A consistent backup is made before schema upgrade; a failed migration rolls back, and an older application version refuses a newer schema rather than modifying it. Early delivery supports forward migration only, not automatic downgrade. SQLite was selected because it is local-first, single-file, transactional, indexed, and suited to the relational Profile/Rule/Repository model. Ignoring this would leave P2 and P6 without a safe persistence boundary and make schema upgrades or recovery ambiguous. JSON/YAML was rejected because integrity, atomic updates, and schema evolution become fragile; a key-value store was rejected because it adds manual relationship and indexing work; a separate database service was rejected because it exceeds the desktop local-first scope.

## D10 - Operation-aware multi-remote routing - decided

This closes Q2. GitCompass never selects `origin` merely by name. It determines the current branch's effective fetch/pull remote and default push remote from Git's native configuration precedence, while evaluating every configured remote URL against Remote rules. If all matched remotes resolve to one Profile, routing remains deterministic. If remotes resolve to different Profiles and no Exact Repository rule exists, GitCompass reports a Routing Conflict with `Unknown` health and blocks materialization. An Exact Repository rule fixes the expected repository Profile and permits materialization, but it does not declare every remote healthy: an incompatible effective fetch or push target is `Mismatch`, while only incompatible secondary remotes produce `Warning`. Remote URLs, credential bindings, and account selection are never rewritten silently; each change remains an explicit Fix. This policy follows Git's operation-specific remote behavior and preserves the separation between commit identity and authentication context. Ignoring it could materialize the Profile for one remote while ordinary Git uses another. Hard-coding `origin`, choosing the first remote, treating every multi-remote repository as broken, and allowing ambiguous materialization with only a warning were rejected because each either misstates Git behavior or permits unsafe configuration.

## D11 - Conditional-include-first materialization - decided

This closes Q3. GitCompass uses conditional includes as the default and orders its owned fragments from Default to broad Folder, narrow Folder, Remote, then Exact Repository so the established precedence remains effective. Default uses an unconditional managed include; Folder and normal Exact Repository rules use case-insensitive `gitdir/i` conditions on normalized Windows paths; Remote rules use `hasconfig:remote.*.url` only after a capability check, and their included fragments must not define remote URLs. Folder rules re-evaluate after a move, Remote rules continue only while the URL still matches, and a moved Exact Repository rule becomes stale and requires an explicit rebind. GitCompass does not silently fall back to repository-local writes when a conditional rule cannot be represented. A local write is allowed only after the user approves the displayed file, key or managed include, current value, proposed value, and reason; the change must be minimal, ownership-tracked, and reversible. Existing local overrides remain user-owned and are removed or replaced only through that Fix flow. This keeps ordinary Git working while the app is closed without making GitCompass the owner of `.git/config`. Ignoring it would make rule precedence dependent on accidental file order, let local overrides silently defeat intent, or leave moved repositories attached to stale paths. Directly scattering Profile values into local config, automatically rebinding moved repositories, and silently using local writes as a compatibility fallback were rejected because they obscure ownership and make cleanup unsafe.

## Q5 - Discovery and background behavior - open

User-added roots, recently inspected repositories, manual registration, and known workspaces are candidate discovery inputs. Continuous whole-machine scanning is not required. Before P6, define opt-in defaults, performance limits, privacy boundaries, and how moved/deleted repositories age out of the inventory.
