MODULE ?= broadcast-3a
BINARY ?= broadcast

run:
	go run ./$(MODULE)

build:
	go build -o ./$(MODULE)/$(BINARY) github.com/DeamonDev/gossip-glomers-$(MODULE)
