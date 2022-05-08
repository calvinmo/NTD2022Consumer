FROM alpine

COPY install-alpine.sh /opt/install.sh

RUN /opt/install.sh

COPY dist /dist

WORKDIR /dist

USER nobody

ENTRYPOINT ["./service"]
