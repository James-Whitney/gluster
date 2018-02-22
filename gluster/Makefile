all: client
	go build ./src/*.go

client:
	go build -o slave ./src/slave/slave.go

clean:
	rm runner test