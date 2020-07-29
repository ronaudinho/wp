# using alpine image since it is lighter than go
FROM golang:alpine AS build
RUN apk add --no-cache ca-certificates && update-ca-certificates
RUN mkdir /app/
WORKDIR /app/
COPY cmd/wp-sqlite/wp .
# create a single layer image
FROM scratch
# put binary in app folder
COPY --from=build /app/ /
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["./wp"]
LABEL Name=wp Version=0.1.0
EXPOSE 3195
