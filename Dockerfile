FROM ubuntu:22.04

WORKDIR /music-store

RUN apt-get update && apt-get install -y ca-certificates
RUN apt-get clean &&   rm -rf /var/lib/apt/lists/* 

COPY models/ models/
ADD bin/music_store.bin /music-store/music_store.bin
