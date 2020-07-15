FROM golang

COPY ./main.go ./main.go

RUN go build -o exporter main.go

CMD ["./exporter"]