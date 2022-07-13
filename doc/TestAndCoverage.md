

```
go test -v  -coverprofile coverage.out ./... -coverpkg=./... && go tool cover -html coverage.out -o coverage.html
```