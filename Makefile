build:
	go build -o bin/s3upload main.go

run:
	go run main.go

compile:
	# FreeBDS
	GOOS=freebsd GOARCH=amd64 go build -o bin/s3upload-freebsd-amd64 main.go
	# MacOS
	GOOS=darwin GOARCH=amd64 go build -o bin/s3upload-darwin-amd64 main.go
	# Linux
	GOOS=linux GOARCH=amd64 go build -o bin/s3upload-linux-amd64 main.go
	# Windows
	GOOS=windows GOARCH=amd64 go build -o bin/s3upload-windows-amd64 main.go