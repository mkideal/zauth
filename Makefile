osx:
	GOOS=darwin go build -o authc ./cmd/authc/main.go
	GOOS=darwin go build -o authd ./cmd/authd/main.go
linux:
	GOOS=linux go build -o authc ./cmd/authc/main.go
	GOOS=linux go build -o authd ./cmd/authd/main.go
win:
	GOOS=windows go build -o authc.exe ./cmd/authc/main.go
	GOOS=windows go build -o authd.exe ./cmd/authd/main.go
