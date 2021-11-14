FROM golang:alpine as builder

	WORKDIR /github.com/lukasdietrich/webzipd
	COPY . .

	RUN go build ./cmd/webzipd

FROM alpine

	WORKDIR /app
	COPY --from=builder /github.com/lukasdietrich/webzipd/webzipd .

	USER nobody
	VOLUME /data
	EXPOSE 8080

	CMD ["/app/webzipd", "-address", ":8080", "-mode", "hostname", "-folder", "/data"]

