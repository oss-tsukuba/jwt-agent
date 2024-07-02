% jwt-agent(1)
%
% September 2, 2023

# NAME

jwt-agent - Obtain and keep refreshing a JSON Web Token

# SYNOPSIS

**jwt-agent** [-s _URL_] ... [-l _user_] [-f] [-t _timeout_]
**jwt-agent** --status
**jwt-agent** --stop [-t _timeout_]
**jwt-agent** --version

# DESCRIPTION

**jwt-agent** obtains a JSON Web Token (JWT) from a [JWT
server](https://github.com/oss-tsukuba/jwt-server.git), and keep
refreshing it not to expire.  It is running in the background unless
the -f option is specified.  When the -s option is not specified,
JWT_SERVER_URL environment variable is used.  When the -l option is
not specified, LOGNAME environment variable is used.

The jwt-agent asks a passphrase at the start up to obtain a JWT, which
is provided by a JWT server.  The jwt-agent also accepts the
passphrase by the standard input.  The jwt-agent does not stop unless
it is explicitly shut down, or some error happens.

By default, the jwt-agent stores a JWT at
/tmp/jwt_user_u$UID/token.jwt, which can be changed by JWT_USER_PATH
environment variable.

When multiple -s options are specified, the jwt-agent attempts to
retrieve the JWT in order if the server is not working.
If the passphrase is incorrect, the program terminates at that point.
Also, the jwt-agent exits if all URLs are not working.

# OPTIONS

-s _URL_
: specifies the URL of a JWT server

-l _user_
: specifies a user name

-f
: executes in the foreground not in the background

-t _timeout_
: specifies the timeout in seconds.  Default is 60 seconds.

--status
: checks the running status of the jwt-agent

--stop
: stops the jwt-agent execution

--version
: displays the version number

# ENVIRONMENT

JWT_USER_PATH
: path to the JSON Web Token.  Default is /tmp/jwt_user_u$UID/token.jwt

JWT_SERVER_URL
: URL of a JWT server
