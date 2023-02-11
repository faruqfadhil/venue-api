FROM golang:alpine

RUN mkdir /app

WORKDIR /app

ADD go.mod .
ADD go.sum .

RUN go mod download
ADD . .

# add docker-compose-wait tool to wait db init container ready before run other container.
ENV WAIT_VERSION 2.7.2
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait
RUN chmod +x /wait

# Build
RUN go build -o /venue-api

EXPOSE 8081

# Run
CMD [ "/venue-api" ]
