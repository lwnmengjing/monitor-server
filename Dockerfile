FROM alpine

LABEL authors="lwnmengjing"

COPY ./application /app/application

ENTRYPOINT ["/app/application"]