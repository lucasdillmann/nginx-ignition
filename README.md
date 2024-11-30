<p align="center">
    <img src="docs/readme-screenshots-v2.png" alt="" width="600" />
</p>
<h1 align="center">
    nginx ignition
</h1>

The nginx ignition is a user interface for the nginx web server, aimed at developers and enthusiasts that don't
want to manage configuration files manually for their use-cases. 

Although it isn't the goal to be feature-complete (if your use-case is quite advanced or complex, you probably will not 
use a UI anyway), the project does aim to provide an intuitive and powerful way to configure and run nginx.

Some of the available features include:
- Multiple nginx virtual hosts, each one with its customized set of domain, routes and bindings (port listeners)
- Each host route can act as a proxy, redirection or reply with a static response
- SSL certificates (Let's Encrypt, self-signed or bring your custom one) with automatic renew (when applicable)
- Server and virtual hosts access and error logs
- Multiple users with role-based access control (RBAC)

## Getting started

To run nginx ignition, run the following in your terminal. If you don't have it already, you will need to install Docker
first (more details on hot to do it, follow [this link](https://www.docker.com/get-started/)).

```shell
docker run -p8090:8090 dillmann/nginx-ignition
```

After a few seconds, you can open your favorite browser at http://localhost:8090 and start using it. There's no 
default username or password, the nginx ignition will guide you to create your user.

Please note that in its default configuration the app will start using an embedded H2 database. While this is fine for
testing and some experiments, is not recommended for a long-term scenario. For that, please refer to the 
configuration section below to use PostgreSQL instead.

## Configuration

Check [this documentation file](docs/configuration-properties.md) for more details about the available 
configuration properties and some common use-case examples.

## Contributing and feedback

Feel free to open an issue or pull request here at GitHub or to send me a message through
[my LinkedIn profile](https://linkedin.com/in/lucasdillmann) to share some feedback. Every constructive one is welcome.
