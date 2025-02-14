MODULE ?= broadcast-3b
BINARY ?= broadcast

run:
	go run ./$(MODULE)

build:
	go build -o ./$(MODULE)/$(BINARY) github.com/DeamonDev/gossip-glomers-$(MODULE)
