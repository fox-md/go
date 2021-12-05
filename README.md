# go cheatsheet

Run tests 
```
go clean -testcache && go test ./... -v -cover -coverprofile cover.out
```
Convert tests to html
```
go tool cover -html=cover.out -o cover.html
```
