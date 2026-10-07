# API tokens

API tokens let scripts, CI jobs and other external tools authenticate against the nginx ignition API without sharing
your username and password.

Each token belongs to the user who created it and carries that user's permissions, so a token can do exactly what its
owner can do. Tokens can be created and revoked at any time from the web interface, and every request made with a token
is subject to the same attribute-based access control (ABAC) as the web interface.

> API tokens are available since version **2.47.0**.

## Creating a token

1. Open the user menu in the top-right corner of the web interface.
2. Switch to the **API Tokens** tab.
3. Click **New token**.
4. Give the token a name and, optionally, an expiration date.
5. Click **Create token** and copy the generated token.

Some things to keep in mind:

- The **name** must be unique among your own tokens and can have up to 256 characters. Use something descriptive
  (like `ci-deploy` or `home-assistant`) so you can tell tokens apart later.
- The **expiration date** is optional. When left empty, the token never expires. When set, it must be in the future and
  the token stops working as soon as that moment passes.
- The token value is displayed **only once**, right after creation. nginx ignition doesn't store it, so it
  cannot be recovered later. If you lose it, revoke the token and create a new one.

## Using a token

Send the token in the `Authorization` header using the `Bearer` scheme, exactly like a login session token:

```shell
curl -H "Authorization: Bearer YOUR_API_TOKEN" \
     http://localhost:8090/api/hosts
```

Requests using an invalid, revoked or expired token are answered with `401 Unauthorized`.

> API tokens authenticate API calls only. They cannot be exchanged for a web interface session, so logging into the web
> interface still requires a username and password.

## Permissions

A token inherits the permissions of the user that owns it. There is no way to narrow a token's permissions down
independently of its owner.

If you need a token with a different (usually smaller) set of permissions, create a dedicated **service user** with the
desired access levels and generate the token from that user's account. Creating users requires read-write access to the
`Users` permission, so you'll need an account with that access to set this up.

## Reserved operations

A few operations are reserved for interactive sessions and can't be performed with an API token:

- **Managing API tokens:** an API token can't list, create or revoke tokens, not even its own. Such requests are
  rejected.
- **Logging out:** logging out only invalidates interactive sessions, so the request is rejected when it's made with an
  API token.

Because of this, revoking a token from the web interface is the only way to permanently disable it.

## Token lifecycle

- **Not auto-renewed:** unlike the tokens used by web sessions, API tokens are never automatically extended. The
  `NGINX_IGNITION_SECURITY_JWT_TTL_SECONDS` and `NGINX_IGNITION_SECURITY_JWT_RENEW_WINDOW_SECONDS` properties do not
  apply to them. Only the per-token expiration date matters.
- **Revocation is immediate:** every authenticated request verifies the token against the database, so revoking a token
  stops it from working on the very next call, without waiting for any cache to expire.
- **Disabled users:** disabling a user makes the tokens that user owns stop working, but nothing is deleted or revoked.
  Re-enabling the user makes their tokens work again.

## Notes

- **Upgrading from a previous version:** version 2.47.0 added a claim to the authentication tokens, which makes the access
  tokens issued by earlier versions invalid. Users are therefore logged out once when updating, and need to log in
  again. 
- Tokens are stored in the `user_token` table and are removed automatically along with their owner.
