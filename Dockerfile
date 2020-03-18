FROM golang:latest

RUN mkdir /varq
ADD . /varq/
WORKDIR /varq
RUN go get -d -v
RUN go build -o main .

CMD ["/varq/main"]