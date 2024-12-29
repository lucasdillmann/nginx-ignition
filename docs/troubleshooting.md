# Troubleshooting

nginx ignition was designed to have helpful and useful logs. If something is not working as expected and the feedback
on the UI isn't clear enough, the first step is to check the application logs and start from there.

Such logs are available in the Docker's container logs (stdout of the container). Please check the 
[Docker's documentation](https://docs.docker.com/engine/logging/) for more details on how to view them.

Beyond that, this document also has some common troubleshooting scenarios below. If both logs and they don't fix your 
problem, please raise an issue on the project's GitHub repository.

## Password reset

If you lost your password, you can reset any user's password using by setting the environment variable below in the
nginx ignition's Docker container.

```shell
NGINX_IGNITION_PASSWORD_RESET_USERNAME=<your username>

# Example
NGINX_IGNITION_PASSWORD_RESET_USERNAME=admin
```

When the application starts, it will generate a new and random password for the specified user, printing the new 
password in the logs like the example below. 

```text
[2024-12-29 14:16:43.128] INFO: Password reset completed successfully for the user admin. New password: 7da75d54
```

Please note that application will not finish the startup/boot process with the environment variable above set, it will 
only reset the password and then shut down. If you run the application multiple times with the environment variable set,
the password will be changed every time the application executes.
