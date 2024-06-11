# Dynamic DNS
This is a simple project that aims to help home labs and servers to automatically update their IP address to all the available `A` names in the given domain automatically.

# Preamble
I was using a domain from GoDaddy, but unfortunately GoDaddy has removed API access to customers with less than 50 domains, hence I had to change my nameservers to digital ocean, who doesn't have any limitations to their APIs and also their service is good in my usage.
<br>
`THIS WILL ONLY WORK IF YOU USE DIGITAL OCEAN FOR YOUR NAME SERVERS`

# Build Instructions
```bash
go build -ldflags "-X main.domain=YOUR_DOMAIN_HERE -X main.token=YOUR_DIGITAL_OCEAN_TOKEN_HERE"
```
