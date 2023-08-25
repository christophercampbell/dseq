FROM alpine:3.7

RUN apk update && \
    apk upgrade && \
    apk --no-cache add curl jq file

VOLUME /dseq
WORKDIR /dseq
EXPOSE 26656 26657
ENTRYPOINT ["/usr/bin/wrapper.sh"]
STOPSIGNAL SIGTERM

COPY wrapper.sh /usr/bin/wrapper.sh
COPY config-template.toml /etc/dseq/config-template.toml
