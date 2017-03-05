osx:
	go build -o authc ./cmd/authc/main.go
	go build -o authd ./cmd/authc/main.go
linux:
	GOOS=linux go build -o authc ./cmd/authc/main.go
	GOOS=linux go build -o authd ./cmd/authc/main.go
