FROM golang:1.23.4
WORKDIR /app
# Define build arguments
ARG DB_USER
ARG DB_PW
ARG DB_HOST
ARG DB_PORT
ARG DB_NAME

# Set as environment variables
ENV DB_USER=$DB_USER \
    DB_PW=$DB_PW \
    DB_HOST=$DB_HOST \
    DB_PORT=$DB_PORT \
    DB_NAME=$DB_NAME

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o main ./server
CMD ["./server/main"]