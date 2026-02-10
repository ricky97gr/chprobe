# Use smaller base image as runtime environment
FROM alpine:latest

# Install necessary tools
#RUN apk add --no-cache ca-certificates tzdata

# Set timezone
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" > /etc/timezone

# Create working directory with writeable layer
WORKDIR /app

# Create directory for client downloads
RUN mkdir -p /app/client_downloads && chmod 777 /app/client_downloads

# Copy compiled server binary from local directory
COPY chprobe_server/bin/chprobe_server ./

# Copy compiled client binary from local directory
COPY chprobe_client/bin/chprobe_client /app/client_downloads/

# Copy configuration files
COPY chprobe_server/config/config.yaml ./config/

# Expose ports
EXPOSE 32000 32001

# Set startup command
CMD ["./chprobe_server"]
