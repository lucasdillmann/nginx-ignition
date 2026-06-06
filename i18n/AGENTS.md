# i18n Properties File Guide for AI Agents

This document provides comprehensive instructions for AI models working with the i18n (internationalization) properties files in this project.

## Overview

The `messages_en.properties` file contains all user-facing text strings for the nginx-ignition application. These keys are used in both the Go backend and the TypeScript frontend.

## File Format

### Basic Structure

```properties
key/path/suffix=Value text here
```

- **Key**: Left side of `=`, uses `/` as separator (path-like format)
- **Value**: Right side of `=`, the actual translated text
- **No quotes** around values
- **One key per line** - this is critical and must be preserved
- **No trailing whitespace**

### Key Format

Keys follow a hierarchical path structure that mirrors where they are used in the codebase:

```
{module}/{submodule}/.../{descriptor}
```

#### Examples:
```properties
core/accesslist/in-use=Access list is in use by one or more hosts
frontend/authentication/login-button=Log in
certificate/letsencrypt/dns/azure/client-id=Azure client ID
vpn/tailscale/auth-key=Tailscale auth key
```

## Key Naming Conventions

### 1. Path Prefix = Folder Location

The key prefix MUST match the folder path where the key is used in the codebase:

| Code Location | Key Prefix |
|---------------|------------|
| `core/accesslist/*.go` | `core/accesslist/` |
| `certificate/letsencrypt/dns/azure/*.go` | `certificate/letsencrypt/dns/azure/` |
| `frontend/src/domain/authentication/*.tsx` | `frontend/authentication/` |
| `frontend/src/core/components/shell/*.tsx` | `frontend/components/shell/` |
| `api/common/authorization/*.go` | `api/common/authorization/` |
| `integration/docker/*.go` | `integration/docker/` |
| `vpn/tailscale/*.go` | `vpn/tailscale/` |
| `frontend/src/domain/trafficstats/*.tsx` | `frontend/traffic-stats/` |

**Frontend Exception**: For frontend paths, omit `src/domain/` or `src/core/`:
- `frontend/src/domain/accesslist/` → `frontend/accesslist/`
- `frontend/src/core/components/shell/` → `frontend/components/shell/`

### 2. Suffix = Descriptive Name

The suffix (after the last `/`) should be:
- **Concise**: Use the minimum words needed to describe the purpose
- **Descriptive**: Clearly indicate what the text is for
- **Non-redundant**: Do NOT repeat information already in the path

#### ✅ GOOD Examples:
```properties
certificate/letsencrypt/dns/azure/client-id=Azure client ID
core/accesslist/in-use=Access list is in use
integration/docker/name=Docker
vpn/tailscale/auth-key=Tailscale auth key
frontend/authentication/login-button=Log in
```

#### ❌ BAD Examples (redundant):
```properties
# BAD: "azure" already in path, don't repeat it
certificate/letsencrypt/dns/azure/azure-client-id=Azure client ID

# BAD: "error" is redundant when context implies it
core/accesslist/error-in-use=Access list is in use

# BAD: "validation" is redundant
core/binding/validation-invalid-ip=Value is not a valid IP address

# BAD: "lets-encrypt-dns-azure" all redundant with path
certificate/letsencrypt/dns/azure/lets-encrypt-dns-azure-client-id=Azure client ID
```

### 3. Keys Used in Multiple Places

If a key is used in multiple folders, use the `common/` prefix:

```properties
common/cannot-be-empty=Value cannot be empty
common/invalid-url=Value is not a valid URL
common/value-missing=A value is required
```

Important: If you create a new key and code generation fails with the message that such value already exists, DO NOT 
USE THE KEY THAT ALREADY EXISTS. Move the key to the `common/` group following the examples above.

### 4. Common Suffix Patterns

| Pattern | Usage | Example |
|---------|-------|---------|
| `name` | Display name of a feature/provider | `integration/docker/name=Docker` |
| `description` | Longer description text | `integration/docker/description=Enables...` |
| `{field}` | Form field label | `vpn/tailscale/auth-key=Tailscale auth key` |
| `{field}-help` | Help text for a field | `certificate/letsencrypt/dns/acmedns/allow-list-help=Comma-separated...` |
| `in-use` | Resource is being used | `core/cache/in-use=Cache is in use...` |
| `not-found` | Resource was not found | `core/user/not-found=User not found` |
| `invalid-{thing}` | Validation: invalid value | `core/binding/invalid-ip=Value is not a valid IP` |
| `{thing}-required` | Validation: required field | `vpn/tailscale/auth-key-required=Auth key is required` |

## Value Guidelines

### 1. Placeholder Variables

Use `${variable}` syntax for dynamic values:

```properties
core/cache/invalid-status-code=Invalid status code ${value}: must be between ${min} and ${max}
core/host/duplicated-route-priority=Priority ${priority} is duplicated
```

### 2. Text Style

- Use sentence case (capitalize first word only, except proper nouns)
- Use consistent terminology throughout
- Keep messages user-friendly and actionable
- For error messages, explain what went wrong and ideally how to fix it

### 3. Capitalization

- **Do not use uppercase letters at will**: Follow standard capitalization rules for the language.
- **Sentence case**: Always use sentence case for UI labels and messages. Use uppercase only at the beginning of the sentence or for proper nouns (e.g. brand names, technical terms like 'Nginx').
- **Avoid Title Case**: Do not capitalize every word. E.g., use "Upstream server" instead of "Upstream Server".
- **Context matters**: Ensure capitalization fits the grammatical context.

### 4. Length Considerations

- UI labels: Keep short (1-3 words)
- Error messages: Be descriptive but concise
- Help text: Be descriptive but concise

## Adding New Keys

### Step 1: Determine the Folder Location

Find where in the codebase the key will be used. The folder path becomes the key prefix.

### Step 2: Create a Descriptive Suffix

Choose a suffix that:
- Describes the purpose
- Does NOT include words already in the path
- Follows existing patterns in the same folder

### Step 3: Add to the File

Add the new key in the file. You don't need to worry where to place it, just add it at the bottom of the file, at the 
top or anything else (whichever is the fastest one) and the `make generate-i18n` command will sort it automatically.

### Step 4: Generate Code

After adding keys, run the code generation:
```bash
make .generate-i18n-files
```

This generates:
- `i18n/keys.generated.go` - Go constants
- `i18n/en.generated.go` - Go dictionary
- `frontend/src/core/i18n/model/MessageKey.generated.ts` - TypeScript enum

## Critical Rules

### 🚨 NEVER Use Raw Key Strings in Code

All message keys referenced in Go or frontend TypeScript must use generated constants — never raw path strings.

**Go** — use `i18n.K.*` from `i18n/keys.generated.go`:

```go
// ✅ CORRECT
i18n.M(ctx, i18n.K.CoreUserNotFound)
i18n.DetachedMessage{Key: i18n.K.CoreNotificationCategoryCertificateRenewed}

// ❌ WRONG
i18n.M(ctx, "core/user/not-found")
i18n.DetachedMessage{Key: "core/notification/category/certificate-renewed"}
```

**Frontend** — use `MessageKey.*` from `frontend/src/core/i18n/model/MessageKey.generated.ts`:

```tsx
// ✅ CORRECT
<I18n id={MessageKey.FrontendUserListSubtitle} />
const text = i18n(MessageKey.FrontendUserNewButton)

// ❌ WRONG
<I18n id="frontend/user/list-subtitle" />
<I18n id={'frontend/user/list-subtitle'} />
const text = i18n("frontend/user/new-button")
```

Use `i18n.Static("...")` only for non-localized text (test fixtures, dynamic values) — not for keys in
`.properties` files.

### 🚨 NEVER Use Dots in Keys

Keys use `/` as separator, NOT `.`:
```properties
# ✅ CORRECT
core/user/not-found=User not found

# ❌ WRONG
core.user.not-found=User not found
```

### 🚨 NEVER Duplicate Keys

Each key must be unique. Before adding a new key, search the file to ensure it doesn't already exist.

### 🚨 NEVER duplicate values

Each value must be unique. If the i18n generator fails with a message that a value already exists,
then move the message to the `common/` key group/prefix. Do not reuse the key that already exists, if a message is used
twice in different folders/places, moving it to the `common/` must be done.

Important: ONLY USE common/ IF DUPLICATED. Otherwise, create new keys in the appropriate folder.

### 🚨 Keep Values on Same Line

Do not break values across multiple lines. Each key=value pair must be on a single line.

## File Organization

The file is organized by module/feature area:

1. **Core modules** (`core/...`) - Backend business logic
2. **Frontend** (`frontend/...`) - UI-specific text
3. **API** (`api/...`) - API layer messages
4. **Certificate** (`certificate/...`) - Certificate management
5. **Integration** (`integration/...`) - Third-party integrations
6. **VPN** (`vpn/...`) - VPN features
7. **Database** (`database/...`) - Database layer
8. **Common** (`common/...`) - Shared across modules

## Localization

This file (`messages_en.properties`) is the source of truth. Other locale files should:
- Have the SAME keys in the SAME order
- Only differ in the values (translated text)

### Available Languages

| File | Language | Dialect/Variant | Script | Region/Notes |
|------|----------|-----------------|--------|--------------|
| `messages_en.properties` | English | Standard | Latin | Source of truth |
| `messages_bn.properties` | Bengali | Standard | Bengali (বাংলা) | Bangladesh/India |
| `messages_de.properties` | German | Standard | Latin | Germany/Austria/Switzerland |
| `messages_es.properties` | Spanish | Standard | Latin | General Spanish |
| `messages_fr.properties` | French | Standard | Latin | France/General French |
| `messages_hi.properties` | Hindi | Standard | Devanagari (हिंदी) | India |
| `messages_ja.properties` | Japanese | Standard | Kanji/Hiragana/Katakana | Japan |
| `messages_pt.properties` | Portuguese | Brazilian (pt-BR) | Latin | Brazil |
| `messages_ru.properties` | Russian | Standard | Cyrillic (русский) | Russia |
| `messages_vi.properties` | Vietnamese | Standard | Latin with diacritics | Vietnam |
| `messages_zh.properties` | Mandarin Chinese | Simplified (zh-CN) | Simplified Hanzi (简体中文) | Mainland China (PRC) |

### Script Notes

- **CJK Languages**: Japanese (`ja`) and Chinese (`zh`) use complex character sets
- **Indic Languages**: Bengali (`bn`) uses the Bengali script derived from Brahmi; Hindi (`hi`) uses Devanagari script

## Usage in Code

### Go Backend

```go
import "dillmann.com.br/nginx-ignition/i18n"

// Get a message
msg := i18n.M(ctx, i18n.K.CoreAccesslistInUse)

// With parameters
msg := i18n.M(ctx, i18n.K.CoreCacheInvalidStatusCode, 
    i18n.P("value", "abc"),
    i18n.P("min", 100),
    i18n.P("max", 599))
```

### TypeScript Frontend

```typescript
import MessageKey from "@/core/i18n/model/MessageKey.generated"
import { i18n } from "@/core/i18n/I18n"

// Get a message
const text = i18n(MessageKey.CoreAccesslistInUse)

// In JSX
<I18n id={MessageKey.FrontendAuthenticationLoginButton} />
```

Never pass raw key path strings in `<I18n id={...}>` or `i18n(...)` — always use `MessageKey.*` constants
(see **NEVER Use Raw Key Strings in Code** above).

Important frontend notes:
- Always use the `I18N` component. Only use `i18n()` function if the output is required to be a `string`.
- A param must be a type that can be rendered as a string, like a number, boolean, string itself, etc. A param cannot be
  a dynamic value (e.g. a component, function, link, HTML tag and so on). If you need to render a dynamic value in the
  middle of a message, split the message into two or more parts.

## Summary Checklist

When adding or modifying keys:

- [ ] Key prefix matches the folder where it's used
- [ ] Suffix is descriptive but not redundant
- [ ] No dots in the key (use `/` only)
- [ ] Value uses `${var}` for placeholders
- [ ] Key is unique (no duplicates)
- [ ] Key stays on its original line (if editing)
- [ ] Go code references the key via `i18n.K.*` — never a raw string literal
- [ ] Frontend code references the key via `MessageKey.*` — never a raw string literal
- [ ] Run `make generate-i18n` after changes
