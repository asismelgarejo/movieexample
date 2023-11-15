FROM alpine:latest
COPY exe/main .
COPY configs/. ./configs/.
EXPOSE 8086
ENV mode="production"
CMD ["/main"]