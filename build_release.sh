#!/bin/bash

cd ./objects

go get -d -u
go get -t github.com/Lunkov/go-ecos-client/objects
go test

cd ..


cd ./messages

go get -d -u
go get -t github.com/Lunkov/go-ecos-client/objects
go test

cd ..

go get -d -u
go get -t github.com/Lunkov/go-ecos-client/objects
go test
