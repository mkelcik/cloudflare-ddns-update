FROM golang:1.20 as build

# Copy project sources
COPY . /opt/project/
WORKDIR /opt/project

RUN apt update
RUN apt-get install -y ca-certificates

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /cloudflare-ddns-updater

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /cloudflare-ddns-updater /cloudflare-ddns-updater
ENTRYPOINT ["/cloudflare-ddns-updater"]
CMD ["cloudflare-ddns-updater"]