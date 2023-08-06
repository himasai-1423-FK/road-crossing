FROM golang:1.21rc2-alpine3.17

WORKDIR /app/

COPY . /app/

RUN go get .

RUN go build -o myapp

EXPOSE 3000

CMD ["./myapp"]