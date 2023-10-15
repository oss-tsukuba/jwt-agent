# jwt-agent

**jwt-agent** obtains a JSON Web Token (JWT) from a [JWT
server](https://github.com/oss-tsukuba/jwt-server.git), and keep
refreshing it not to expire.  It is running in the background unless
the -f option is specified.

## Usage
```
Usage: jwt-agent [-s URL] [-l user] [-f] [-t timeout]
       jwt-agent --status
       jwt-agent --stop [-t timeout]
       jwt-agent --version
```

When the -s option is not specified, `JWT_SERVER_URL` environment
variable is used.  When the -l option is not specified, `LOGNAME`
environment variable is used.

The jwt-agent asks a passphrase at the start up to obtain a JWT, which
is provided by a JWT server.  The jwt-agent also accepts the
passphrase by the standard input.  The jwt-agent does not stop unless
it is explicitly stopped, or some error happens.

By default, the jwt-agent stores a JWT at
`/tmp/jwt_user_u$UID/token.jwt`, which can be changed by
`JWT_USER_PATH` environment variable.

jwt-agent --status checks the running status of the jwt-agent.

jwt-agent --stop stops the jwt-agent execution.

jwt-agent --version displays the version number.

## How to build and install

jwt-agent-core is written in Go.  Go 1.18 or later is required.

    % make
    % sudo make PREFIX=/usr/local install

The default `PREFIX` is /usr.
