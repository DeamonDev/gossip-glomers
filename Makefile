MODULE ?= echo

run:
	go run ./$(MODULE)

build:
	go build -o ./$(MODULE)/$(MODULE) github.com/DeamonDev/gossip-glomers-$(MODULE)
