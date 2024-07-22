FROM golang:1.22.5

WORKDIR /app

COPY . .

EXPOSE 8080

RUN go build -o api cmd/main.go

CMD [ "./api" ]
