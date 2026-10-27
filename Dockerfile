
FROM golang:1.21.1

WORKDIR /app

RUN go install github.com/revel/cmd/revel@latest

COPY . .

RUN rm -rf /app/target && revel build -a /app -t /build/target -m prod

WORKDIR /build/target

RUN chmod +x run.sh && mkdir -p log

EXPOSE 9000

ENTRYPOINT ["./run.sh"]
