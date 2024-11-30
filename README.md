<p align="center">
    <img src="docs/readme-screenshots-v2.png" alt="" width="600" />
</p>
<h1 align="center">
    nginx ignition
</h1>

The nginx ignition is a user interface for the nginx web server, aimed for developers and enthusiasts that don't
want to manage configuration files manually for their use-cases. 

Although isn't the goal to be feature-complete (if your use case is quite advanced or complex, you probably will not 
use a UI anyway), the project does aim to provide a intuitive and powerful way to configure and run nginx.

Some of the available features include:
- Multiple nginx virtual hosts, each one with its customized set of domain, routes and bindings (port listeners)
- SSL certificates (Let's Encrypt, self-signed or bring your custom one) with automatic renew (when applicable)
- Server and virtual hosts access and error logs
- Multiple users with role based access control (RBAC)

## Getting started

To run nginx ignition, run the following in your terminal. If you don't have it already, you will need to install Docker
first (more details on hot to do it, follow [this link](https://www.docker.com/get-started/)).

```shell
docker run -p8090:8090 dillmann/nginx-ignition
```

After a few seconds, you can open your favorite browser at http://localhost:8090 and start using it. There's no 
default username or password, the nginx ignition will guide you to create your user.

Please note that in its default configuration, the app will start using an embedded database. While this is fine for
testing and some experiments, is not recommended for a production-like scenario. For that, please refer to the 
configuration section below to use PostgreSQL instead.

## Configuration

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
