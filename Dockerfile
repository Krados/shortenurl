FROM golang:1.19

COPY . /app
WORKDIR /app

RUN make build

EXPOSE 8080

CMD [ "bin/cmd", "-conf", "./configs" ]