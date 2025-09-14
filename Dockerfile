# Stage 1: Common builder with all dependencies
FROM golang:1.23-bullseye as builder

# Install all necessary dependencies for OpenCV and the Go build
# Using non-interactive to prevent prompts
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y \
    build-essential cmake git pkg-config \
    libgtk-3-dev libavcodec-dev libavformat-dev libswscale-dev \
    libv4l-dev libxvidcore-dev libx264-dev libjpeg62-turbo-dev \
    libpng-dev libtiff-dev gfortran openexr libatlas-base-dev \
    python3-dev python3-numpy libtbb2 libtbb-dev libdc1394-dev

# Download and compile OpenCV
WORKDIR /opt
RUN git clone --branch 4.11.0 https://github.com/opencv/opencv.git
RUN git clone --branch 4.11.0 https://github.com/opencv/opencv_contrib.git
WORKDIR /opt/opencv/build
RUN cmake -D CMAKE_BUILD_TYPE=Release \
          -D CMAKE_INSTALL_PREFIX=/usr/local \
          -D OPENCV_GENERATE_PKGCONFIG=ON \
          -D OPENCV_EXTRA_MODULES_PATH=/opt/opencv_contrib/modules \
          -D BUILD_EXAMPLES=OFF ..
RUN make -j$(nproc) && make install

# Set up environment for CGO
ENV PKG_CONFIG_PATH="/usr/local/lib/pkgconfig"

# Copy Go modules and source code
WORKDIR /app
COPY . .

# Download dependencies and add the missing one
RUN go mod download
RUN go get github.com/aws/aws-lambda-go/lambda

# Build the API binary
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /app/api-bootstrap ./cmd/app/main.go

# Build the Cronjob binary
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /app/cron-bootstrap ./cmd/cronjob/main.go

# --------------------------------------------------------------------

# Final stage for the API function
FROM amazonlinux:2023 as api-v1

# Install required runtime libraries
RUN dnf install -y freetype fontconfig && dnf clean all

# Copy the necessary shared libraries from the builder stage
COPY --from=builder /usr/local/lib /usr/local/lib

# Copy the compiled API binary from the builder stage
COPY --from=builder /app/api-bootstrap /app/bootstrap

# Copy the email templates into the final image
COPY --from=builder /app/internal/emails/templates /app/internal/emails/templates

# Set the working directory
WORKDIR /app

# Set the command to run the binary
CMD ["./bootstrap"]

# --------------------------------------------------------------------

# Final stage for the Cronjob function
FROM public.ecr.aws/lambda/provided:al2023 as cron-v1

# Copy the necessary shared libraries from the builder stage
COPY --from=builder /usr/local/lib /usr/local/lib

# Copy the compiled Cronjob binary from the builder stage
COPY --from=builder /app/cron-bootstrap /var/runtime/bootstrap

# Add a unique identifier to make this image different
RUN echo "CRONJOB_IMAGE" > /var/runtime/.image_type

# Set the command to run the binary
CMD ["/var/runtime/bootstrap"]
