nagios-check-graphite
==============

A basic Graphite check for Nagios.

```
Usage:
  check-graphite [OPTIONS]

Application Options:
  -H, --host=HOST                          Graphite host url (required)
  -m, --metric=METRIC                      Graphite metric name (required)
  -z, --zero                               Convert 'None' values to 0.
  -s, --scale=SCALE                        Set the desired numeric scale for the values. (default: 2)
  -d, --duration=SECONDS                   Number of minutes of data to aggregate. (default: 10)
  -f, --function=(min|max|avg|sum|last)    The aggregation function to apply. (default: last)
  -w, --warning=WARNING                    Warning threshold of aggregated value.
  -c, --critical=CRITICAL                  Critical threshold of aggregated value.
  -i, --invert                             Invert thresholds to alert below metric value.

Help Options:
  -h, --help      Show this help message
```

Example Nagios command:
```
define command {
    command_name    check-graphite
    command_line    $USER1$/check-graphite -H http://graphite.mydomain -m "$ARG1$" -w $ARG2$ -c $ARG3$
}
```

Example Nagios service:
```
define service {
    service_description    load average
    host_name              graphite
    check_command          check-graphite!collectd.graphite.load.load.longterm!1.0!2.5
}
```
