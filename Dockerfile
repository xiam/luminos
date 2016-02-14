FROM debian:jessie

ENV VERSION 0.9.3

RUN apt-get update

RUN apt-get install -y curl

ENV LUMINOS_URL https://github.com/xiam/luminos/releases/download/v$VERSION/luminos_linux_amd64.gz

RUN curl --silent -L ${LUMINOS_URL} | gzip -d > /bin/luminos

RUN chmod +x /bin/luminos

EXPOSE 9000

ENTRYPOINT [ "/bin/luminos" ]
