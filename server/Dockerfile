FROM golang:latest

RUN mkdir /app 

WORKDIR /app 

COPY . /app/server
# COPY .env .
RUN cd /app/server && go mod download all 

RUN cd /app/server/cmd && go build -o /go-app 

EXPOSE 8080

CMD ["/go-app"]