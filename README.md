# go cheatsheet

Init module
```
go mod init MODULENAME
```
Install all deps
```
go get ./...
```
Run tests 
```
go clean -testcache && go test ./... -v -cover -coverprofile cover.out
```
Run single test with 10 min timeout
```
go test -run TestR53 -timeout 10m
```
This will show you the code coverage for every single package withing your project and down at the bottom the total coverage.
```
go tool cover -func cover.out
```
Convert tests to html
```
go tool cover -html=cover.out -o cover.html
```

Enable linux build on Windows
```
set GOARCH=amd64
set GOOS=linux
go tool dist install -v pkg/runtime
go install -v -a std
```
