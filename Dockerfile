FROM golang:1.20-alpine

# WORKDIR creates the specified directory inside the container.
# All futher commands are executed relative to the WORKDIR
WORKDIR /app

# COPY go.mod & go.sum into the WORKDIR before downloading dependencies 
# to utilize caching and reduce image build times
COPY go.mod go.sum ./
RUN go mod download

# COPY all go source files into the WORKDIR
COPY . .

# Build the executable
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o ./listr-backend ./cmd/listr-server/


ENV DOMAIN_SELF_SIGNED_CERTS_PATH=/root/.ssh/tls_certs

# Expose port 8080
EXPOSE 8080

# Primary command to execute the server that just got built
CMD [ "./listr-backend" ]
