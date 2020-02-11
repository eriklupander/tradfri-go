FROM golang:1.13

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 80
EXPOSE 81

CMD ["tradfri-go", "--server"]
