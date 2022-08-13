#!/bin/bash
cd user/userpb/
protoc --go_out=. --go-grpc_out=. user.proto 
cd ..
cd .. 
cd advice/advicepb/
protoc --go_out=. --go-grpc_out=. advice.proto