# Tracer

Tracer is a tool to query and aggregate metrics and traces from an Aporeto control plane.

## Usage

You will need client certificate to query the monitoring API.

```console
export TRACER_MONITORING_CERT=/path/to/cert.pem
export TRACER_MONITORING_CERT_KEY=/path/to/key.pem
export TRACER_MONITORING_CERT_KEY_PASS="p0ul3t"
export TRACER_MONITORING_URL=https://monitoring.poulet.com
```

Then you can play, see `--help`

```console
Usage:
      --code string                       Filters: The code to filter ex:200-300,400-422,500
      --direction string                  Logs: Direction of the logs [allowed: forward,backward] (default "forward")
      --errors-only                       Traces: Look only for trace in error
      --follow                            Logs: Follow logs stream in almost real time
      --from string                       From date
      --help                              Show full help with examples
      --limit int                         Traces: The number of traces to display (default 1)
      --lines int                         Logs: Number of lines to print (default 10)
      --log                               Logs: Enable log mode to get logs from services
      --log-filter string                 Logs; Optional log filter to append to log query
      --log-format string                 Log format (default "console")
      --log-level string                  Log level (default "info")
      --monitoring-ca-path string         Path to the monitoring CA certificate
      --monitoring-cert string            Path to the monitoring cert
      --monitoring-cert-key string        Path to the monitoring cert key
      --monitoring-cert-key-pass string   Password for the monitoring cert key
      --monitoring-url string             The monitoring url to use
      --namespace string                  Traces: Lookg for traces matching that namespace
      --no-labels                         Logs: Do not display labels with logs
      --open string                       Traces: Open a given trace to your browser.
      --profile-file string               Profile file: the profile file pathto use. (default "~/.tracer/default.yaml")
      --service strings                   Filters: The service to filter (repeatable)
      --since duration                    Since duration (will compute From and To with currrent date) (default 1h0m0s)
      --slower-than duration              Traces: Look for traces slower than the provided duration
      --stack string                      Stack: The stack name to use if any. (default "default")
      --to string                         To date
      --url strings                       Filters: The url to filter (repeatable)
  -v, --version                           Display the version


> Display all queries with traces from the last 1h

  ./tracer --since 1h

> Display all queries for a service from the last 1h

  ./tracer --since 1h --service squall

> Display all queries for a service in a given namespace from the last 1h

  ./tracer --since 1h --service squall --namespace /foo/bar

> Display all queries for a service in a given namespace that took more than 2s from the last 1h

  ./tracer --since 1h --service squall --namespace /foo/bar --slower-than 2s

> Display all requests that returns with an error for the past hour

  ./tracer --since 1h --errors-only

> Display all requests that return with a code 200 or 400-422 in the past hour

  ./tracer --since 1h --code 200,400-422

> Display all requests made to /flowreports

  ./tracer --since 1h --url /flowreports

> Display all 400-403 requests on service squall, cid and /issue between two dates

  ./tracer --code 400-403 --service squal --service cid --url /issue --from 2020-10-21T17:56:17Z --to 2020-10-22T17:56:17Z

> Display logs for 2 services between two dates

  ./tracer --log --service squal --service cid --from 2020-10-21T17:56:17Z --to 2020-10-22T17:56:17Z

> Display logs with a custom filter for a given service

  ./tracer --log --service squall --log-filter '|~"ERROR"'

> Displau logs with a custom filter

  ./tracer --log --log-filter '{type="aporeto",app!~"squall|wutai.*"}|~"ERROR|WARNING"'

Some queries are not providing traces (like reports because this is too much for jaeger to handle).
In general errors are logged in the service in debug mode. Use the switch-debug <service name>  command to enable it.
And look at the logs either through Grafana->Explore->Loki or with the k get log <pod_name> command.
```

Example of output:

```console
./tracer --since 1m

  count |    service    |       identity       |   operation   |            url             | code |         traces (limit=2)
--------+---------------+----------------------+---------------+----------------------------+------+------------------------------------
      2 | squall        | processingunit       | info          | /processingunits           |  204 | 6a113d0efa9b259b,491cfcaef5a8e343
      2 | midgard       | issue                | create        | /issue                     |  500 |
      2 | squall        | enforcer             | info          | /enforcers                 |  204 | 1ed3127e33c0cd05,746b8d9d280bac7a
      4 | squall        | datapathcertificate  | create        | /datapathcertificates      |  403 | 2d729ad76e3a7c48,32fa0a50a091607c
      4 | meteor        | graphedge            | retrieve-many | /graphedges                |  200 | 587a4de05be76a7e,59d453fa6da22c95
      4 | zack          | counterreport        | create        | /counterreports            |  403 |
      4 | meteor        | graphnode            | retrieve-many | /graphnodes                |  200 | 587a4de05be76a7e,59d453fa6da22c95
      4 | gaga          | poke                 | retrieve-many | /enforcers/:id/poke        |  403 | 0b390744000f683a,070368d07ad4bc54
      4 | jenova        | dependencymap        | retrieve-many | /dependencymaps            |  200 | 587a4de05be76a7e,59d453fa6da22c95
      4 | zack          | enforcerreport       | create        | /enforcerreports           |  403 |
      6 | squall        | externalnetwork      | retrieve      | /externalnetworks/:id      |  200 | 587a4de05be76a7e,59d453fa6da22c95
      8 | sephiroth-api | alarm                | create        | /alarms                    |  403 | 44bfe6894099ae75,31bdc4ad168758fc
      8 | squall        | processingunit       | retrieve-many | /processingunits           |  403 | 411c321e7f3621c6,2dbef2645d32ac14
      8 | leon          | eventlog             | create        | /eventlogs                 |  403 | 7b2c065d74f82dfa,399ccabbf13a4e05
     11 | cid           | authz                | create        | /authz                     |  200 | 72f2c675a9e06544,40f20ffdcdba37c8
     16 | midgard       | issue                | create        | /issue                     |  200 | 211c4e34e7b643ff,2db1a90e21745544
     22 | barret        | x509certificatecheck | retrieve      | /x509certificatechecks/:id |  204 | 587a4de05be76a7e,6a113d0efa9b259b
     34 | jenova        | statsquery           | create        | /statsqueries              |  200 | 713dfa716fb9ce51,14d0c5ec6a428964
     72 | zack          | enforcerreport       | create        | /enforcerreports           |  204 |
     78 | gaga          | poke                 | retrieve-many | /enforcers/:id/poke        |  204 | 714b6134c9c6bcfc,6f9be5fa82241917
     84 | zack          | flowreport           | create        | /flowreports               |  204 |
    184 | zack          | counterreport        | create        | /counterreports            |  204 |
    278 | zack          | dnslookupreport      | create        | /dnslookupreports          |  204 |


> 23 results found. You can read the traces from https://monitoring.poulet.com/explore and select the jaeger datasource.
```

> Note: some query are not generating traces, in general the reports because there is too much of them.

## Profiles

You can create profiles see the `--profile-file` flag with default value `~/.tracer/default.yaml` as:

```yaml
datasources:
  - name: default
    monitoringCAPath: /path/to/ca-chain-public.pem
    monitoringCertPath: /path/to/cert.pem
    monitoringCertKeyPath: /path/to/key.pem
    monitoringURL: https://monitor.default.poulet.com
  - name: foo
    monitoringCAPath: /path/to/ca-chain-public.pem
    monitoringCertPath: /path/to/cert.pem
    monitoringCertKeyPath: /path/to/key.pem
    monitoringURL: https://monitor.foo.poulet.com
```

Then use select a profile with `--stack <name>` flag.
