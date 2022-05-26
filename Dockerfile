FROM golang:1.18-alpine as BUILD
RUN mkdir /app
WORKDIR /app
COPY src /app
RUN go build -tags musl 

FROM alpine:3.15.4 
RUN mkdir /app
WORKDIR /app
COPY --from=BUILD /app/uptime-robot-exporter /app
CMD [ "/app/uptime-robot-exporter" ]
