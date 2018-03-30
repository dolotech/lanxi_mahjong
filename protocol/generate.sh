#!/bin/bash
./protoc --go_out=../src/protocol ./*.proto
read -p "Press any key to continue." var
