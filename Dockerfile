FROM golang:1.22 AS builder
WORKDIR /app
RUN chmod 777 /app && echo "will going to build"
COPY . .
ENV CGO_ENABLED=1
# ENV GinMode=release
RUN go mod download && go build  -o /app/main .
EXPOSE 7860
RUN echo "build done"
CMD ["/app/main"]