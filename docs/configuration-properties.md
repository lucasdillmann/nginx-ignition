# Configuration properties 

The following configuration properties are available through environment variables. Use them freely to customize
nginx ignition to suit you better, if needed.

| Environment variable                                   | Description                                                                                        | Example                                    | Default value                                |
|--------------------------------------------------------|----------------------------------------------------------------------------------------------------|--------------------------------------------|----------------------------------------------|
| NGINX_IGNITION_DATABASE_URL                            | Connection URL (JDBC formatted) to the database                                                    | jdbc:postgresql://localhost/nginx_ignition | jdbc:h2:mem:nginx-ignition;DB_CLOSE_DELAY=-1 |
| NGINX_IGNITION_DATABASE_USERNAME                       | Database username                                                                                  | postgres                                   | sa                                           |
| NGINX_IGNITION_DATABASE_PASSWORD                       | Database username                                                                                  | postgres                                   |                                              |
| NGINX_IGNITION_DATABASE_CONNECTION_POOL_MINIMUM_SIZE   | Minimum amount of database connections tha the app will keep open                                  | 1                                          | 1                                            |
| NGINX_IGNITION_DATABASE_CONNECTION_POOL_MAXIMUM_SIZE   | Maximum amount of database connections tha the app will open at any time                           | 10                                         | 10                                           |
| NGINX_IGNITION_SECURITY_JWT_SECRET                     | Secret key (64 chars long) for the authentication tokens                                           |                                            |                                              |
| NGINX_IGNITION_SECURITY_JWT_TTL_SECONDS                | Amount of seconds that an authentication token will be valid before logout by inactivity           | 3600                                       | 3600                                         |
| NGINX_IGNITION_SECURITY_JWT_RENEW_WINDOW_SECONDS       | Amount of seconds that an authentication token will be automatically renewed before its expiration | 900                                        | 900                                          |
| NGINX_IGNITION_SECURITY_JWT_CLOCK_SKEW_SECONDS         | Amount of seconds that the token's dates can variate from the server dates                         | 60                                         | 60                                           |
| NGINX_IGNITION_CERTIFICATE_AUTO_RENEW_INTERVAL_MINUTES | Amount of minutes between the SSL certificates auto renew procedures executions                    | 60                                         | 60                                           |

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

nginx ignition can be used with its embedded database (powered by H2), but its strongly recommended that any long-term
installation is used alongside a PostgreSQL database. To do so, the following environment variables must be informed:

- `NGINX_IGNITION_DATABASE_URL` with the connection URL
- `NGINX_IGNITION_DATABASE_USERNAME` with the database username
- `NGINX_IGNITION_DATABASE_PASSWORD` with the database password

To exemplify, let's imagine that you have your SGDB running at the IP  `192.168.1.150`. The database name is 
`my_custom_db`, username is `my_username` and password `supersecretpassword`. In such scenario, the environment 
variables should have the following values:

```shell
NGINX_IGNITION_DATABASE_URL="jdbc:postgresql://192.168.1.150/my_custom_db"
NGINX_IGNITION_DATABASE_USERNAME="my_username"
NGINX_IGNITION_DATABASE_PASSWORD="supersecretpassword"
```

nginx ignition will create all the required tables, indexes and alike when it boots. In future updates, any changes
will be applied automatically also.
