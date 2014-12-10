FROM golang:latest

RUN go get github.com/braintree/manners
RUN go get github.com/googollee/go-middleware
RUN go get github.com/julienschmidt/httprouter
RUN go get labix.org/v2/mgo
