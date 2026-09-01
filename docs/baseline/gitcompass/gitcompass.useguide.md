---
baseline_schema: "2.0"
pack: "gitcompass"
document: "useguide"
status: "active"
updated: "2026-09-02"
code_ref: "unknown"
---

# Planned product and operator contract

## Developer workflow

1. Create a Profile with a commit identity, HTTPS context and credential-helper reference, optional SSH reference, optional signing configuration, and metadata.
2. Add routing rules. Exact repository wins over remote, remote wins over folder, and folder wins over the Default Profile.
3. Apply managed configuration. GitCompass writes only its owned fragments and the minimal include; normal Git commands remain the operational interface.
4. Inspect a repository to see its path, remote, resolved Profile, effective identity, authentication mechanism, signing, health, matched rule, and value origins.
5. When Identity Guard reports a mismatch, review expected versus actual values, the source of the difference, and the recommended minimal change. Approve the Fix explicitly before it is applied.

## Repository health contract

`Healthy` means routing is deterministic, effective identity matches the Profile, the HTTPS/SSH context and remote are coherent, signing matches where configured, and no unexpected override breaks intent. `Warning` means the configuration may work but is ambiguous. `Mismatch` means observed context differs from expected. `Broken` means required configuration is invalid, such as a missing credential helper, missing SSH key, or malformed config. `Unknown` means GitCompass cannot confidently determine the effective context. A color or label alone is insufficient; every state must explain its evidence and cause.

## Safety contract

GitCompass must never silently rewrite user-owned configuration. A proposed write shows the target file, key, current value, new value, and reason, then makes the smallest approved change. A typical repair removes one repository-local `user.email` override so the routed Profile can inherit; it does not rewrite the whole `.git/config`. Uninstall removes the GitCompass include and GitCompass-owned files only.

## Desktop information contract

Home summarizes Profiles, known repositories, healthy count, and attention items. Profiles manage identity, HTTPS, SSH, signing, and metadata. Rules display their precedence and conflict explanation. Repositories expose resolution, effective context, health, and origins. Settings cover Git executable discovery, Git Credential Manager discovery, managed location, discovery roots, diagnostic export, and future SSH settings.

## Failure handling contract

The product must handle Git missing or undiscoverable, unavailable Git Credential Manager, misconfigured helper, invalid config, broken include, duplicate/conflicting rules, moved repository, changed remote, local override masking a Profile, incompatible multiple remotes, missing SSH key, unsupported credential behavior, and locked config files with actionable explanations. It must not represent a configured credential helper as proof that an external authentication attempt succeeded.
