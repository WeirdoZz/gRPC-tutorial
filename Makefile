gen:
	protoc --go_out=plugins=grpc:pb proto/*.proto
clean:
	del pb/*.go
run:
	go run main.go
