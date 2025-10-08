MODULE ?= echo
BINARY ?= echo

run:
	go run ./$(MODULE)

build:
	go build -o ./$(MODULE)/$(BINARY) github.com/DeamonDev/gossip-glomers-$(MODULE)
