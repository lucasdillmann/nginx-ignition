# Configuration properties

## Common configuration scenarios

### Defining a custom authentication token secret

The nginx ignition uses JSON Web Tokens (JWT) to authenticate you and all the users in the application. A JWT is signed
using a private key, allowing the application to check if any received value is trustworthy.

The JWT secret can be defined using the `NGINX_IGNITION_SECURITY_JWT_SECRET` environment variable, which must be
a string value with exactly 64 chars-long. If no value is provided, nginx ignition will still work by generating a
random value every time the app boots (which can be fine, but will force the users to log-in again everytime the
app restarts).

An example of value to the environment variable follows. Please do not use this value in your installation, but rather
generate a value just for you.

```shell
NGINX_IGNITION_SECURITY_JWT_SECRET="e54rVg9NX5moIP6k2xmUwT0bauAG7pvkR7XI7ygJ6jz0T50huvujCdW4ym6mOjAy"
```

### Connecting to a PostgreSQL database

nginx ignition can be used with its embedded database (powered by SQLite), but its strongly recommended that any long-term
installation is used alongside a PostgreSQL database. To do so, the following environment variables must be informed:

- `NGINX_IGNITION_DATABASE_DRIVER` with the database type
- `NGINX_IGNITION_DATABASE_HOST` with the database hostname or IP
- `NGINX_IGNITION_DATABASE_PORT` with the database port
- `NGINX_IGNITION_DATABASE_NAME` with the database name
- `NGINX_IGNITION_DATABASE_SSL_MODE` with the database encryption mode, being either `require` or `disable` (defaults to
  `require` if not informed or an invalid value is provided)
- `NGINX_IGNITION_DATABASE_USERNAME` with the database username
- `NGINX_IGNITION_DATABASE_PASSWORD` with the database password

To exemplify, let's imagine that you have your SGDB running at the IP  `192.168.1.150`. The database name is
`my_custom_db`, username is `my_username` and password `supersecretpassword`. In such scenario, the environment
variables should have the following values:

```shell
NGINX_IGNITION_DATABASE_DRIVER=postgres
NGINX_IGNITION_DATABASE_HOST=192.168.1.150
NGINX_IGNITION_DATABASE_PORT=5432
NGINX_IGNITION_DATABASE_NAME=my_custom_db
NGINX_IGNITION_DATABASE_SSL_MODE=disable
NGINX_IGNITION_DATABASE_USERNAME=my_username
NGINX_IGNITION_DATABASE_PASSWORD=supersecretpassword
```

nginx ignition will create all the required tables, indexes and alike when it boots. In future updates, any changes
will be applied automatically also.

## All configurations properties available

The following configuration properties are available through environment variables. Use them freely to customize
nginx ignition to suit you better, if needed.

| Environment variable                                     | Description                                                                                           | Example      | Default value                                                                 |
|----------------------------------------------------------|-------------------------------------------------------------------------------------------------------|--------------|-------------------------------------------------------------------------------|
| NGINX_IGNITION_SERVER_PORT                               | Port number where the nginx ignition should listen for requests                                       | 1234         | 8090                                                                          |
| NGINX_IGNITION_SERVER_ADDRESS                            | Address/IP where the nginx ignition should listen for requests                                        | 192.168.0.1  | 0.0.0.0                                                                       |
| NGINX_IGNITION_NGINX_BINARY_PATH                         | Path to the nginx's binary that the nginx ignition should use                                         | /bin/nginx   | nginx                                                                         |
| NGINX_IGNITION_NGINX_CONFIG_PATH                         | Path on where the nginx ignition should store the generated nginx's configuration files               | /etc/nginx   | /tmp/nginx-ignition/nginx (`C:\Windows\Temp\nginx-ignition\nginx` on Windows) |
| NGINX_IGNITION_VPN_CONFIG_PATH                           | Path on where the nginx ignition should store the generated vpn configuration files                   | /etc/vpn     | /tmp/nginx-ignition/vpn (`C:\Windows\Temp\nginx-ignition\vpn` on Windows)     |
| NGINX_IGNITION_DATABASE_DRIVER                           | The type of the database, being either `postgres` or `sqlite`                                         | postgres     | sqlite                                                                        |
| NGINX_IGNITION_DATABASE_HOST                             | Hostname or IP of the database server                                                                 | 192.168.0.1  |                                                                               |
| NGINX_IGNITION_DATABASE_PORT                             | Port on where the database is listening for connections                                               | 5432         |                                                                               |
| NGINX_IGNITION_DATABASE_NAME                             | Name of the database to be used                                                                       | 5432         |                                                                               |
| NGINX_IGNITION_DATABASE_SSL_MODE                         | Definition if the connection to the database should be encrypted, being either `require` or `disable` | disable      | require                                                                       |
| NGINX_IGNITION_DATABASE_USERNAME                         | Database username                                                                                     | postgres     |                                                                               |
| NGINX_IGNITION_DATABASE_PASSWORD                         | Database username                                                                                     | postgres     |                                                                               |
| NGINX_IGNITION_DATABASE_SCHEMA                           | Schema name (PostgreSQL only)                                                                         | example      | public                                                                        |
| NGINX_IGNITION_DATABASE_DATA_PATH                        | Folder on where the database file should be stored. Applicable only for the `sqlite` database.        | /opt/example | /tmp/nginx-ignition/data (`C:\Windows\Temp\nginx-ignition\data` on Windows)   |
| NGINX_IGNITION_SECURITY_JWT_SECRET                       | Secret key (64 chars long) for the authentication tokens                                              |              |                                                                               |
| NGINX_IGNITION_SECURITY_JWT_TTL_SECONDS                  | Amount of seconds that an authentication token will be valid before logout by inactivity              | 3600         | 3600                                                                          |
| NGINX_IGNITION_SECURITY_JWT_RENEW_WINDOW_SECONDS         | Amount of seconds that an authentication token will be automatically renewed before its expiration    | 900          | 900                                                                           |
| NGINX_IGNITION_SECURITY_JWT_CLOCK_SKEW_SECONDS           | Amount of seconds that the token's dates can variate from the server dates                            | 60           | 60                                                                            |
| NGINX_IGNITION_SECURITY_USER_PASSWORD_HASHING_ALGORITHM  | Which algorithm should be use to hash the user's passwords                                            | SHA-512      | SHA-512                                                                       |
| NGINX_IGNITION_SECURITY_USER_PASSWORD_HASHING_SALT_SIZE  | The amount of random bytes that should be appended to the user's passwords (improves security)        | 64           | 64                                                                            |
| NGINX_IGNITION_SECURITY_USER_PASSWORD_HASHING_ITERATIONS | How many times the passwords should be hashed (improves security)                                     | 1024         | 1024                                                                          |
| NGINX_IGNITION_FRONTEND_CODE_EDITOR_API_KEY              | Custom Codeium API key for the frontend's code editors (optional)                                     |              |                                                                               |
| NGINX_IGNITION_HEALTH_CHECK_ENABLED                      | Defines if the health check endpoints should be enabled or not                                        | false        | true                                                                          |

## Configuration file

nginx ignition can be configured using a configuration properties file. The app will try to read the
configuration from the file specified by the `NGINX_IGNITION_CONFIG_FILE_PATH` environment variable, from the 
`--config` command line argument or from the `nginx-ignition.properties` file in the current working directory.

Examples of the configuration file:
- [Linux](../dist/linux/nginx-ignition.properties)
- [macOS](../dist/macos/nginx-ignition.properties)
- [Windows](../dist/windows/nginx-ignition.properties)
