FROM golang:1.16-alpine

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY . ./

RUN go mod tidy
RUN go mod download
RUN go build -o bin/cyp2px

CMD [ "./bin/cyp2px" ]
