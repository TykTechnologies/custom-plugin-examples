# How to start gRPC server:


```bash
cp .env.example .env
```


modify .env



```bash
docker-compose up --build
```

Tyk configurations:

1) tyk.conf:

```
"coprocess_options": {
  "enable_coprocess": true,
  "coprocess_grpc_server": "tcp://host.docker.internal:9111"
},
```

ensure that tyk can reach the gRPC server

2) Use apidef.json in tyk

# Example DB config
Create a collection in mongo with this object (I have tested this with Loopback 3 `AccessToken` model which already has compatible schema)
```json
{
  "_id":"asdfasdfasdf32dsafdasdfas",
  "ttl":1209600,
  "created":"2022-08-30T09:20:20.910+00:00"
}
```

# Testing:
```bash
curl -s 'http://localhost:8082/mongo-auth/get' -H "Authorization: asdfasdfasdf32dsafdasdfas"
```
Output if token is valid:
```
{
  "args": {}, 
  "headers": {
    "Accept": "*/*", 
    "Accept-Encoding": "gzip", 
    "Authorization": "asdfasdfasdf32dsafdasdfas", 
    "Host": "httpbin.org", 
    "User-Agent": "curl/7.81.0", 
    "X-Amzn-Trace-Id": "Root=1-6315b555-4a26f3f923cde5aa3f596b6f"
  }, 
  "origin": "172.28.0.1, <external ip>", 
  "url": "http://httpbin.org/get"
}
```

Output if token is invalid:
```
{
    "error": "Access forbidden"
}
```