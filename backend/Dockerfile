FROM golang:1.23.2-alpine

WORKDIR /backend

COPY . .
RUN  go mod tidy && go build -o backend

EXPOSE 8080
CMD ["./backend"]