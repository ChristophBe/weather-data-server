FROM golang:latest

RUN apt-get update && apt-get install -y --no-install-recommends \
    libssl-dev libclamav-dev libmagic-dev libyara-dev liblzma-dev \
    && rm -rf /var/lib/apt/lists/*

RUN pwd
RUN wget https://github.com/neo4j-drivers/seabolt/releases/download/v1.7.4/seabolt-1.7.4-Linux-ubuntu-18.04.deb
RUN ls -l
RUN dpkg -i seabolt-1.7.4-Linux-ubuntu-18.04.deb

RUN mkdir /app

ADD . /app/
WORKDIR /app


RUN go mod download
RUN go build  -tags static_all -o main .
CMD ["/app/main"]