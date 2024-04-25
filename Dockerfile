FROM golang:1.18.3-alpine3.16
RUN mkdir /app
ADD . /app
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
# ENV EINHORN_FDS="3"
ENV ISLOCAL="1"
# RUN go get github.com/zimbatm/socketmaster 
EXPOSE 10004
RUN go build -o joona /app/cmd/

CMD ["/app/joona"]
