FROM debian:bookworm-slim

WORKDIR /app
COPY ./build/frontend .
COPY ./templates ./templates
COPY ./static ./static

CMD ["./frontend"]

