.PHONY: hbrender

hbrender:
	go build -o bin/hbrender cmd/hbrender/main.go

install:
    go install github.com/ambientsound/hbrender/cmd/hbrender