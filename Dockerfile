FROM golang:alpine

#COPY ./swaggerui ./swaggerui
COPY ./configs ./configs
COPY bin/main /main

CMD ["/main"]