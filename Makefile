gen:
	protoc --go_out=./pb ./proto/*.proto
clean:
	rm  pb/*.go
run:
	go run main.go
