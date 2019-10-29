# nutstat

A simple Network UPS Tools (NUT) to InfluxDB exporter.

## why

To build myself a neat Grafana dashboard like this...

![](https://i.imgur.com/dsTc0lg.png)

## usage
I use this with docker, but using docker-compose. If you want to check your configuration, do something like this:

```bash
docker run --rm -v $(pwd)/nutstat.yml:/config/nutstat.yml --network stats nutstat:local --config /config/nutstat.yml check
```
