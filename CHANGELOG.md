# CHANGELOG

## 2.35.0
- Development pipeline and workflow improvements
- Security fixes and updates

## 2.34.0

- NetBird VPNs are now supported alongside Tailscale, allowing you to expose any host as a peer/subdomain in your
  private networks.
- Tailscale integration has been improved and no longer requires HTTPS certificate support to be enabled on the Tailnet
  coordinator server.
- Improved initial UI page load times.
- Improved feedback on the login and onboarding pages when a field is filled out incorrectly.
- TrueNAS integration now supports the new WebSocket-based APIs. Existing integrations will continue to use the
  deprecated REST APIs, but can be migrated by editing the integration and disabling the legacy switch.
- Fixed an issue where the Stats page would break if statistics were newly enabled and no data was available yet.
- Improved Docker image configuration, which now reports which paths should be made persistent using a volume.
- Improved automatic downloads of GeoIP databases (used when stats are enabled), which could previously fail in some
  scenarios with an unexpected EOF error.
- Other minor improvements and bug fixes.

## 2.33.0

- Nginx Ignition now supports Two-Factor Authentication (2FA)
- New onboarding page, reworked for a better experience
- Performance optimizations, security updates, and other minor fixes and improvements

## 2.32.0

- Nginx Ignition now has integrated traffic statistics
  - Real-time insights into server performance
  - Metrics for request rates, response times, and bandwidth
  - Traffic breakdown by host, domain, and upstream servers
- Improvements for the multi-language support
  - Corrected the Hindi translation
  - Added French and German translations
  - Removed Western Punjabi due to Right-to-Left compatibility issues
- Dependencies and security updates
- Other minor improvements and fixes

## 2.31.2

- Minor fixes and improvements for the logs and tables

## 2.31.1

- Minor fixes and improvements for the logs and tables
- Dependencies and security updates

## 2.31.0

-   Improvements to the Logs page, including a new text viewer and search capabilities with support for viewing
  surrounding lines of a search match.
-   Improvements to Ignition's Tables
    -   Tables now support customizable settings via a new options menu. You can configure the default page size for all
      tables and/or configure the app to remember the last page size used (globally or on a per-table basis).
    -   Tables can now persist your search terms and selected page number within the same session. This allows you to
      return to the exact page and/or filtered results after viewing a row's details.
-   Other minor improvements and bug fixes.

## 2.30.0

- On host routes of the type integration, Nginx Ignition now allows the customization of the base URI on the target
  service

## 2.29.1

- Fix for the SSL certificate import feature failing when the certification chain file ends with an empty line (thanks
  @dpaessler for the [bug report](https://github.com/lucasdillmann/nginx-ignition/discussions/82))

## 2.29.0

- Nginx Ignition now supports multiple languages: Brazilian Portuguese, English, Chinese (Simplified), Hindi, Spanish,
  Vietnamese, Russian, Bengali, Japanese, and Punjabi.
  - This is currently a **beta feature**. If you find any translation issues or have suggestions for improvements,
    please feel free to [raise an issue](https://github.com/lucasdillmann/nginx-ignition/issues).
  - If you'd like to see a new language added, please [start a discussion](https://github.com/lucasdillmann/nginx-
    ignition/discussions).
- General UI/UX improvements, including an updated login screen.
- Security updates.
- Other minor improvements and bug fixes.

## 2.28.0

- Nginx Ignition is now available for Windows (amd64 and arm64). A ZIP file is available for download below; usage
  instructions are included within the archive.
- Updated the README file to be less technical and more concise.
- Added support for 19 new DNS providers for Let's Encrypt challenges: AlibabaCloud ESA, Alwaysdata, Anexia CloudDNS,
  Beget.com, 35.com, EdgeCenter, Gigahost.no, Gravity, Hostinger, Hosting.nl, Ionos Cloud, ISPConfig 3, ISPConfig 3
  DDNS, JD Cloud, Neodigit, Octenium, United Domains, Virtualname, and webnames.ca.
- Other minor improvements and bugfixes.

## 2.27.1

- Fixed an issue in the systemd service file that prevented Ignition from starting when installed via Linux packages
  (.deb, .rpm, etc; thanks [DustPhyte](https://www.reddit.com/user/DustPhyte/) for [the
  report](https://www.reddit.com/r/selfhosted/comments/1nabcy6/comment/nydk8qa/))

## 2.27.0

- Improvement of the SSL certificates UI
- Added the ability to customize the index file (e.g, index.html) when using static files host routes
- Internal code quality improvements
- Other minor improvements and bug fixes

## 2.26.0

- Native support for the nginx content caching
- Optimization of the ignition's performance and memory usage
- Other minor improvements and bug fixes

## 2.25.1

- Minor fixes and improvements for the Docker Swarm integration

## 2.25.0

- Native Docker integration now supports Swarm clusters. When Swarm mode is enabled in the integration settings,
  Ignition will look for deployed services and their host/ingress port bindings instead of the underlying containers
  (thanks @Kegelcizer for the suggestion!).
- Also in the Docker integration, Ignition now offers the option to use the container name as the identifier instead of
  the container ID. For containers managed by third-party tools or those that are constantly recreated, this option
  prevents the integration from failing to locate the container after a recreation.
- Improved application stability and consistency.
- Other minor bug fixes and improvements.

## 2.24.0

- Docker integration now shows and allows the selection of ports exposed in the container (previously limited to ports
  exposed on the host; thanks @Kegelcizer for the suggestion!)
- On the hosts list page, domain names are now links that open in a new tab (thanks @Kegelcizer for the suggestion!)
- Other minor improvements and bug fixes.

## 2.23.0

- The restriction to either use `root` or `nginx` as the nginx runtime user was removed. Now, the settings page allows
  to inform any user as the runtime user.
- Fix for the `nginx-ignition.properties` file where a wrong key was used for the nginx binary path (thanks @maslke for
  catching and [sharing the problem](https://github.com/lucasdillmann/nginx-ignition/discussions/47)!)

## 2.22.0

- The settings page was updated to include new, advanced configuration properties (such as customizable buffer sizes and
  the ability to input custom nginx configuration directives), enabling fine-tuning for specific use cases.
- Minor improvements and bug fixes.

## 2.21.0

- nginx ignition now includes health check endpoints to enable automated monitoring of application status. Refer to
  [this documentation page](https://github.com/lucasdillmann/nginx-ignition/blob/main/docs/health-checks.md) for details
  on how to configure or disable them.
- Added a new [docker-compose.yml](https://github.com/lucasdillmann/nginx-ignition/blob/main/docker-compose.yml) file
  with a suggested deployment configuration for nginx ignition alongside a PostgreSQL database with health checks
  enabled.
- Other minor improvements and bug fixes.

## 2.20.2

- Security fixes and updates

## 2.20.1

- Fix for the database backups not including the ignition's users data

## 2.20.0

- The upload of a custom SSL certificate (issued by a third-party) now allows it to be pasted as PEM-encoded text as an
  alternative to the previously available PEM file selection (thanks @Volt-hf for [the
  suggestion](https://github.com/lucasdillmann/nginx-ignition/discussions/35))
- Other minor improvements and bug fixes

## 2.19.0

- nginx ignition now has a native integration with the Tailscale (Tailnet) network, enabling a simplified way of
  exposing your hosts on your VPNs for remote access. Support for more networks are coming soon.
- Security updates and fixes
- Other minor improvements and bug fixes

## 2.18.1

- Support for multiple Docker and TrueNAS integrations/hosts (previously limited to one of each; [thanks Th3Smok3y for
  the suggestion](https://github.com/lucasdillmann/nginx-ignition/discussions/28))
- Fix for the database migrations failing when using PostgreSQL (thanks @smokey007 for the [bug
  report](https://github.com/lucasdillmann/nginx-ignition/issues/31))
- Other minor improvements and bug fixes

## 2.17.1

- Security fixes and updates

## 2.17.0

- nginx ignition can now be installed on Debian, Ubuntu, Arch Linux, Red Hat, Rocky Linux, Alpine Linux, OpenWrt, and
  other distributions using native packages (.deb, .rpm, .apk, etc.). Download links for amd64 and arm64 are available
  below.
- The nginx Stream, Lua, and JS modules are now optional. nginx ignition can detect if they are available, display a
  warning in the UI if they are not, and still start the nginx server. By making them optional, ignition can be
  installed on a wider range of Linux distributions.
- Static response routes now support large response payloads, which could previously cause the nginx server to fail to
  start.
- Other minor improvements and bug fixes.

## 2.16.0

- Added support for 153 more DNS providers for the Let's Encrypt certificate challenges, bringing the total to 157
  providers (such as AWS Route 53, Azure, Google Cloud, DigitalOcean, Vercel, cPanel, DirectAdmin, and many more)
- Other minor fixes and improvements

## 2.15.1

- nginx ignition can now be executed directly on Linux (amd64/arm64) and macOS (arm64) without Docker. ZIP files for
  download are available below (installation instructions included) with more options (like .deb, .rpm and alike) coming
  soon.
- Other minor fixes and improvements

## 2.14.0

- Improved support for streams. Now nginx ignition is able to configure domain-based routing using the SNI (Server Name
  Indication) from the TLS protocol, forwarding a request to a custom set of backing servers using the requested domain
  name as the qualifier. Support for weights and circuit breakers is also now available.
- Other minor improvements and bugfixes

## 2.13.1

- Minor fixes and improvements for the new backup and export features ([released on the 2.13.0
  version](https://github.com/lucasdillmann/nginx-ignition/releases/tag/2.13.0))

## 2.13.0

- New option to download a copy of the ignition's database, enabling simpler backups for disaster recoveries
- New option to export the nginx configuration files that the ignition generates, allowing them to be reviewed and even
  used into a standalone nginx deployment
- Minor bugfixes and improvements
- Internal project structure and code improvements

## 2.12.0

- Improvements for the static files support. Now is possible to customize if the directory listing should be enabled or
  not.
- Other minor improvements and UI tweaks

## 2.11.1

- Fix for Let's Encrypt SSL certificates not being issued when using Cloudflare DNS (where the TTL wasn't being set)

## 2.11.0

- Automatic retry of the nginx server startup in the ignition startup, improving the scenario where third-party
  dependencies (like the TrueNAS integration) are being initialized alongside ignition and may not be ready yet
- Dependencies and security updates

## 2.10.1

- Dependencies update and other minor fixes

## 2.10.0

- Improvement of the login and onboarding pages
- Fix for the log being returned with the lines in the reverse order
- Other minor fixes and improvements

## 2.9.1

- Fix for host copy action not working as expected

## 2.9.0

- User permissions and access control improvements
- Minor fixes and improvements

## 2.8.0

- New dark mode theme
- Minor fixes and improvements

## 2.7.0

- New "Directory" host route type able to serve static files with directory listing
- Minor improvements and bug fixes

## 2.6.0

- Support for nginx's streams, enabling the proxy of raw TCP, UDP and Unix sockets
- Improvement of the feedback for when a new version is available
- Minor improvements and bug fixes

## 2.5.0

- New settings option for the nginx's runtime user
- Minor improvements and bug fixes

## 2.4.2

- Dependencies and security updates

## 2.4.1

- Fix for the password reset procedure not working after the codebase migration to Golang
- Dependencies and security updates

## 2.4.0

- Better feedbacks and details when the nginx fails to start, stop or reload
- Better UI error handling
- Other minor improvements and bugfixes

## 2.3.1

- Fix for the browser caching definitions
- Security and dependencies updates
- Minor internal improvements

## 2.3.0

- UI/UX improvements

## 2.2.0

- New code editor with AI code completion (powered by Codium) for the host route's static responses and JavaScript/Lua
  source code (free with no API key required, but can be informed if you have one; check the configuration docs for more
  details)
- Improvement of the internal context management
- Improvement of the memory management settings
- Minor bugfixes and improvements

## 2.1.1

- Minor improvements of the memory management
- Fixes the scenario where a host couldn't be deleted
- Fixes the scenario where the nginx wasn't able to start when using a static response with a body that contained a
  double quote

## 2.1.0

- UI performance improvements
- Code and development dependencies upgrades
- Minor bugfixes

## 2.0.2

- Minor improvements and bugfixes

## 2.0.1

- Fixes the scenario where a host couldn't be updated if it contains any static response route with no custom headers
- Fixes the SSL certificates that couldn't be issued if the domain name contains a wildcard
- Fixes the enable/disable toggle not working in the host's list page
- Other minor improvements

## 2.0.0

Major rewrite of the application code, migrating from Kotlin to Go. This release is more an internal one that prepares
the project for its future iterations, also consuming way less memory and CPU but still keeping all of its features.

Some minor and simple breaking changes were introduced, please check [this migration
guide](https://github.com/lucasdillmann/nginx-ignition/blob/main/docs/migration-guide.md) for more details. Original
source code in Kotlin is available at [this public archive repository](https://github.com/lucasdillmann/legacy-nginx-
ignition) (or in the [git commit history](https://github.com/lucasdillmann/nginx-ignition/tree/1.7.1)).

## 1.7.1

- Minor improvements and fixes

## 1.7.0

- Option to disable a host's route without the need to remove it
- New troubleshooting documentation with instructions like how to reset a password
- Upgrade of the project's dependencies (React, Gradle and other libs/SDKs)
- General optimizations, simplifications and improvements (like the replacement of the Ktor framework by the Java's
  native HTTP server)
- Other minor improvements and bugfixes

## 1.6.0

- Support for JavaScript and Lua code as route handlers
- Other minor improvements and bugfixes

## 1.5.1

- Fix for the H2 database migrations failing while creating the access lists' new tables and constraints
- Other minor improvements and bugfixes

## 1.5.0

- Improvement of the host's routes settings with new configuration options
- Access lists for easy control of who can access what by using source IP address checks and/or username/password
  authentication
- Other minor improvements and bugfixes

## 1.4.0

- New home page, featuring a quick start guide and some useful information
- Native Docker integration
- Improvement of the login flow and feedback for when the session expires
- Better error handling and feedback through empty state pages
- Possibility to reset the settings to the default values
- Other minor improvements and bug fixes

## 1.3.0

- New settings page for the nginx server with definition for the maximum body/upload size, timeouts, server tokens, log
  level and more
- Support for automatic (scheduled) log rotation
- Support for global bindings (port listeners) with the possibility to override them in each host
- UX improvements and other minor fixes

## 1.2.0

- Support for Azure DNS, Google Cloud DNS and Cloudflare DNS as DNS challenge providers for the Let's Encrypt SSL
  certificates
- Support for TrueNAS integrations where the apps are exposed in an address other than the TrueNAS console/API
- Improvement of the application startup times
- UX improvements and other minor fixes

## 1.1.0

- Native integration with TrueNAS Scale
- Other minor improvements and bugfixes

## 1.0.0

Initial release with the initial set of features

- Multiple nginx virtual hosts, each one with its customized set of domain, routes and bindings (port listeners)
- Each host route can act as a proxy, redirection or reply with a static response
- SSL certificates (Let's Encrypt, self-signed or bring your custom one) with automatic renew (when applicable)
- Server and virtual hosts access and error logs
- Multiple users with role-based access control (RBAC)
