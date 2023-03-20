#!/usr/bin/env bash

dict="qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM0123456789"

function rand() {
  res=""
  for((i=1;i<=20;i++));
  do
    # shellcheck disable=SC2006
    # shellcheck disable=SC2003
    num=`expr $RANDOM % 62`
    res=$res${dict:num:1}
  done

  cd ./src
  go mod tidy
  go build -o haru .
  cd ..
  ./src/haru -jwt=$res
}

rand