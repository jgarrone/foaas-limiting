# FOAAS Limiting

## Overview

This repository contains a server that fetches messages from [FOOAS](https://www.foaas.com/) and returns them. The
server implements a fixed-window rate-limiter, which denies service to users making more requests than allowed.

## Getting Started

The server is implemented in the [Go language](https://go.dev/). It was tested using Go 1.16. To install Go, head
to https://go.dev/doc/install and follow the instructions.

To run the server, go to the root of the repository and execute:

```
make all
make run
```

That's it. The first command will test and build the code. The second command starts the server.

If you have Docker installed and prefer to run it in a containerised environment, run:

```
make docker-build
make docker-run
```

After setting up the server, it can be tested with `curl` or a similar tool. The only available path on this server
is `/message`. The request must contain a header named `userId` and map to the name of the individual using this API.
The limiter is a bit naive and uses this header to do its job.
**Warning**: Using a name that isn't yours is a felony and could get you into trouble.

Example:

```
curl localhost:8080/message -H 'userId: <my-name>'
```

## Customising the server

By default, the server listens at `localhost:8080`. To specify the address where the server should listen:

```
./foaas-limiting serve --listen_address=<host>:<port>
```

The limiter can be tweaked with specific arguments. One for the time window and another for the request count limit. The
default parameters allow 5 requests ever 10 seconds. To change to, for example, 10 requests per second:

```
./foaas-limiting serve --rate_limit_count=10 --rate_limit_window_ms=1000
```

The server tracks the latency when making requests to FOAAS API. This can be observed by enabling the debug log level:

```
./foaas-limiting serve --log_level=debug
```

## Future work

- Implement autn/authz
- Leverage the `/operations` endpoint of FOAAS to pick random messages at run time  
- Support customising limiter configuration for different endpoints
- Use a remote storage like Redis to support limiting requests in a distributed environment
