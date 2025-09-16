# Restoring your nginx ignition database

This guide explains how to restore the database you downloaded from the Export page. The process differs depending on 
the database you use:
- SQLite: the download is a single .db file (nginx-ignition.db)
- PostgreSQL: the download is a plain SQL file (nginx-ignition.sql)

Important safety notes
- Stop the app (container) or ensure it's not writing to the database during the restore to avoid corruption.
- Always keep an extra copy of your backup before overwriting anything.

## SQLite restore
Use this if your nginx ignition is running with the embedded SQLite database (the default for fresh and default installs).

Official Docker image expects the file to be placed at `/opt/nginx-ignition/data` using the file name `nginx-ignition.db`,
where the full path is `/opt/nginx-ignition/data/nginx-ignition.db`. This path can be changed using the env var 
`NGINX_IGNITION_DATABASE_DATA_PATH` (see [this documentation file](configuration-properties.md) for more details).

Option A: Copy back the file into the running Docker container
1) Identify the container name (example: `nginx-ignition`)
   - Example: `docker ps`
2) Stop the app to avoid writes during replacement
   - Example: `docker stop <container-name>`
3) Copy the downloaded database file into the expected path
   - Example: `docker cp /path/to/nginx-ignition.db <container-name>:/opt/nginx-ignition/data/nginx-ignition.db`
4) Start the container again
   - Example: `docker start <container-name>`

Option B: Using a Docker bind mount for the data folder
1) Stop the container
   - Example: `docker stop <container-name>`
2) Replace file on host where the bind mount lives 
   - Example: `cp /path/to/nginx-ignition.db /path/to/bind/mount/nginx-ignition.db`
3) Start the container again 
   - Example: `docker start <contaienr-name>`

## PostgreSQL restore
Use this if your nginx ignition is configured to use a PostgreSQL server.

The database restore process can be done using the official PostgreSQL client by running the `psql` command. You may
need to install it first, see https://www.postgresql.org/download/ for instructions.

Please note that this guide assumes that you have already created a PostgreSQL database, has an user with privileges to
create/alter/drop objects present in the dump, and that such database is currently running and empty. Please check the
official PostgreSQL documentation at https://www.postgresql.org/docs/ for more details and instructions.

Steps to restore the database:
1) Identify the container name (example: `nginx-ignition`)
   - Example: `docker ps`
2) Stop the app to avoid writes during replacement
   - Example: `docker stop <container-name>`
3) Execute the psql to restore the database
   - Example: `psql --host=<your-db-host> --port=5432 --username=<your-username> --dbname=<your-db-name> --file=/path/to/nginx-ignition.sql`
4) Fill the database password if and when prompted
5) Start the container again
   - Example: `docker start nginx-ignition`
