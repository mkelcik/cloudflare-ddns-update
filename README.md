![Code and security checks](https://github.com/mkelcik/cloudflare-ddns-update/actions/workflows/quality-checks.yml/badge.svg)

## What is Cloudflare Dynamic DNS?
DNS records are static, and it does not play well with dynamic IP addresses. Now, to solve that problem, you’ll need to set up dynamic DNS. Cloudflare provides an API that allows you to manage DNS records programmatically.

To set up a Cloudflare dynamic DNS, you’ll need to run a process on a client inside your network that does two main actions: get your network’s current public IP address and automatically update the corresponding DNS record.

This simple updater do the job.

## How to run
### Environment variables

Before run, you need configure this environment variables.

 - `CLOUDFLARE_DNS_TO_CHECK` - (required) dns records that will be automatically checked and modified based on the current public IP address. Multiple entries are separated by commas. For example: `domain.com,sub1.domain.com,sub2.domain.com`
 - `CLOUDFLARE_API_KEY` - (required) your cloudflare api key, with access rights to edit selected domains. See: [https://developers.cloudflare.com/fundamentals/api/get-started/create-token/](https://developers.cloudflare.com/fundamentals/api/get-started/create-token/)
 - `CLOUDFLARE_ZONE` - (required) zone name with domain you want to check. See: [https://developers.cloudflare.com/fundamentals/get-started/concepts/accounts-and-zones/#zones](https://developers.cloudflare.com/fundamentals/get-started/concepts/accounts-and-zones/#zones)
 - `ON_CHANGE_COMMENT` - (optional) in the event that the ip address of the dns record changes, this comment will be added to the record
 - `CHECK_INTERVAL_SECONDS` - (optional) how often will the ip address of the records be checked (default: `300`)
 - `PUBLIC_IP_RESOLVER` - (optional) public ip address resolver. (default: `ifconfig.me`) Available: `ifconfig.me`, `v4.ident.me`

### Building from source

```shell
  go build -o /cloudflare-ddns-updater
```

### Go install

Install via go install

```shell
go install github.com/mkelcik/cloudflare-ddns-update
```

### Running

```shell
CLOUDFLARE_DNS_TO_CHECK="domain.com" CLOUDFLARE_API_KEY="my_key" CLOUDFLARE_ZONE="domain.com" cloudflare-ddns-update
```

### Via `docker-compose`
```yaml
version: "3"
services:
  cf-dns-updater:
    image: mkelcik/cloudflare-ddns-update:latest
    restart: unless-stopped
    environment:
      - CLOUDFLARE_DNS_TO_CHECK=my.testdomain.com,your.testdomain.com
      - CLOUDFLARE_API_KEY=your_cloudflare_api_key
      - CLOUDFLARE_ZONE=testdomain.com
      - ON_CHANGE_COMMENT="automatically updated"
      - CHECK_INTERVAL_SECONDS=300
```

### Via `docker run`
```shell
docker run -e CLOUDFLARE_DNS_TO_CHECK=my.testdomain.com,your.testdomain.com -e CLOUDFLARE_API_KEY=your_cloudflare_api_key -e CLOUDFLARE_ZONE=testdomain.com -e ON_CHANGE_COMMENT="automatically updated" -e CHECK_INTERVAL_SECONDS=300 mkelcik/cloudflare-ddns-update:latest 
```

### Contributing 

Feel free to contribute and pls report bugs. Thanks
