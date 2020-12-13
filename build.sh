#!/bin/bash

usage="Usage: build.sh [-t|-h] \n -t (mac|linux|win|all)\n -h show help"

version=$(cat ./version.txt)

if [ -z "$version" ];then
	version="v1.0.5"
fi

cmd=""

function win() {
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./build/rpc-router-win64_$version.exe main.go
	echo "生成win的执行文件成功"
}

function linux() {
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/rpc-router-linux64_$version main.go
	echo "生成linux的执行文件成功"
}

function mac() {
	 CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -o ./build/rpc-router-darwin64_$version main.go
	echo "生成mac的执行文件成功"
}


while getopts ":t:h:" opt
do
    case $opt in
        t)
        cmd="$OPTARG"
	;;
	h)
	echo "$usage"
	exit 1;;
        ?)
        echo -e "$usage"
        exit 1;;
    esac
done

case $cmd in
	(win)
	win
	;;
	(linux)
	linux
	;;
	(mac)
	mac
	;;
	(all)
	win
	linux
	mac
	;;
	(*)
	echo -e "$usage"
	exit 1;;
esac
