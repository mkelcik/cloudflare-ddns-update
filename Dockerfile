FROM golang:1.20 as build

# Copy project sources
COPY . /opt/project/
WORKDIR /opt/project

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates=20210119

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /cloudflare-ddns-updater

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /cloudflare-ddns-updater /cloudflare-ddns-updater
ENTRYPOINT ["/cloudflare-ddns-updater"]
CMD ["cloudflare-ddns-updater"]