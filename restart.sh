#!/bin/bash
echo "go build"
go mod tidy
go build -o main main.go
chmod +x ../gotklivedemo

res=$(netstat -anp|grep 8000|grep -v grep|wc|awk '{print $1}')
if [ "$res" -gt 1 ]; then
  echo "http start $res"
  supervisorctl restart tklivedemo_main
  echo "run gotklivedemo success"
else
   supervisorctl restart tklivedemo_main
   echo "run gotklivedemo success"
fi

