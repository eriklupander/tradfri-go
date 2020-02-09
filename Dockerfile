FROM golang:1.13

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 8080
EXPOSE 8081

CMD ["tradfri-go", "--server"]
