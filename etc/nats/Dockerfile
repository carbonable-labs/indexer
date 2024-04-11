FROM nats:2.10-alpine

RUN apk add --no-cache bind-tools

WORKDIR /srv

COPY ./etc/nats/nats-server.conf /etc/nats/nats-server.conf
COPY ./etc/nats/entrypoint.sh /srv/entrypoint.sh
RUN chmod +x /srv/entrypoint.sh

ENTRYPOINT ["/srv/entrypoint.sh"]
CMD ["nats-server", "--config", "/etc/nats/nats-server.conf", "--jetstream"]
