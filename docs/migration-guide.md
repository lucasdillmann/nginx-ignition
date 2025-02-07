# Migration guide from 1.x to 2.0.0
Version 2.0.0 of the nginx ignition introduces some minor breaking changes. This guide details such changes and how to 
migrate your existing configuration to the new version.

## Important notes
nginx ignition supports an embedded database for quick tests and experimentation of the app, but it isn't recommeded
for production on long-term use cases. In version 1.x, such embedded database was H2 and in version 2.0.0 it was 
changed to SQLite.

The migration from 1.x to 2.0.0 is not possible if you're using the embedded database, but will work just fine if you're
using the PostgreSQL database, even if you're not comming from the latest 1.x version.

## Configuration changes
Most of the changes are in the name of the environment variables that the application uses, specifically the ones with 
the credentials to connect to the database.

1. Variable `NGINX_IGNITION_DATABASE_URL` is now split into four variables:
   - `NGINX_IGNITION_DATABASE_HOST` with the hostname or IP of the database
   - `NGINX_IGNITION_DATABASE_DRIVER` with the type of the database, being either `sqlite` or `postgres`
   - `NGINX_IGNITION_DATABASE_PORT` with the port where the database is listening for connections
   - `NGINX_IGNITION_DATABASE_NAME` with the name of the database
   - `NGINX_IGNITION_DATABASE_SSL_MODE` with the SSL mode to connect to the database, being either `disable` or `require`
     (value `require` is the default one by omission or if the value is not recognized)

2. Variables `NGINX_IGNITION_DATABASE_CONNECTION_POOL_MINIMUM_SIZE` and 
   `NGINX_IGNITION_DATABASE_CONNECTION_POOL_MAXIMUM_SIZE` that previously defined the connection pool to the database
   no longer take effect and can be removed from your configuration.

The remaining variables are the same as before and work in the same way. If you wish to get the full list of environment
variables that nginx ignition supports, please refer to the [configuration properties](configuration-properties.md) 
documentation file.

## Example of migration
Let's imagine that you previously had the following environment variables defined in your configuration:

```shell
NGINX_IGNITION_DATABASE_URL="jdbc:postgresql://192.168.1.150:5432/my_custom_db"
NGINX_IGNITION_DATABASE_USERNAME="my_username"
NGINX_IGNITION_DATABASE_PASSWORD="supersecretpassword"
```

Most of the information in the new variables are the ones in the previous `NGINX_IGNITION_DATABASE_URL` variable. By
migration to version 2.0.0, the configuration variables would become the following ones:

```shell
NGINX_IGNITION_DATABASE_HOST=192.168.1.150
NGINX_IGNITION_DATABASE_DRIVER=postgres
NGINX_IGNITION_DATABASE_PORT=5432
NGINX_IGNITION_DATABASE_NAME=my_custom_db
NGINX_IGNITION_DATABASE_DRIVER=disable
NGINX_IGNITION_DATABASE_USERNAME=my_username
NGINX_IGNITION_DATABASE_PASSWORD=supersecretpassword
```
