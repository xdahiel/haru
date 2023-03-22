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

  # shellcheck disable=SC2164
  cd ./src
  go mod tidy
  go build -o haru .
  # shellcheck disable=SC2103
  cd ..
  mv ./src/haru .
  ./haru -jwt=$res
}

rand