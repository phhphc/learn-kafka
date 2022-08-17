all:
	go build -o ./cmd/ consumer/*.go
	go build -o ./cmd/ producer/*.go

makeproducer:
	go build -o ./cmd/ producer/*.go

makeconsumer:
	go build -o ./cmd/ producer/*.go
	

