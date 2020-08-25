FROM golang:alpine
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go mod download
RUN go build  github.com/ChristophBe/weather-data-server
CMD ["/app/weather-data-server"]