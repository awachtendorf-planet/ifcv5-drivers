## proxy.toml

```toml
# access blacklist file name, path must exist
blacklist = "./.config/blacklist.txt"

[listen]
port = 6000
#host = "127.0.0.1"

[tls]
key = "./.proxy.key"
certificate = "./.proxy.crt"
insecure = true
#servername = "ifc.evilcorp"
#rootCAs = "/home/ullmann/.step/certs"

[api]
port = 7000
#host = "127.0.0.1"

# hardcoded tls
# "./.api-cert.pem"
# "./.api-key.pem"

[metric]
port = 8000

[redis]
host = "127.0.0.1"
port = 11001
auth = "topsecret"
usetls = true

[redis.tls]
insecure = false
servername = "kvdb"
#rootCAs = ["/home/ullmann/.local/share/mkcert/rootCA.pem"]


```

## blacklist.txt

```
pcon-local	2021-07-08T12:00:00Z                // block peer 'pcon-local' until date
pcon-* 		2021-07-09T14:00:00Z 217.6.121.0/24 // block peers starts with 'pcon-' until date AND peer ip is in subnet
pcon-NO* 	217.6.121.164                       // block peers starts with 'pcon-NO' with exact ip address match
pcon-SE* 	                                    // block peers starts with 'pcon-SE'
```