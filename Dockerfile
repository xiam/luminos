FROM debian:jessie

RUN apt-get update

RUN apt-get install -y curl

ENV LUMINOS_URL https://github.com/xiam/luminos/releases/download/v0.9.2/luminos_linux_amd64.gz

RUN curl --silent -L ${LUMINOS_URL} | gzip -d > /bin/luminos

RUN chmod +x /bin/luminos

EXPOSE 9000

ENTRYPOINT [ "/bin/luminos" ]
