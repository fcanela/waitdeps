
![WaitDeps](https://www.fcanela.com/waitdeps_logo.png)

# waitdeps

Waits for dependent services to become available.

This CLI tool ensures that dependent services such as databases, SMTP servers, or any other TCP-based services are available before proceeding with the application startup or tests. By waiting for these dependencies, it helps prevent the application from starting with invalid configurations, which could result in inconsistent states or failures.

It can be used:
- In production, to guarantee that all services are properly configured and ready
- In development and testing, to ensure that your application does not fail because the dependencies took a while to start.

## Installation

### Docker (Alpine)

```dockerfile
ENV WAITDEPS_VERSION 0.0.1
RUN apk update --no-cache \
  && apk add --no-cache wget openssl \
  && wget -nv -O - https://github.com/fcanela/waitdeps/releases/download/$WAITDEPS_VERSION/waitdeps-$WAITDEPS_VERSION-linux-amd64.tar.gz | tar xzf - -C /usr/local/bin \
  && apk del wget
```

### Docker (Ubuntu-based)

```dockerfile
ENV WAITDEPS_VERSION 0.0.1
RUN apt-get update \
    && apt-get install -y wget \
    && wget -nv -O - https://github.com/fcanela/waitdeps/releases/download/$WAITDEPS_VERSION/waitdeps-$WAITDEPS_VERSION-linux-amd64.tar.gz | tar xzf - -C /usr/local/bin \
    && apt-get autoremove -yqq --purge wget && rm -rf /var/lib/apt/lists/*
```

### Manual Installation

1. Find the latest release in the [releases section](https://github.com/fcanela/waitdeps/releases) of the repository.
2. Download the right version for your architecture and operating system.
3. Place the binary file in your system `PATH`.
## Documentation

Use `waitdeps wait` command. You can use the following flags:

- `--timeout`: how much time to wait to dependencies to become ready. Use a number followed by an unit (like 200ms, 30s, 1m, 1h).
- `--tcp-connect <service>`: checks the service readiness using a TCP connection attempt. The service can be in `host:port` format (`smtp.domain.com:25`) or URI (`postgres://user@host/db_name`). You can place this flag multiple times.

## Usage/Examples

Wait 10 seconds until PostgreSQL is available:
```
waitdeps wait --timeout 10s --tcp-connect db.myserver.com:5432
```

You can also wait other units of time:
```
waitdeps wait --timeout 500ms --tcp-connect db.myserver.com:5432
waitdeps wait --timeout 1m --tcp-connect db.myserver.com:5432
waitdeps wait --timeout 1h --tcp-connect db.myserver.com:5432
```

The timeout flag is optional, defaulting to 30s
```
waitdeps wait --tcp-connect db.myserver.com:5432
```

You can also use URIs. The port is automatically resolved from the URI schema:
```
waitdeps wait --tcp-connect https://www.google.com
waitdeps wait --tcp-connect postgresql://pguser:pgpass@db.myserver.com
```

To wait for multiple servers, you can place the `--tcp-connect` multiple times:
```
waitdeps wait --tcp-connect https://www.google.com --tcp-connect smtp.myserver.com:25
```

Alternatively, you can also use comma separated values:
```
waitdeps wait --tcp-connect https://www.google.com,smtp.myserver.com:25
```

## Authors

- [@fcanela](https://www.github.com/fcanela) - Francisco Canela


## License

[MIT](https://choosealicense.com/licenses/mit/)


## Related

- [dockerize](https://github.com/jwilder/dockerize) - The solution that I previously used and in which this tool is based. It is great, but it does not parse URIs for TCP checks and does too many things.
