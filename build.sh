#!/bin/sh

GIT_TAG=`git describe --always --dirty`                                                                                                                                                    
BTIME=`date -u +%s`                                                                                                                                                                        
# -w and -s diasables debugging stuff leading to a                                                                                                                                         
# reduction of binaries sizes/                                                                                                                                                             
#godep go build -ldflags "-w -X main.Version=${GIT_TAG} -X main.BuildTime=${BTIME}" -o ./build/bin/twiph
go build -ldflags "-w -X main.Version=${GIT_TAG} -X main.BuildTime=${BTIME}" -o ./build/bin/twiph
cp -r ./templates ./build
tar -czvf /tmp/linux_X86-64.tar.gz ./build
# Build for Mac OSX
env GOOS=darwin GOARCH=amd64 go build -ldflags "-w -X main.Version=${GIT_TAG} -X main.BuildTime=${BTIME}" -o ./build/bin/twiph
tar -czvf /tmp/darwin_X86-64.tar.gz ./build
