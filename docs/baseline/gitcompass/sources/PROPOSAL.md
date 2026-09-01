# GitCompass — Product & Business Analysis Proposal

> **Document role:** Source of Truth  
> **Checkpoint:** Updated Baseline — 2026-08-26  
> **Product:** GitCompass  
> **Domain:** Multi-account Git identity, authentication context, and repository routing  
> **Primary product surface:** Desktop utility  
> **Primary platform:** Windows  
> **Future platforms:** Linux, macOS  
> **Recommended implementation:** Go + Wails 3 + Vue 3  
> **Authentication strategy:** HTTPS-first, SSH as a first-class supported alternative  
> **Git integration principle:** Native Git remains the source of truth  
> **Public CLI:** Not planned

---

# 1. Executive Summary

GitCompass is a local-first developer utility that ensures every Git repository is used with the correct Git identity and authentication context.

A developer may simultaneously work with:

- a company-hosted GitLab CE instance;
- a personal GitHub account;
- client repositories;
- open-source repositories;
- additional Git providers or organizations;
- different commit identities;
- different HTTPS credential contexts;
- different SSH keys;
- different signing keys.

Git already provides the primitives needed to represent most of these configurations, but the developer is still responsible for remembering which identity belongs to which repository, which account should authenticate against which remote, and how the effective Git configuration is being resolved.

GitCompass turns that manual responsibility into an explicit system:

```text
Git Profiles
     ↓
Routing Rules
     ↓
Repository Context Resolution
     ↓
Native Git Configuration
     ↓
Identity Guard
```

The core product promise is:

> **The right Git identity for the right repository, automatically and visibly.**

GitCompass does not replace Git, does not intercept ordinary Git operations, and does not introduce a new command-line workflow that developers must learn.

Instead, it configures and verifies Git so that normal commands continue to work:

```bash
git commit
git pull
git push
```

even when GitCompass itself is closed.

---

# 2. Background

A professional development machine commonly contains more than one Git context.

Example:

```text
Personal
├── GitHub
├── personal commit identity
├── personal HTTPS credentials
└── optional personal SSH/signing keys

Company
├── self-hosted GitLab CE
├── company commit identity
├── company HTTPS credentials
└── optional company SSH/signing keys
```

The same machine, editor, shell, and Git installation are used for both environments.

Typical friction occurs after:

```text
git clone
git init
moving a repository
changing remotes
working across company and personal folders
```

The developer must repeatedly remember:

```text
Which email should this repo use?
Which account should authenticate?
Which credential is Git actually selecting?
Is a local override masking my intended config?
```

This reduces DX and creates preventable mistakes.

---

# 3. Problem Statement

## 3.1 Wrong commit identity

A company repository may accidentally receive commits authored with a personal identity.

Example:

```text
Repository:
git.company.internal/backend/payment

Expected:
kai@company.com

Actual:
personal@example.com
```

Consequences:

- privacy leakage;
- company identity leakage;
- incorrect attribution;
- incorrect contribution statistics;
- signing inconsistencies;
- cleanup and history-rewrite pressure.

---

## 3.2 Wrong authentication context

Commit identity and authentication are separate concerns.

A repository may use the correct:

```text
user.name
user.email
```

while Git authenticates with the wrong account.

Examples:

```text
Personal GitHub credential
→ company remote
```

or:

```text
Company GitLab credential
→ personal remote
```

This may cause:

- authentication failures;
- wrong-account access;
- confusing credential prompts;
- accidental use of the wrong account;
- repeated troubleshooting.

---

## 3.3 Manual repository setup

After cloning or initializing a repository, developers often run commands such as:

```bash
git config user.name ...
git config user.email ...
```

They may also need to manage:

- credential helpers;
- HTTPS credential scopes;
- SSH identities;
- signing keys;
- repository-specific overrides.

This is exactly the repeated CLI friction GitCompass is intended to remove.

---

## 3.4 Fragmented configuration

The effective Git context may come from multiple sources:

```text
system config
global config
conditional include
repository-local config
worktree config
environment
credential helper
SSH config
```

When something is wrong, Git does not provide a cohesive product UI for understanding:

```text
What profile should apply?
What actually applies?
Why?
Where did the effective value come from?
```

GitCompass exists to make this visible and explainable.

---

# 4. Product Thesis

GitCompass should model the problem around a first-class abstraction:

> **Git Profile**

A Git Profile represents the intended Git context for a category of repositories.

```text
Git Profile
├── Identity
├── Authentication
├── Signing
└── Metadata
```

Routing is modeled separately:

```text
Routing Rule
├── matcher
├── target profile
└── priority
```

Together:

```text
Repository
     ↓
Routing Engine
     ↓
Resolved Profile
     ↓
Configuration Materialization
     ↓
Native Git
     ↓
Identity Guard
```

The separation is intentional:

- **Profile** answers: “Who am I in this Git context?”
- **Routing Rule** answers: “Where does this Profile apply?”
- **Materialization** answers: “How is that intent expressed through native Git configuration?”
- **Identity Guard** answers: “Is the repository actually using the intended context?”

---

# 5. Core Product Principles

## 5.1 Configure Git, do not replace Git

GitCompass should work with the user's installed Git.

Normal Git operations remain normal Git operations:

```bash
git commit
git pull
git push
```

GitCompass must not require:

```bash
gitcompass commit
gitcompass push
```

---

## 5.2 Setup once, use automatically

After a Profile and Routing Rule are configured, the developer should not repeatedly set identity by hand.

---

## 5.3 Active identity must be visible

The user should be able to inspect a repository and immediately know:

```text
Expected Profile
Effective identity
Authentication method
Remote context
Signing state
Health status
```

---

## 5.4 Routing must be deterministic

A repository should never mysteriously resolve to a Profile.

The UI must explain:

```text
Matched rule
Rule priority
Competing matches if relevant
```

---

## 5.5 GitCompass should not own secrets by default

Prefer:

```text
credential helper
OS-backed secure storage
SSH key reference
```

over:

```text
raw PAT stored in app database
raw password stored in app database
private key copied into app storage
```

---

## 5.6 GitCompass should own only its own configuration

The application must not rewrite the user's entire Git configuration.

It should manage isolated config fragments and only modify user-owned config when explicitly authorized.

---

# 6. Target User

The primary user is a developer who uses one Windows machine for multiple Git contexts.

Typical setup:

```text
Windows
Git for Windows
GitHub personal account
Company GitLab CE
HTTPS authentication
Git Credential Manager
VS Code or another IDE
```

Secondary/future setups may use:

```text
SSH
Linux
macOS
other Git providers
```

---

# 7. Authentication Strategy

## 7.1 HTTPS First

HTTPS is the primary authentication path for GitCompass.

The preferred architecture is:

```text
GitCompass Profile
        ↓
Native Git configuration
        ↓
Git Credential Manager / Git credential helper
        ↓
Windows Credential Manager
        ↓
GitHub / GitLab / Git provider
```

GitCompass should orchestrate and explain this relationship rather than becoming a credential vault.

A conceptual Profile may contain:

```text
Company

Identity
├── Name: Kai Nguyen
└── Email: kai@company.com

Authentication
├── Protocol: HTTPS
├── Remote host: git.company.internal
└── Credential mechanism: Git Credential Manager
```

Another Profile:

```text
Personal

Identity
├── Name: Personal Developer
└── Email: personal@example.com

Authentication
├── Protocol: HTTPS
├── Remote host: github.com
└── Credential mechanism: Git Credential Manager
```

---

## 7.2 HTTPS Credential Context

GitCompass must distinguish between:

```text
commit identity
```

and:

```text
credential selection
```

For HTTPS, the product should inspect and surface:

```text
remote URL
credential helper
credential scope/context
provider host
effective authentication mechanism
```

Where possible, GitCompass should rely on Git and Git Credential Manager behavior rather than duplicating authentication logic.

---

## 7.3 SSH Support

SSH remains an important first-class alternative.

GitCompass should support:

```text
SSH key path
SSH host alias
core.sshCommand
IdentitiesOnly
signing key references
```

Example:

```ini
[core]
    sshCommand = ssh -i ~/.ssh/id_ed25519_company -o IdentitiesOnly=yes
```

The proposal does not position GitCompass as SSH-first.

SSH is supported as a major authentication method alongside the HTTPS-first happy path.

---

## 7.4 Secret Ownership

Preferred:

```text
GitCompass
→ reference authentication mechanism
```

Avoid:

```text
GitCompass DB
→ raw token / password
```

GitCompass should not read or display private key contents.

---

# 8. Git Profile Model

Recommended conceptual shape:

```text
Profile
├── id
├── name
├── Identity
│   ├── userName
│   └── userEmail
├── Authentication
│   ├── protocol
│   ├── host
│   ├── credential mechanism
│   └── optional SSH settings
├── Signing
│   ├── enabled
│   ├── format
│   └── key reference
└── metadata
```

Example:

```yaml
id: company
name: Company

identity:
  userName: Kai Nguyen
  userEmail: kai@company.com

authentication:
  protocol: https
  host: git.company.internal
  credentialProvider: git-credential-manager

signing:
  enabled: false
```

This is a business/domain model, not a committed on-disk schema.

---

# 9. Identity vs Authentication

GitCompass must keep these separate in both UI and implementation.

## Identity

Controls commit authorship:

```ini
[user]
    name = Kai Nguyen
    email = kai@company.com
```

## Authentication

Controls access to the remote.

Typical GitCompass priority:

```text
HTTPS
→ Git Credential Manager / credential helper

SSH
→ key / host / sshCommand
```

A repository is not healthy simply because `user.email` is correct.

Both identity and authentication context must be coherent.

---

# 10. Routing Engine

Recommended precedence:

```text
1. Exact Repository Rule
2. Remote Rule
3. Folder Rule
4. Default Profile
```

---

## 10.1 Exact Repository Rule

Most specific.

Example:

```text
C:\Projects\SpecialClient
→ Client A
```

---

## 10.2 Remote Rule

Useful when repository ownership is primarily associated with a provider/server.

Example:

```text
git.company.internal/**
→ Company
```

or:

```text
github.com/**
→ Personal
```

Remote rules are particularly important for HTTPS-first authentication.

---

## 10.3 Folder Rule

Useful for workspace conventions.

Example:

```text
C:\Projects\Company\**
→ Company
```

---

## 10.4 Default Profile

Fallback when no specific rule matches.

---

# 11. Routing Conflict Example

Rules:

```text
Folder:
C:\Projects\Company\**
→ Company

Remote:
github.com/**
→ Personal
```

Repository:

```text
C:\Projects\Company\open-source-project

Remote:
https://github.com/user/project.git
```

Resolution:

```text
Remote rule
→ Personal
```

The UI should explain:

```text
Resolved Profile
Personal

Matched rule
Remote: github.com/**

Also matched
Folder: C:\Projects\Company\**
```

This makes routing behavior transparent.

---

# 12. Public CLI Decision

## Public CLI is not planned

GitCompass is specifically intended to reduce repeated Git-related CLI typing.

A product interface such as:

```bash
gitcompass profile
gitcompass resolve
gitcompass doctor
gitcompass apply
```

would add another command vocabulary for users to learn while providing limited product ROI.

Therefore:

```text
Public Product Interface
└── Desktop GUI
```

The initial product should not include a public CLI.

---

## 12.1 No CLI does not mean business logic belongs in the GUI

The architecture must still remain layered.

```text
Vue UI
  ↓
Go Application Services
  ↓
Git / Routing / Guard Engines
  ↓
System Git
```

Example service layer:

```text
ProfileService
RoutingService
RepositoryService
MaterializationService
IdentityGuardService
GitService
```

This allows:

- clean unit testing;
- separation of concerns;
- future integrations if justified;
- UI replacement without rewriting business logic.

---

## 12.2 Future CLI reconsideration

A CLI may only be reconsidered if concrete demand appears for:

```text
enterprise automation
CI integration
agent integration
scripting
headless operations
```

It should not be built speculatively.

---

# 13. Native Git Configuration Strategy

The core rule is:

> **GitCompass expresses intent using configuration normal Git understands.**

Folder-based routing may use native conditional includes.

Example:

```ini
[includeIf "gitdir:C:/Projects/Company/"]
    path = C:/Users/User/.gitcompass/profiles/company.gitconfig
```

Profile fragment:

```ini
[user]
    name = Kai Nguyen
    email = kai@company.com
```

For repository-specific cases, GitCompass may use repository-local configuration.

The business model must stay independent from a particular materialization technique.

---

# 14. Intent vs Materialization

GitCompass should explicitly separate:

```text
Intent
Repository X should use Company Profile
```

from:

```text
Materialization
How that intent is represented in Git configuration
```

This matters because different rules may require different implementations.

Example:

```text
Folder Rule
→ includeIf

Exact Repo Rule
→ repository-local config

Remote Rule
→ resolved intent + appropriate managed config/local materialization
```

Routing stays a domain concept.

Git config remains an implementation detail.

---

# 15. Managed Configuration Ownership

This is a critical safety design.

GitCompass must not assume ownership of the user's entire `.gitconfig`.

---

## 15.1 Bad model

Do not do:

```text
GitCompass
→ parse whole ~/.gitconfig
→ rewrite whole file
→ become owner of everything
```

This risks:

- destroying unrelated aliases;
- overwriting user edits;
- conflicting with other tools;
- making uninstall/recovery difficult.

---

## 15.2 Preferred model

Keep user configuration user-owned.

Example:

```text
~/.gitconfig
```

GitCompass adds only a minimal include:

```ini
[include]
    path = C:/Users/User/.gitcompass/gitconfig
```

Then GitCompass owns:

```text
C:\Users\User\.gitcompass\
├── gitconfig
└── profiles\
    ├── personal.gitconfig
    └── company.gitconfig
```

The effective model becomes:

```text
User-owned Git config
        +
GitCompass-owned config fragments
        ↓
Effective Git configuration
```

---

## 15.3 Ownership rule

```text
User config
→ user owns

GitCompass generated files
→ GitCompass owns

Repository-local manual overrides
→ user owns until explicit Fix is approved
```

---

## 15.4 Uninstall behavior

A clean uninstall should ideally require only:

```text
remove GitCompass include
remove GitCompass-owned files
```

The user's unrelated Git configuration remains intact.

---

# 16. Identity Guard

Identity Guard is the primary safety feature.

Its job:

```text
Expected Profile
      VS
Effective Git Context
```

Checks may include:

```text
user.name
user.email
HTTPS credential mechanism
credential helper
remote host
SSH route if applicable
signing configuration
config origin
unexpected local overrides
```

---

# 17. Identity Guard Example

Repository:

```text
C:\Projects\Company\payment
```

Expected:

```text
Profile
Company

Email
kai@company.com

Remote
https://git.company.internal/backend/payment.git
```

Actual:

```text
user.email
personal@example.com
```

Cause:

```text
.git/config
contains repository-local override
```

GitCompass shows:

```text
⚠ Identity mismatch

Expected
Company
kai@company.com

Actual
personal@example.com

Cause
Repository-local user.email overrides routed profile.

Recommended Fix
Remove repository-local override and inherit Company Profile.
```

---

# 18. Safe Fix Philosophy

GitCompass should not silently rewrite user-owned config.

The flow:

```text
Detect
   ↓
Explain
   ↓
Recommend
   ↓
User approves Fix
   ↓
Apply minimal change
```

Example:

```text
Remove one local user.email override
```

not:

```text
Rewrite entire repository config
```

---

# 19. Repository Health Model

Suggested states:

```text
Healthy
Warning
Mismatch
Broken
Unknown
```

## Healthy

Everything matches expected Profile.

## Warning

Configuration may work but is ambiguous.

## Mismatch

Effective context is known to differ from expected Profile.

## Broken

Required configuration is invalid.

Examples:

```text
credential helper missing
SSH key path missing
Git config malformed
```

## Unknown

GitCompass cannot confidently determine the effective context.

Health must always be explainable.

---

# 20. Repository Inventory

The desktop app should provide a known-repository inventory.

Example:

```text
Repositories

✓ payment-service       Company
✓ infrastructure        Company
⚠ portfolio             Personal
✗ internal-tool         Company
```

Key columns/fields:

```text
Repository
Path
Remote
Resolved Profile
Effective Identity
Authentication Method
Health
```

---

# 21. Repository Discovery

GitCompass does not need to continuously scan the entire machine.

Possible discovery methods:

```text
user-added root folders
recently inspected repositories
manual repository registration
known workspace folders
```

Future optional background discovery can be considered only if it improves DX without unnecessary system cost.

---

# 22. Desktop UX

Primary navigation:

```text
Home
Profiles
Rules
Repositories
Settings
```

---

## 22.1 Home

Example:

```text
4 Profiles
26 Known Repositories
24 Healthy
2 Need Attention
```

Recent issues:

```text
portfolio
Wrong email

backend-api
Credential helper issue
```

---

## 22.2 Profiles

Manage:

```text
Identity
HTTPS authentication context
SSH fallback/support
Signing
Metadata
```

---

## 22.3 Rules

Manage:

```text
Exact Repo
Remote
Folder
Default
```

Show precedence and conflict explanations.

---

## 22.4 Repositories

Inspect:

```text
Resolved Profile
Effective Identity
Remote
Credential Mechanism
Signing
Health
Config Origin
```

---

## 22.5 Settings

Potential settings:

```text
Git executable
Git Credential Manager discovery
managed config location
repository discovery roots
diagnostic export
future SSH settings
```

---

# 23. Repository Detail UI Concept

```text
┌────────────────────────────────────────────────────┐
│ inventory-service                                  │
├────────────────────────────────────────────────────┤
│ Path                                               │
│ C:\Projects\Company\inventory-service               │
│                                                    │
│ Remote                                             │
│ https://git.company.internal/backend/inventory.git │
│                                                    │
│ Expected Profile                                   │
│ Company                                            │
│                                                    │
│ Effective Identity                                 │
│ Kai Nguyen <kai@company.com>                       │
│                                                    │
│ Authentication                                     │
│ HTTPS                                              │
│ Git Credential Manager                             │
│                                                    │
│ ✓ Configuration healthy                           │
│                                                    │
│ Matched Rule                                       │
│ Remote: git.company.internal/**                    │
└────────────────────────────────────────────────────┘
```

---

# 24. Profile UI Concept

```text
┌────────────────────────────────────────────────────┐
│ Company                                            │
├────────────────────────────────────────────────────┤
│ Identity                                           │
│ Name      Kai Nguyen                               │
│ Email     kai@company.com                          │
│                                                    │
│ Authentication                                     │
│ Protocol  HTTPS                                    │
│ Host      git.company.internal                     │
│ Provider  Git Credential Manager                   │
│                                                    │
│ SSH Support                                        │
│ Optional / not configured                          │
│                                                    │
│ Signing                                            │
│ Disabled                                           │
│                                                    │
│ Used by                                            │
│ 8 repositories · 2 rules                          │
└────────────────────────────────────────────────────┘
```

---

# 25. Architecture

Recommended:

```text
GitCompass
│
├── Git Engine
├── Profile Engine
├── Routing Engine
├── Materialization Engine
├── Identity Guard
├── Repository Inventory
├── Persistence
└── Desktop App
```

No public CLI layer is required.

---

# 26. Git Engine

Responsibilities:

```text
discover Git for Windows
execute Git commands safely
read Git version
read config
read config origin/scope
read remotes
inspect repository
inspect effective identity
inspect credential-helper configuration
validate Git availability
```

System Git remains the source of truth.

---

# 27. Profile Engine

Responsibilities:

```text
Profile CRUD
Identity model
Authentication model
HTTPS context model
SSH support model
Signing model
Profile validation
```

---

# 28. Routing Engine

Responsibilities:

```text
evaluate matching rules
apply precedence
detect conflicts
explain resolution
return expected Profile
```

---

# 29. Materialization Engine

Responsibilities:

```text
translate Profile + Routing intent
→ Git-native configuration
```

It may manage:

```text
GitCompass-owned include fragments
conditional includes
repository-local values when required
```

Requirements:

```text
idempotent
minimal
reversible
ownership-aware
```

---

# 30. Identity Guard Engine

Responsibilities:

```text
compare expected vs effective config
classify repository health
identify config source
detect manual overrides
explain mismatch
apply approved minimal repair
```

---

# 31. Persistence

GitCompass needs local app state.

Suggested stored data:

```text
Profiles
Routing Rules
Known Repositories
Discovery Roots
Managed Config Metadata
Health Cache
UI Preferences
Migration Metadata
```

SQLite is a suitable candidate.

Reasons:

```text
local-first
transactional
easy migrations
suitable for desktop state
inspectable
```

Git-native configuration itself remains outside the database.

---

# 32. Security Principles

## 32.1 HTTPS credential safety

GitCompass should prefer:

```text
Git Credential Manager
Windows Credential Manager
Git credential helpers
```

and avoid owning raw credentials.

---

## 32.2 SSH safety

Store only key references/paths.

Never display private key contents.

---

## 32.3 Config write transparency

For potentially destructive writes, show:

```text
Target file
Key
Current value
New value
Reason
```

---

## 32.4 Diagnostics redaction

Diagnostic exports may include:

```text
Git version
Git path
Profile names
Remote hosts
Routing decisions
Config origin
Health results
```

but must redact:

```text
tokens
passwords
private key contents
credential payloads
```

---

# 33. Platform Strategy

## Tier 1 — Windows

GitCompass v1 is Windows-first.

Primary environment:

```text
Windows 10/11
Git for Windows
Git Credential Manager
Windows Credential Manager
Wails desktop runtime
```

All architectural decisions should be validated first against Windows behavior.

---

## 33.1 Windows-specific concerns

First-class support must account for:

```text
drive letters
backslashes
case-insensitive paths
Git for Windows locations
Windows OpenSSH
Git Credential Manager
Windows Credential Manager
junctions
long paths
Unicode paths
file locking
```

---

## 33.2 Future — Linux

Linux support is future work.

Likely concerns:

```text
credential helpers
desktop keyrings
filesystem permissions
case-sensitive paths
SSH defaults
desktop packaging
```

---

## 33.3 Future — macOS

macOS support is future work.

Likely concerns:

```text
Keychain
filesystem behavior
SSH
app packaging/signing
native webview differences
```

Linux and macOS should influence abstractions enough to avoid hard-coded Windows-only domain logic, but they should not dilute Windows-first implementation focus.

---

# 34. Technology Stack

## Core

```text
Go
```

Why:

```text
filesystem
process execution
Git CLI orchestration
config handling
Windows-native binaries
straightforward testing
direct Wails integration
```

---

## Desktop

```text
Wails 3
```

Wails is a strong fit because:

```text
Go application services
        ↓
Wails bindings
        ↓
Vue UI
```

No additional backend language is required.

---

## Frontend

```text
Vue 3
TypeScript
Vite
```

Vue is appropriate for:

```text
settings
forms
repository tables
health dashboards
profile management
rule management
```

without excessive frontend complexity.

---

# 35. Suggested Repository Structure

```text
gitcompass/
├── internal/
│   ├── git/
│   ├── profile/
│   ├── routing/
│   ├── materialize/
│   ├── guard/
│   ├── repository/
│   ├── persistence/
│   └── config/
│
├── apps/
│   └── desktop/
│       ├── backend/
│       └── frontend/
│
└── docs/
```

Business logic must not live inside Vue components.

---

# 36. Service Layer

Recommended application services:

```text
ProfileService
RoutingService
RepositoryService
MaterializationService
IdentityGuardService
SettingsService
DiagnosticsService
```

Example flow:

```text
Vue
 ↓
IdentityGuardService.CheckRepository()
 ↓
Routing Engine
 ↓
Git Engine
 ↓
Health Result
```

This preserves testability without requiring a public CLI.

---

# 37. Testing Strategy

## Unit Tests

Focus:

```text
routing precedence
rule conflict resolution
Profile validation
materialization planning
health classification
path normalization
ownership rules
```

---

## Integration Tests

Use temporary Git repositories to test:

```text
global config
GitCompass managed include
conditional includes
repository-local overrides
multiple remotes
HTTPS credential helper configuration
SSH optional configuration
doctor/guard behavior
repair behavior
```

---

## Windows Tests

Priority scenarios:

```text
Git for Windows discovery
Windows path handling
Git Credential Manager
Windows Credential Manager integration assumptions
case-insensitive paths
drive-letter normalization
junction behavior
locked config files
```

---

## Safety Tests

Verify that GitCompass:

```text
does not delete unrelated Git config
does not expose secrets
does not rewrite user-owned config silently
applies managed config idempotently
can remove its managed config cleanly
```

---

# 38. Failure Modes

GitCompass should explicitly handle:

```text
Git not installed
Git executable not found
Git Credential Manager unavailable
credential helper misconfigured
invalid config file
broken include path
duplicate/conflicting routing rules
repository moved
remote changed
local override masks Profile
multiple remotes imply different contexts
missing SSH key
unsupported credential behavior
```

Error messages should be actionable.

---

# 39. Product Roadmap

## Stage 1 — Windows Git Foundation

Deliver:

```text
Git for Windows discovery
Git command execution
repository detection
config origin/scope inspection
remote inspection
effective identity inspection
credential helper inspection
```

---

## Stage 2 — Profile Engine

Deliver:

```text
Profile CRUD
Identity
HTTPS authentication context
credential helper references
optional SSH configuration
signing configuration
Profile validation
```

---

## Stage 3 — Routing Engine

Deliver:

```text
Exact Repository Rule
Remote Rule
Folder Rule
Default Profile
deterministic precedence
resolution explanation
conflict detection
```

---

## Stage 4 — Managed Configuration Materialization

Deliver:

```text
GitCompass-owned config root
managed profile fragments
conditional includes
repository-local materialization when required
idempotent apply
ownership metadata
clean removal
```

---

## Stage 5 — Identity Guard

Deliver:

```text
expected vs effective comparison
health states
config origin explanation
wrong identity detection
HTTPS auth-context diagnostics
SSH diagnostics
safe Fix workflow
```

---

## Stage 6 — Desktop Product

Deliver:

```text
Home
Profiles
Rules
Repositories
Settings
Health views
Fix flows
```

---

## Stage 7 — Windows Hardening

Deliver:

```text
installer
config migrations
diagnostics export
Git Credential Manager edge cases
Windows path edge cases
security review
recovery from malformed config
```

---

## Stage 8 — Future Platform Preparation

Only after Windows product quality is satisfactory:

```text
Linux research
macOS research
platform abstraction validation
packaging strategy
credential-store mapping
```

---

# 40. Product Success Criteria

GitCompass is successful when a developer can:

1. create at least two distinct Git Profiles;
2. define commit identity for each Profile;
3. use HTTPS credentials through native Git/Git Credential Manager;
4. optionally configure SSH where needed;
5. assign Profiles through deterministic routing rules;
6. clone or initialize repositories without repeatedly setting identity by hand;
7. use normal Git commands without GitCompass intercepting each operation;
8. immediately see which Profile a repository resolves to;
9. understand why that Profile resolved;
10. detect wrong identity or authentication context;
11. inspect where the effective value came from;
12. safely fix a mismatch;
13. close GitCompass and still have normal Git use the intended native configuration.

---

# 41. Definition of a Healthy Repository

A repository is healthy when:

```text
Routing resolves deterministically
        +
Effective commit identity matches Profile
        +
HTTPS/SSH authentication context is coherent
        +
Remote context matches expectation
        +
Signing configuration matches Profile
        +
No unexpected override breaks intent
```

Health must be explainable, not merely represented by a green icon.

---

# 42. Non-Goals

GitCompass should not become:

```text
a Git hosting platform
a GitHub Desktop replacement
a full Git GUI
a password manager
a secret vault
a custom Git implementation
a new mandatory Git CLI
a Git proxy daemon
an authentication server
```

The product remains focused on:

```text
Profile
Routing
Visibility
Guard
Repair
```

---

# 43. Future Possibilities

Possible future directions, not part of the current baseline:

```text
Linux support
macOS support
system tray health indicator
repository auto-discovery
enterprise-managed profile templates
team policy import
IDE integrations
agent integrations
headless automation
public CLI if real demand emerges
```

These must be justified by actual usage rather than pre-built speculatively.

---

# 44. Final Product Definition

> **GitCompass is a Windows-first, local-first Git identity and context manager that maps reusable developer Profiles to repositories through deterministic routing rules, expresses those rules through native Git configuration, uses HTTPS credential workflows as the primary authentication path, supports SSH where needed, and verifies the effective repository context through Identity Guard.**

The essential product loop is:

```text
Create Profile
      ↓
Define Routing Rule
      ↓
GitCompass Materializes Native Config
      ↓
Repository Resolves Automatically
      ↓
Normal Git Uses Correct Context
      ↓
Identity Guard Verifies
      ↓
Explain / Fix if needed
```

GitCompass exists to turn multi-account Git from a manual, memory-based convention into a visible, deterministic, and safe developer workflow.

---

# 45. Current Baseline Decisions

```text
Brand
GitCompass

Product category
Local-first Git identity/context manager

Primary platform
Windows

Future platforms
Linux
macOS

Primary UX
Desktop utility

Public CLI
Not planned

Core language
Go

Desktop framework
Wails 3

Frontend
Vue 3 + TypeScript

Git semantics
System Git CLI

Authentication strategy
HTTPS First

Primary HTTPS mechanism
Git Credential Manager / native Git credential helpers

SSH
First-class supported alternative

Primary abstraction
Profile

Routing precedence
Exact Repo → Remote → Folder → Default

Core safety feature
Identity Guard

Configuration philosophy
Native Git

Configuration ownership
GitCompass-owned isolated config fragments

Secrets philosophy
Reference / native secure storage, do not own raw credentials

Interaction philosophy
Configure Git, do not intercept Git
```

These decisions form the current Source of Truth baseline for all subsequent GitCompass product and implementation work.
