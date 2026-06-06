# nginx-ignition Guide for AI Agents

Project-wide architecture, conventions, and expected behaviour for nginx-ignition. Follow these patterns when adding
features, fixing bugs, or writing tests in any module.

Domain-specific supplements (e.g. `core/notification/AGENTS.md`) extend this document — they do not replace it.

**When i18n is in scope** — adding or modifying message keys, `.properties` files, or user-facing text in Go or
frontend — agents **must** also read and follow [`i18n/AGENTS.md`](i18n/AGENTS.md).

## Architecture

### Module layout

The repository is a Go workspace (`go.work`) with separate modules per layer and feature area:

| Layer               | Path                                                    | Role                                                 |
|---------------------|---------------------------------------------------------|------------------------------------------------------|
| Core business logic | `core/`                                                 | Services, validators, domain models, scheduled tasks |
| Persistence         | `database/`                                             | Repositories, DB models, converters                  |
| HTTP API            | `api/`                                                  | Gin handlers, DTOs, route registration               |
| Feature plugins     | `certificate/`, `integration/`, `vpn/`, `notification/` | Driver/provider implementations in sibling modules   |
| i18n                | `i18n/`                                                 | Properties files and generated message keys          |
| Frontend UI         | `frontend/`                                             | React SPA: domain pages, services, shared components |
| Application         | `application/`                                          | Composition root, server startup                     |

Each domain (user, host, vpn, nginx, notification, …) has matching packages across `core/`, `database/`, and `api/` when
it exposes HTTP or persistence.

### Core package file roles

| File                        | Role                                                                                                    |
|-----------------------------|---------------------------------------------------------------------------------------------------------|
| `service.go`                | Struct definition, constructor, shared helpers                                                          |
| `service_{scope}.go`        | Scoped service methods (see below)                                                                      |
| `commands.go`               | `Commands` interface consumed by other layers; catalog DTOs (`AvailableDriver`, `AvailableProvider`, …) |
| `model.go`                  | Domain entities, enums, domain helper functions                                                         |
| `constants.go`              | Sentinel errors, limits, magic values                                                                   |
| `validator.go`              | Business validation via `validation.ConsistencyValidator`                                               |
| `repository.go`             | Repository interface                                                                                    |
| `provider.go` / `driver.go` | External integration interface (when applicable)                                                        |
| `installer.go`              | DI wiring via `container.Provide` / `container.Run`                                                     |
| `{name}_task.go`            | Scheduled task registration and `Run`/`Schedule` (delegates to `Commands`)                              |

### Service decomposition

When a service grows beyond a handful of methods, split it across scoped files — **never** consolidate everything into a
monolithic `service.go`:

| Pattern             | Example                                                              |
|---------------------|----------------------------------------------------------------------|
| Single-scope module | `core/user/service.go` — all methods in one file when small          |
| Multi-scope module  | `core/nginx/service.go` + `service_stats.go` + `service_metadata.go` |
| Multi-scope module  | `core/notification/service_inbox.go` + `service_publish.go` + …      |

Each scoped source file gets a matching scoped test file (see Test Organization).

### Commands interface

`commands.go` defines the public contract for a domain. Other modules depend on `Commands`, not on the concrete
`service` struct.

Catalog DTOs for pluggable drivers/providers live here, mirroring `vpn.AvailableDriver`:

```go
type AvailableDriver struct {
  Name                  *i18n.Message
  ID                    string
  ImportantInstructions []*i18n.Message
  ConfigurationFields   []dynamicfields.DynamicField
}
```

### Validation

Business validation belongs in `validator.go`, not in API handlers. Use `validation.ConsistencyValidator` (same pattern
as `core/vpn`, `core/user`, `integration`):

- Collect field errors via `delegate.Add(field, message)`
- Return via `delegate.Error()` at the end
- Validate `dynamicfields` parameters when a driver/provider defines `ConfigurationFields`

### DI / installer pattern

Each module exposes `Install() error` in `installer.go`:

```go
func Install() error {
    return container.Provide(newCommands)
}

func newCommands(deps ...) (*service, Commands) {
  svc := newService(deps...)
  return svc, svc
}
```

Complex modules chain sub-installers and startup hooks:

```go
container.Run(registerStartup, registerScheduledTask, registerShutdown)
```

The `database/installer.go` registers repositories; `core/installer.go` and `api/installer.go` compose all domain
installers.

### Scheduled tasks

Follow `core/nginx/log_rotation_task.go`:

| File                  | Responsibility                                               |
|-----------------------|--------------------------------------------------------------|
| `{name}_task.go`      | Scheduler registration, `Run`, `Schedule`, interval constant |
| `{name}_task_test.go` | Tests for schedule interval, `Run` wiring, registration      |

The task delegates to `Commands` — it does not contain business logic itself.

## Test Organization

Violations of these rules cause significant review friction.

### Scoped test files

Test files mirror their source files:

| Source                      | Test                                               |
|-----------------------------|----------------------------------------------------|
| `service.go` (single-scope) | `service_test.go` → `Test_service`                 |
| `service_{scope}.go`        | `service_{scope}_test.go` → `Test_service_{scope}` |
| `validator.go`              | `validator_test.go`                                |
| `{name}_task.go`            | `{name}_task_test.go` → `Test_{name}Task`          |
| `model.go` (pure helpers)   | `model_test.go` when needed                        |
| `converter.go`              | `converter_test.go`                                |

### Forbidden patterns

| Forbidden                                                      | Why                                                          |
|----------------------------------------------------------------|--------------------------------------------------------------|
| `service_test.go` when scoped `service_{scope}.go` files exist | Monolithic — use scoped test files                           |
| Orphan test files without matching source                      | e.g. `delivery_test.go` when source is `service_delivery.go` |
| Merging scoped tests into one file                             | Each scope gets its own test file                            |
| Assertion-bearing tests in `artifacts_test.go`                 | Helpers only                                                 |

### Nested test structure

Three levels — reference `core/user/service_test.go`:

```
Test_service_{scope}(t)          // or Test_service, Test_logRotationTask, Test_Repository
  └── t.Run("{MethodName}", ...)
        └── t.Run("{scenario description}", ...)
              └── ctrl := gomock.NewController(t)  // controller HERE, in leaf subtest
```

Rules:

- One `gomock.NewController` per **leaf** scenario subtest — never share across siblings.
- Top-level name matches the source scope (`Test_service`, `Test_service_inbox`, `Test_listHandler`).
- Second level: method name (`Get`, `Save`, `handle`, `Schedule`).
- Third level: behavioural scenario in plain English.

### Test helpers (`artifacts_test.go`)

Reusable fixtures only — no test functions that assert behaviour:

- `core/user/artifacts_test.go` — `newUser()`, `newSaveRequest()`
- `api/user/artifacts_test.go` — `newUserPage()`, sample DTOs
- `database/user/artifacts_test.go` — `newUser()` for repository tests

### Database repository tests

```
Test_Repository(t)
  └── testutils.RunWithMockedDatabases(t, runRepositoryTests)

runRepositoryTests(t, db)
  └── t.Run("{MethodName}", ...)
        └── t.Run("{scenario}", ...)
```

Use helpers from `database/{domain}/artifacts_test.go`. Both SQLite and PostgreSQL are exercised via
`testutils.RunWithMockedDatabases`.

### API handler tests

Structure: `Test_{handlerName}(t)` → `t.Run("handle", ...)` → scenario subtests.

Set the authenticated subject via middleware:

```go
engine.Use(func (ginContext *gin.Context) {
  ginContext.Set("ABAC:Subject", &authorization.Subject{User: subject})
  ginContext.Next()
})
```

One gomock controller per leaf scenario, same as service tests.

## API Conventions

Handlers live in `api/{domain}/`. They are thin — validation and business rules belong in `core/{domain}/validator.go`
and service methods.

### Authentication and authorization

- Read the current user via `authorization.CurrentSubject(ctx)` (see `api/common/authorization/subject.go`).
- Register routes with `authorizer.AllowAllUsers(...)` when any authenticated user may access their own scoped data;
  enforce ownership by passing `userID` to core commands.
- Admin-only routes use permission-based authorization instead.

### List endpoints — pagination only

List handlers accept **only** standard pagination parameters via `pagination.ExtractPaginationParameters(ctx)`:

- `pageSize`
- `pageNumber`
- `searchTerms`

This matches `api/user`, `api/host`, `api/cache`, and other list handlers.

**Do not** add domain-specific filter query parameters to public list APIs unless explicitly required. Filtering beyond
search belongs in core/repository if ever needed.

### Route ordering — static before `/:id`

Register static path segments **before** parameterized routes to avoid Gin shadowing:

```go
group.GET("/available-providers", ...) // static first
group.GET("/unread-count", ...)
group.GET("", ...) // list

byIDPath := group.Group("/:id")
byIDPath.GET("", ...)
byIDPath.PUT("", ...)
byIDPath.DELETE("", ...)
```

### Thin handlers

Handlers bind JSON, convert DTOs, call `Commands`, return responses. They do **not**:

- Run business validation (core validator handles this)
- Contain i18n resolution logic
- Merge sensitive fields (service layer handles this)

Converters deserialize JSON to domain types; invalid values are caught by the core validator.

### Sensitive configuration fields

On read, strip sensitive provider/driver parameters via `dynamicfields.RemoveSensitiveFields` in the service layer. On
save, merge without overwriting unchanged secrets:

```go
configuration.Parameters = dynamicfields.MergeSensitiveFields(
  configuration.Parameters,
  existing.Parameters,
  provider.ConfigurationFields(ctx),
)
```

Always assign the return value — `MergeSensitiveFields` returns the merged map.

## Frontend

The UI is a React 19 + TypeScript SPA in `frontend/`, built with Vite and Ant Design. Domain screens live under
`frontend/src/domain/`; shared infrastructure under `frontend/src/core/`.

**When i18n is in scope** — read [`i18n/AGENTS.md`](i18n/AGENTS.md) for key naming (`frontend/{domain}/` prefix),
generation workflow, and `<I18n>` usage.

### Module layout

| Path                             | Role                                                        |
|----------------------------------|-------------------------------------------------------------|
| `frontend/src/domain/{domain}/`  | List/form pages, domain services, gateways, models, actions |
| `frontend/src/core/components/`  | Reusable UI (DataTable, shell, forms, access control, …)    |
| `frontend/src/core/apiclient/`   | `ApiClient` (fetch), response helpers                       |
| `frontend/src/core/i18n/`        | `I18n` component, context, generated `MessageKey`           |
| `frontend/src/domain/Routes.tsx` | Route + sidebar menu registration                           |

Each domain with a management UI typically includes:

| File                 | Role                                    |
|----------------------|-----------------------------------------|
| `{Name}ListPage.tsx` | Paginated list via `DataTable`          |
| `{Name}FormPage.tsx` | Create (`/new`) / edit (`/:id`) form    |
| `{Name}Gateway.ts`   | HTTP calls to `/api/{name}`             |
| `{Name}Service.ts`   | Response unwrapping, orchestration      |
| `model/*.ts`         | TypeScript interfaces matching API JSON |
| `actions/*Action.ts` | Confirm → API → toast flows             |

### React conventions

The codebase uses **class components** — not hooks.

| Pattern            | Example                                                 |
|--------------------|---------------------------------------------------------|
| List page          | `React.PureComponent` + `DataTable` ref                 |
| Form page          | `React.Component` + `formRef` + `this.state.formValues` |
| Destructive action | Singleton class exported as `new DeleteXAction()`       |

```tsx
// ✅ CORRECT — shell title/actions in componentDidMount
componentDidMount() {
    AppShellContext.get().updateConfig({
        title: MessageKey.CommonUsers,
        actions: [{ description: MessageKey.FrontendUserNewButton, onClick: "/users/new" }],
    })
}

// ❌ DO NOT USE, EVER! — functional components with hooks (not used in this project)
export default function UserListPage() {
    const [data, setData] = useState(...)
}
```

### API integration

Gateway → Service → Page. No axios; no generated OpenAPI client.

```typescript
// Gateway — query params only for list pagination/search
async getPage(pageSize ? : number, pageNumber ? : number, searchTerms ? : string) {
    return this.client.get(undefined, undefined, { pageSize, pageNumber, searchTerms })
}

// Service — unwrap payloads
async list(...args) {
    return this.gateway.getPage(...args).then(requireSuccessPayload)
}
```

- Base path per domain: `new ApiClient("/api/users")`.
- `ApiClient` sends `Accept-Language` from `I18nContext`.
- Use `requireSuccessPayload`, `requireNullablePayload` (404), `requireSuccessResponse` consistently.

### List pages — pagination and search only

List UIs must mirror backend list API conventions (see **List endpoints — pagination only** above):

- `pageSize`, `pageNumber`, `searchTerms` — **no extra filter query params** on the frontend.
- Use `DataTable` with a stable `id` (persists user preferences per table).
- Pass `dataProvider={(pageSize, pageNumber, searchTerms) => service.list(...)}`.

```tsx
// ✅ CORRECT
<DataTable
    id="users"
    columns={this.buildColumns()}
    dataProvider={(pageSize, pageNumber, searchTerms) =>
        this.service.list(pageSize, pageNumber, searchTerms)
    }
    rowKey={item => item.id}
/>

// ❌ AVOID — domain-specific filters on list fetch
dataProvider = { (pageSize, pageNumber, enabledOnly) => ... }
```

Search is debounced in `DataTableHeader`; changing search resets to page 0.

### Routing and menu

Register routes in `frontend/src/domain/Routes.tsx`:

- **Static paths before parameterized routes** (same rule as Gin).
- `menuItem` → sidebar entry; `activeMenuItemPath` for form routes under a list.
- `requiresAuthentication: false` for `/login`, `/onboarding`; `fullPage: true` for pages without shell.

```typescript
// Static before :id
{
    path: "/certificates/new", 
    activeMenuItemPath: "/certificates",
    ...
},
{
    path: "/certificates/:id", 
    activeMenuItemPath : "/certificates",
    ...
},
{
    path: "/certificates", 
    menuItem: { ... }, 
    ...
},
```

Navigate imperatively: `navigateTo("/hosts/new")`, read params: `routeParams().id`, query: `queryParams()`.

### Forms and validation

- Ant Design `Form` with `FormLayout.FormDefaults` / `FormLayout.LabeledItem`.
- Hold editable state in `this.state.formValues`; sync via `onValuesChange`.
- On save error, parse `consistencyProblems` from API body:

```typescript
if (error instanceof UnexpectedResponseError) {
    const validationResult = ValidationResultConverter.parse(error.response)
    if (validationResult != null) this.setState({ validationResult })
}
```

- Per field: `validateStatus={validationResult.getStatus("name")}` and `help={validationResult.getMessage("name")}`.
- Driver/provider config: render `DynamicInput` for each `configurationFields` entry; merge `parameters` on change.
- Sensitive fields: backend strips on read; frontend sends full form — backend merges secrets (same as API conventions).

### Access control

- Wrap list pages in `<AccessControl requiredAccessLevel={READ_ONLY} permissionResolver={...} />`.
- Gate write actions with `isAccessGranted(READ_WRITE, ...)` → disable shell buttons or show `AccessDeniedModal`.
- Form pages without list wrapper: return `<AccessDeniedPage />` when read access missing.

### Styling

- Co-located **plain `.css`** files imported in the component — **not** CSS modules.
- Global variables in `frontend/src/index.css`; theme toggle sets `data-theme` on `<html>`.
- Use `themedColors()` for semantic icon/button colors.
- Prefer Ant Design layout primitives (`Flex`, `Form`, `Table`) before custom CSS.

### i18n (frontend)

See [`i18n/AGENTS.md`](i18n/AGENTS.md) for keys and generation. Quick rules:

- Import `MessageKey` from `core/i18n/model/MessageKey.generated`.
- Prefer `<I18n id={MessageKey.FrontendUserListSubtitle} />` in JSX.
- Use `i18n()` only when a plain string is required (input placeholder, dynamic confirmation text).
- **Never** pass raw key path strings — always use generated `MessageKey.*` constants.
- Key prefix: `frontend/src/domain/user/` → `frontend/user/...`.

```tsx
// ✅ CORRECT
<I18n id={MessageKey.FrontendUserListSubtitle} />
const text = i18n(MessageKey.FrontendUserNewButton)

// ❌ WRONG — raw key strings are forbidden
<I18n id="frontend/user/list-subtitle" />
const text = i18n("frontend/user/new-button")
```

### DTO / types

- Define request/response interfaces in `domain/{name}/model/` matching `api/{name}` JSON field names.
- `PageResponse<T>` shape: `pageSize`, `pageNumber`, `totalItems`, `contents`.
- Complex forms may use separate form types + converters (e.g. `HostFormValues` ↔ `HostRequest` via `HostConverter`).

### Build, lint, format

```bash
cd frontend && pnpm install
pnpm run start          # dev server :8080, proxies /api → :8090
pnpm run build          # tsc && vite build → frontend/build/
pnpm run check          # prettier --check + eslint
```

From repo root: `make .frontend-lint`, `make .frontend-format`, `make .frontend-build` (build requires
`make .generate-i18n-files`).

### Testing

There is no frontend unit/component test suite today. Validate via `pnpm run check` and manual testing. Do not add
Vitest unless the project adopts it repo-wide.

### Forbidden patterns

| Forbidden                      | Why                                                    |
|--------------------------------|--------------------------------------------------------|
| Hooks in new pages             | Project uses class components throughout               |
| Extra list filter query params | Must match backend list API contract                   |
| Pre-translated user strings    | Use `MessageKey` + `<I18n>`                            |
| Raw i18n key path strings      | Use `MessageKey.*` — never `"frontend/..."` literals   |
| CSS modules                    | Project uses plain co-located CSS                      |
| Business validation only in UI | Server returns `consistencyProblems`; UI displays them |

## Go Code Style (Go 1.26)

### Use `new(expr)` for pointer fields

```go
// ✅ CORRECT
submission.LastAttemptAt = new(time.Now())
request.Password = new("password123")

// ❌ AVOID
now := time.Now()
submission.LastAttemptAt = &now
```

### Merge error and nil checks

```go
// ✅ CORRECT
user, err := s.repository.FindByID(ctx, id)
if err != nil || user == nil {
    return nil, err
}
```

### General rules

- No unnecessary comments — code should be self-explanatory
- Fix lint issues; do not dismiss them as pre-existing
- Match surrounding naming: full words, no abbreviated variables (`user`, not `u`)
- Keep diffs minimal and scoped to the task

## i18n

**When i18n is in scope, read [`i18n/AGENTS.md`](i18n/AGENTS.md) first.** That guide is authoritative for:

- **Key naming** — path prefix matches code location (`core/user/` → `core/user/not-found`); `common/` for keys shared
  across folders
- **Properties format** — `key/path/suffix=Value`, `/` separators only (no dots), one key per line
- **Add-key workflow** — edit `messages_en.properties`, then run `make .generate-i18n-files` to regenerate
  `keys.generated.go`, `en.generated.go`, and `MessageKey.generated.ts`
- **Value rules** — `${variable}` placeholders, sentence case, unique keys and values
- **Localization** — keep all locale files in sync with `messages_en.properties`
- **Usage in code** — `i18n.M(ctx, i18n.K....)` in Go; `I18n` component / `MessageKey` in frontend

Quick rules:

- Key prefix must match the folder where the key is used (`core/user/` → `core/user/not-found`)
- Run `make .generate-i18n-files` after adding keys
- Use `i18n.M(ctx, i18n.K....)` in Go; use the `I18n` component / `MessageKey.*` in frontend
- **Never** pass raw key path strings — always use generated constants (`i18n.K.*` in Go, `MessageKey.*` in frontend)

```go
// ✅ CORRECT — Go
i18n.M(ctx, i18n.K.CoreUserNotFound)
i18n.DetachedMessage{Key: i18n.K.CoreNotificationCategoryCertificateRenewed}

// ❌ WRONG — Go raw key strings are forbidden
i18n.M(ctx, "core/user/not-found")
i18n.DetachedMessage{Key: "core/notification/category/certificate-renewed"}
```

```tsx
// ✅ CORRECT — frontend
<I18n id={MessageKey.FrontendUserListSubtitle} />
const text = i18n(MessageKey.FrontendUserNewButton)

// ❌ WRONG — frontend raw key strings are forbidden
<I18n id="frontend/user/list-subtitle" />
const text = i18n("frontend/user/new-button")
```

Use `i18n.Static("...")` only for non-localized text (test fixtures, dynamic user input) — never for keys
defined in `.properties` files.

### DetachedMessage at async boundaries

When a producer cannot resolve the recipient's language at call time (notifications, async events), pass
`i18n.DetachedMessage` values — not pre-translated strings. Resolution happens later per recipient context.

## Database

### Migrations

- Scripts live in `database/common/migrations/scripts/postgres/` and `database/common/migrations/scripts/sqlite/`
- Numbered sequentially: `NNN_description.up.sql`
- **Both** postgres and sqlite variants are required for every migration
- Keep schema changes in sync across dialects

### Repository layer

Each domain has `database/{domain}/`:

| File            | Role                                        |
|-----------------|---------------------------------------------|
| `model.go`      | DB row structs (may differ from core model) |
| `converter.go`  | Core ↔ DB conversion                        |
| `repository.go` | Implements core `Repository` interface      |

Repository tests use `testutils.RunWithMockedDatabases` to verify behaviour against both SQLite and PostgreSQL.

## Summary Checklist

When working on any module:

- [ ] Business logic in `core/`, persistence in `database/`, HTTP in `api/`
- [ ] Large services split into `service_{scope}.go`; matching scoped test files
- [ ] `Commands` interface in `commands.go`; validation in `validator.go`
- [ ] Test structure: top-level scope → method → scenario; gomock controller in leaf subtest
- [ ] `artifacts_test.go` for helpers only — no assertion tests
- [ ] Database tests: `Test_Repository` → `RunWithMockedDatabases` → method → scenario
- [ ] API tests: `Test_{handler}` → `handle` → scenario; set `ABAC:Subject` in middleware
- [ ] List APIs use only `pageSize`, `pageNumber`, `searchTerms`
- [ ] Static routes registered before `/:id` groups
- [ ] Handlers thin; validation in core
- [ ] Use `new(expr)` for pointer assignment
- [ ] i18n changes → read `i18n/AGENTS.md`
- [ ] i18n keys follow folder-prefix convention (see `i18n/AGENTS.md`)
- [ ] i18n keys via `i18n.K.*` (Go) or `MessageKey.*` (frontend) — never raw strings
- [ ] Migrations in both postgres and sqlite
- [ ] Fix lint issues before finishing
- [ ] Domain UI in `frontend/src/domain/{domain}/` with Gateway + Service + model types
- [ ] List pages use `DataTable` with only `pageSize`, `pageNumber`, `searchTerms`
- [ ] Routes in `Routes.tsx`: static paths before `/:id`; `menuItem` / `activeMenuItemPath` set
- [ ] Forms use `ValidationResult` from API `consistencyProblems`
- [ ] User-facing text via `<I18n>` and `MessageKey.*` — never raw key strings (see `i18n/AGENTS.md`)
- [ ] Co-located plain CSS; no CSS modules
- [ ] Class components (no hooks) unless project direction changes
- [ ] `make .generate-i18n-files` before frontend build when keys change
- [ ] `pnpm run check` / `make .frontend-lint` before finishing
