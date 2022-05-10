gen:
	protoc --go-grpc_out=./pb ./proto/*.proto
	protoc --go_out=./pb ./proto/*.proto
clean:
	rm  pb/*.go
server:
	go run cmd/server/main.go -port 8080
client:
	go run cmd/client/main.go -address 0.0.0.0:8080
evans:
	evans -r repl -p 8080
run:
	go run main.go

