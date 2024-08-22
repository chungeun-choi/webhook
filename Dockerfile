FROM alpine:latest

# Install ca-certificates
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /app


# ARG BINARY=myapp
ARG BINARY
COPY ${BINARY} /app/myapp

# TODO: Change time that copy the app.yml and token.txt
COPY app.yml /app/app.yml
COPY token/token.txt /app/token.txt

# This container exposes port 8080 to the outside world
RUN chmod +x /app/myapp

# Run the binary program produced by `go install`
EXPOSE 8080
EXPOSE 443

# Command to run the executable
CMD ["./myapp"]
