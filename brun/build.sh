#!/bin/bash

SOURCE_FILE_NAME=main #源文件名称
TARGET_FILE_NAME=reskd #目标文件名称

rm -rf ${TARGET_FILE_NAME}*

build(){ #构建脚本
  echo $GOOS $GOARCH
  tname=${TARGET_FILE_NAME}_${GOOS}_${GOARCH}${EXT}
  # 编译 反斜杠连接命令 下面三行命令其实是在一行运行的
  env GOOS=${GOOS} GOARCH=${GOARCH} \
  go build -o ${tname} \
  -v ${SOURCE_FILE_NAME}.go
  #添加可执行权限
  chmod +x ${tname}
  mv ${tname} ${TARGET_FILE_NAME}${EXT}
  #打包
  if [ ${GOOS} = "windows" ];then
    zip ${tname}.zip ${TARGET_FILE_NAME}${EXT} config.ini ../public/*
  else
    tar --exclude=*.gz --exclude=*.zip --exclude=*.git -zcvf ${tname}.tar.gz ${TARGET_FILE_NAME}${EXT} config.ini *.sh ../public/ -C ./ .
  fi

  mv ${TARGET_FILE_NAME}${EXT} ${tname}
}

CGO_ENABLED=0

# linux操作系统
GOOS=linux
 # amd64 架构
GOARCH=amd64
build # 执行编译

# mac os操作系统
GOOS=darwin
 # amd64 架构
GOARCH=amd64
build # 执行编译

# windows操作系统
GOOS=windows
 # amd64 架构
GOARCH=amd64
build # 执行编译