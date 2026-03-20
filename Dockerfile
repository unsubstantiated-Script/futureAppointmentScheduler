FROM ubuntu:latest
LABEL authors="skripty"

ENTRYPOINT ["top", "-b"]