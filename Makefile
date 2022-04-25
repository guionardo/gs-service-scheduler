setup:
	go install golang.org/x/tools/cmd/cover

coverage:	
	go test -v -coverprofile=.coverage ./... && go tool cover -func=.coverage && unlink .coverage