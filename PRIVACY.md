# Privacy policy

nginx ignition is a self-hosted application that runs locally as a single binary. It does **not** collect any usage
statistics, tracking, analytics or telemetry, and it shares nothing with the developer or any third party.

## nginx request statistics

nginx ignition may collect nginx request/traffic statistics (such as request rates, response times and a traffic
breakdown by host, domain and upstream server) **only if you enable that feature**. Those statistics:

- are stored locally on the host machine where nginx ignition is running, and
- can be viewed through the web interface or API by users with permission to read traffic statistics, including
  from remote devices.

When this feature is enabled, nginx ignition may also download GeoIP databases (GeoLite2 Country and City) from the
[P3TERX/GeoLite.mmdb](https://github.com/P3TERX/GeoLite.mmdb) project on GitHub to enrich the statistics with
geographic data. These downloads are made anonymously (without authentication) and, as with any HTTPS request, GitHub
may observe your network connection details (for example your IP address), which is governed by GitHub's own privacy
policy.

## Third-party services

Some features, when you use them, may make requests to external, third-party tools:

- **SSL certificates** (for example Let's Encrypt/ACME and DNS providers)
- **VPN integrations** (for example Tailscale and NetBird)
- other configured integrations (for example Docker, Docker Swarm and TrueNAS)

No requests are made to those tools unless you are actively using the corresponding feature. Those third-party tools
have their own privacy policies, which must be read separately, as they may collect and/or store some details.

## Update checks

nginx ignition may call the GitHub API anonymously to query the latest stable version of the project and, when an
update is available, present a warning to the user. This request is made without authentication and does not include
any of your data; as with any HTTPS request, GitHub may observe your network connection details (for example your IP
address), which is governed by GitHub's own privacy policy.

## Data storage

All data is stored locally in an embedded SQLite database or in your own PostgreSQL server, plus configuration and
data directories under your control. nginx ignition does not transfer any information to other networked systems
beyond the requests described above.
