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

## Q2 - Multi-remote policy - open

The proposal identifies repositories whose multiple remotes imply different contexts but does not define whether routing should choose one remote, require an explicit exact rule, mark the repository Warning, or block materialization. This decision affects P3, P4, and health classification.

## Q3 - Conditional-include versus local-write selection - open

The materialization design allows conditional includes and repository-local values when required, but no decision defines the exact selection rule. Before P4, specify which repository paths and cases can safely use `includeIf`, how moved repositories are handled, and when local writes are permitted.

## Q4 - Local persistence schema and migration strategy - open

SQLite is a suitable candidate for Profiles, rules, known repositories, discovery roots, metadata, health cache, preferences, and migrations, but it is not yet selected or designed. Choose the schema, locking behavior, backup/recovery expectations, and migration guarantees before P2/P6 persistence implementation.

## Q5 - Discovery and background behavior - open

User-added roots, recently inspected repositories, manual registration, and known workspaces are candidate discovery inputs. Continuous whole-machine scanning is not required. Before P6, define opt-in defaults, performance limits, privacy boundaries, and how moved/deleted repositories age out of the inventory.
