FROM golang:1.13

# Copy gcp credentials json
ENV GOOGLE_APPLICATION_CREDENTIALS /workspace/auth/client_secret.json

# Move to working directory
WORKDIR /app/

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum . 
RUN go mod download

# Copy the code into the container
COPY . .

# Adds credentials
ADD ./creds.json $GOOGLE_APPLICATION_CREDENTIALS

# Build the application
RUN go build -o main.go .

# Execute main
ENTRYPOINT ["/main"]