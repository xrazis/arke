FROM golang:latest

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go build -o marvin .

CMD [ "/app/marvin" ]
