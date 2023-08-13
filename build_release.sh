#!/bin/bash

cd ./utils

rm go.mod go.sum
go mod init github.com/Lunkov/go-ecos-client/utils
go get .
go get -t github.com/Lunkov/go-ecos-client/utils
go test

cd ..

cd ./objects

go mod init github.com/Lunkov/go-ecos-client/objects
go get .
go get -t github.com/Lunkov/go-ecos-client/objects
go test

cd ..


cd ./messages

go mod init github.com/Lunkov/go-ecos-client/messages
go get .
go get -t github.com/Lunkov/go-ecos-client/messages
go test

cd ..

rm go.mod go.sum
go mod init github.com/Lunkov/go-ecos-client
go get .
go get -t github.com/Lunkov/go-ecos-client
go test
