# Movies persons service #
[![Go Report Card](https://goreportcard.com/badge/github.com/Falokut/movies_persons_service)](https://goreportcard.com/report/github.com/Falokut/movies_persons_service)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/Falokut/movies_persons_service)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Falokut/movies_persons_service)
[![Go](https://github.com/Falokut/movies_persons_service/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/Falokut/movies_persons_service/actions/workflows/go.yml) ![](https://changkun.de/urlstat?mode=github&repo=Falokut/movies_persons_service)
[![License](https://img.shields.io/badge/license-MIT-green)](./LICENSE)
---

# Content
+ [Configuration](#configuration)
    + [Params info](#configuration-params-info)
        + [Database config](#database-config)
        + [Jaeger config](#jaeger-config)
        + [Prometheus config](#prometheus-config)
        + [time.Duration](#timeduration-yml-supported-values)
+ [Metrics](#metrics)
+ [Docs](#docs)
+ [Author](#author)
+ [License](#license)
---------

# Configuration

1. [Configure movies_persons_db](movies_persons_db/README.md#Configuration)
2. Create .env on project root dir  
Example env:
```env
REDIS_PASSWORD=redispass
REDIS_AOF_ENABLED=no
DB_PASSWORD=Password
```

## Configuration params info
if supported values is empty, then any type values are supported

| yml name | yml section | env name | param type| description | supported values |
|-|-|-|-|-|-|
| log_level   || LOG_LEVEL  |   string   |      logging level        | panic, fatal, error, warning, warn, info, debug, trace|
| healthcheck_port   |      | HEALTHCHECK_PORT  |   string   |     port for healthcheck       | any valid port that is not occupied by other services. The string should not contain delimiters, only the port number|
| host   |  listen    | HOST  |   string   |  ip address or host to listen   |  |
| port   |  listen    | PORT  |   string   |  port to listen   | The string should not contain delimiters, only the port number |
| server_mode   |  listen    | SERVER_MODE  |   string   | Server listen mode, Rest API, gRPC or both | GRPC, REST, BOTH|
| service_name   |  prometheus    | PROMETHEUS_SERVICE_NAME | string |  service name, thats will show in prometheus  ||
| server_config   |  prometheus    |   | nested yml configuration  [metrics server config](#prometheus-config) | |
|db_config|||nested yml configuration  [database config](#database-config) || configuration for database connection | |
|jaeger|||nested yml configuration  [jaeger config](#jaeger-config)|configuration for jaeger connection | |
|network|movies_persons_cache|MOVIES_PERSONS_CACHE_NETWORK|string| network| tcp or udp|
|addr|movies_persons_cache|MOVIES_PERSONS_CACHE_ADDR|string|ip address(or host) with port of redis| all valid addresses formatted like host:port or ip-address:port |
|password|movies_persons_cache|MOVIES_PERSONS_CACHE_PASSWORD|string|password for connection to the redis||
|db|movies_persons_cache|MOVIES_PERSONS_CACHE_DB|int|the number of the database in the redis||
|movies_persons_ttl|movies_persons_cache|| time.Duration with positive duration | the time delay between person deleted from cache|[supported values](#timeduration-yml-supported-values)|
| base_photo_url   |      | BASE_PHOTO_URL  |   string   | base url for pictures | url like http://host or https://host|
| photo_category   |      | PHOTO_CATEGORY  |   string   |category name for photos in images storage server|only latin characters, without special characters|

### Database config
|yml name| env name|param type| description | supported values |
|-|-|-|-|-|
|host|DB_HOST|string|host or ip address of database| |
|port|DB_PORT|string|port of database| any valid port that is not occupied by other services. The string should not contain delimiters, only the port number|
|username|DB_USERNAME|string|username(role) in database||
|password|DB_PASSWORD|string|password for role in database||
|db_name|DB_NAME|string|database name (database instance)||
|ssl_mode|DB_SSL_MODE|string|enable or disable ssl mode for database connection|disabled or enabled|

### Jaeger config

|yml name| env name|param type| description | supported values |
|-|-|-|-|-|
|address|JAEGER_ADDRESS|string|hip address(or host) with port of jaeger service| all valid addresses formatted like host:port or ip-address:port |
|service_name|JAEGER_SERVICE_NAME|string|service name, thats will show in jaeger in traces||
|log_spans|JAEGER_LOG_SPANS|bool|whether to enable log scans in jaeger for this service or not||

### Prometheus config
|yml name| env name|param type| description | supported values |
|-|-|-|-|-|
|host|METRIC_HOST|string|ip address or host to listen for prometheus service||
| port|METRIC_PORT|string|port to listen for  of prometheus service| any valid port that is not occupied by other services. The string should not contain delimiters, only the port number|

### time.Duration yml supported values
A Duration value can be expressed in various formats, such as in seconds, minutes, hours, or even in nanoseconds. Here are some examples of valid Duration values:
- 5s represents a duration of 5 seconds.
- 1m30s represents a duration of 1 minute and 30 seconds.
- 2h represents a duration of 2 hours.
- 500ms represents a duration of 500 milliseconds.
- 100µs represents a duration of 100 microseconds.
- 10ns represents a duration of 10 nanoseconds.

# Metrics
The service uses Prometheus and Jaeger and supports distribution tracing

# Docs
[Swagger docs](swagger/docs/movies_persons_service_v1.swagger.json)

# Author

- [@Falokut](https://github.com/Falokut) - Primary author of the project

# License

This project is licensed under the terms of the [MIT License](https://opensource.org/licenses/MIT).

---