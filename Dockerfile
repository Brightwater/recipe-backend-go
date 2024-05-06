FROM golang:1.22

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o recipe_backend .

CMD ["./recipe_backend"]