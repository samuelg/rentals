FROM golang:1.20-alpine

WORKDIR /app

# default environment variables
ENV HOST 0.0.0.0
ENV PORT 8080
ENV GIN_MODE release
ENV ENV production
ENV LOG_LEVEL INFO
ENV DB_HOST localhost
ENV DB_PORT 5432
ENV DB_USER root
ENV DB_PASSWORD root
ENV DB_NAME testingwithrentals

# install dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy all files (.dockerignore set to ignore some files / dirs)
COPY . .

RUN go build -o rentals

# setup waiting for POSTGIS to come up
COPY --from=ghcr.io/ufoscout/docker-compose-wait:latest /wait /wait

# infor wait script of command to run once done waiting
ENV WAIT_COMMAND="/app/rentals"
ENTRYPOINT ["/wait"]

# default port
EXPOSE 8080
