all: client
	go build *.go

client:
	go build -o runner slave/runner.go


clean:
	rm runner test