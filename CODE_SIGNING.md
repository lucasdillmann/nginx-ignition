# Code signing policy

Free code signing provided by [SignPath.io](https://about.signpath.io), certificate by
[SignPath Foundation](https://signpath.org).

## What is signed

The following Windows artifacts are Authenticode-signed with a certificate issued by the SignPath Foundation:

- `nginx-ignition.exe` (amd64 and arm64), including the copy shipped inside the Windows ZIP archives and the
  installers;
- the Windows NSIS installers (`nginx-ignition-<version>.windows-<arch>-setup.exe`).

No other artifacts are signed.

## Where signing happens

Signing happens exclusively in the project's tag-triggered CI release workflow
(`.github/workflows/release.yml`). Binaries and installers are built from the source code in this repository inside
GitHub Actions and submitted to SignPath using the official SignPath GitHub Action. Local or manual builds are never
signed.

## Official repository

The only official source repository is <https://github.com/lucasdillmann/nginx-ignition>.

## Code signing roles

- **Authors:** Lucas Dillmann ([@lucasdillmann](https://github.com/lucasdillmann))
- **Reviewers:** Lucas Dillmann ([@lucasdillmann](https://github.com/lucasdillmann))
- **Approvers:** Lucas Dillmann ([@lucasdillmann](https://github.com/lucasdillmann))

## Technical constraints

nginx ignition accepts the technical constraints of the SignPath.io Open Source subscription and does not attempt to
work around them:

- Binary artifacts are built from the source code in this repository in a verifiable way.
- Every release requires manual approval for signing.

If the SignPath Foundation receives a complaint about a violation of the Code of Conduct, the project team will assist
in verification, investigation and root cause analysis.

## Privacy

See our [privacy policy](PRIVACY.md).
