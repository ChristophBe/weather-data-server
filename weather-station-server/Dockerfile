FROM golang:latest
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go get -u github.com/gorilla/mux
RUN go get -u  github.com/johnnadratowski/golang-neo4j-bolt-driver
RUN go get -u  golang.org/x/crypto/bcrypt
RUN go build -o main .
CMD ["/app/main"]