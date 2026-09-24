# Security policy

First of all, thank you for the time and effort you put into reviewing nginx ignition and sharing security findings, 
bug reports, suggestions, and similar feedback. Your help improves the project and is greatly appreciated.

## Reporting a vulnerability

We encourage you to report suspected vulnerabilities privately. The maintainer will assess each report and decide
whether it should be handled through the private security advisory process or migrated to a regular bug. Reports
assessed as **Medium**, **High** or **Critical** under CVSS v4.0 will generally be handled as security advisories, but
the maintainer's assessment is final.

Please note that submitting a report does not guarantee acceptance, a particular severity rating, a fix, or any specific
outcome. The maintainer may close, reclassify, or migrate a report after the initial review.

Use [GitHub private vulnerability reporting](https://github.com/lucasdillmann/nginx-ignition/security/advisories/new)
and include:

- The affected version and deployment method, such native and its platform (Windows, Linux, macOS), Docker, or both.
- Any relevant non-default user or permission definitions, such as custom roles, access levels, access lists, or account
  settings.
- Any other non-default environment details or configuration required to reproduce the issue.
- A clear description of the vulnerability and its potential impact.
- Reproduction steps and a working proof of concept.
- Relevant logs or configuration with credentials, tokens, personal information, and other sensitive data redacted.
- Any known workaround or mitigation you have identified.

The maintainer assigns the final CVSS v4.0 rating after the initial review and will explain whether the report will
remain in the private security advisory process or be migrated to a regular bug. Reports may be reclassified after
further triage. If a regular bug is appropriate, the maintainer will ask you to migrate the report (do not migrate it
proactively). If migration is requested, review the report for secrets, personal data, and any proof-of-concept content 
that would be unsafe to publish before creating the public issue.

Until the maintainer explicitly asks you to migrate the report to a regular bug, keep it private and do not open a public
issue or disclose the report. If the report remains in the advisory process, keep it private until a fix and advisory have
been published.

## Supported versions

Security fixes are provided for the latest stable release only. Use the [releases page](https://github.com/lucasdillmann/nginx-ignition/releases) to check the current stable version and available downloads.

| Version               | Supported |
|-----------------------|-----------|
| Latest stable release | Yes       |
| Older releases        | No        |

Users running an older release should upgrade to the latest stable release before reporting an issue that is not present
here.

## AI-assisted reports

AI tools may help organize information, explain results, improve wording, investigate the issue, review evidence, analyze
code, and develop or refine reproduction steps, but it cannot be used end-to-end without a human in the loop as the lead.

The human reporter must remain in the lead, verify the evidence and conclusions, inspect the relevant behavior, reproduce
the scenario, determine its impact, and submit an account they understand. AI-generated suggestions, analysis, and
conclusions must be reviewed and verified by the reporter rather than accepted without human validation. Regardless of
whether AI was used, the reporter is fully responsible for the report, including its accuracy, completeness, technical
claims, and compliance with this policy.

A report generated and submitted end-to-end by AI without a human in the loop as the lead, or without a real,
human-led and human-reproduced security scenario, will be discarded without further review.

## Scope

Security issues in nginx ignition, its official release artifacts, and its official container images are in scope. These
include vulnerabilities in how the project integrates or handles third-party components.

Vulnerabilities that exist solely in nginx, dependencies, container base images, or third-party services should be
reported to their respective maintainers. If nginx ignition's use of an upstream component creates an additional
vulnerability, that issue is in scope and should be reported here.

## Responsible research

Good-faith security research is welcome, but this policy does not promise that the project or its maintainer will not pursue action for bad-faith abuse or other harmful conduct. We ask researchers to:

- Test only systems, accounts, and data they own or are explicitly authorized to test.
- Avoid accessing unrelated user data, introducing persistence, or changing systems outside the scope of the test.
- Use bounded proof-of-concept activity and avoid sustained denial-of-service attacks against shared services.
- Minimize collection and retention of sensitive data.
- Report vulnerabilities privately and allow reasonable time for remediation and coordinated disclosure before public
  disclosure.

## Response and disclosure

For private reports that enter review, the maintainer aims to acknowledge receipt, request additional information when
needed, and keep the reporter informed while the issue is investigated. The target is to provide an initial assessment
typically within 7 days, and a fix or final response typically within 30 days. This is a planning target, not a 
guaranteed response or resolution deadline (timing may vary depending on the complexity of the report).

Confirmed issues will be addressed in a fix and communicated through a release, the changelog, and a GitHub Security
Advisory whenever appropriate. If applicable, disclosure may be coordinated with the reporter to help users understand 
the impact and upgrade before technical details are made public.

## Credit

Reporters are credited by default in the published GitHub Security Advisory, release, and changelog. Credit will use 
the name or GitHub handle provided by the reporter.

**To have credit omitted from both the GitHub Security Advisory and the changelog, the reporter must explicitly request
anonymity before either is published.** Silence or a failure to mention a preference is not considered a request for
anonymity. If you wan't to opt-out, we suggest adding a thread/comment stating it in the Security advisory as soon it 
is opened.

## Bounty, compensations, and payments

nginx ignition is a free, open-source, non-commercial, non-profit homelab project. It does not generate revenue, and no 
compensation, payment, or bounty will be provided for security reports or other contributions.
