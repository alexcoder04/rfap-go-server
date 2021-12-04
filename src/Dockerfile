FROM golang:1.17

WORKDIR /go/src/app
RUN mkdir out

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["make", "windows"]
