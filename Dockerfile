FROM golang:1.16

WORKDIR /app

COPY . .

WORKDIR /app/cmd
RUN go build -o hexsatisfaction_purchase .

CMD ["./hexsatisfaction_purchase"]


