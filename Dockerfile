FROM alpine:3.6
COPY bad-http-server /bad-http-server
ENTRYPOINT ["/bad-http-server"]
CMD ["0.0.0.0:8888"]
