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

This `docker-compose.yml` should be good to get you started. You need `nutstat.yml` to be configured with your values. After saving both, start importing stats with `docker-compose up --build -d`.

```yml
version: "3"
services:
  nutstat:
    build:
      context: https://github.com/leonjza/nutstat.git
      dockerfile: Dockerfile
    image: nutstat:local
    command: --config /config/nutstat.yml update
    container_name: nutstat
    volumes:
      - ./config/nutstat.yml:/config/nutstat.yml
    restart: unless-stopped
    networks:
      - stats

networks:
  stats:
    external: true
```
