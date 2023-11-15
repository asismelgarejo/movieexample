FROM alpine:latest
COPY exe/main .
COPY configs/. ./configs/
EXPOSE 8083
ENV mode="production"
CMD ["/main"]