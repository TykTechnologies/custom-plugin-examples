```bash
rm -rf go.mod go.sum vendor/
go1.16.15 mod init tyk-grpc-plugin
go1.16.15 mod edit -replace github.com/jensneuse/graphql-go-tools=github.com/TykTechnologies/graphql-go-tools@v1.6.2-0.20221026084245-1fc4f5ca74bb
go1.16.15 get github.com/TykTechnologies/tyk@767c8b336fa5433174648bc2532edd7611ce9be0
go1.16.15 get go.mongodb.org/mongo-driver
<Write and finalize code>
go1.16.15 mod tidy
go1.16.15 mod vendor
```