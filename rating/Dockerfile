FROM alpine:latest
COPY exe/main .
COPY configs/. ./configs/.
EXPOSE 8082
ENV mode="production"
CMD ["/main"]