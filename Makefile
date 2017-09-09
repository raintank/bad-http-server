.PHONY: build
build:
	docker run --rm -v `pwd`:/go/src/bad-http-server/ golang bash -c "GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o src/bad-http-server/bad-http-server bad-http-server"
	docker build -t bad-http-server .
