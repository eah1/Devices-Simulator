# Build the Go Binary.
FROM golang:1.19.1-alpine AS builder

ENV CGO_ENABLED 0
ARG BUILD_REF

ARG GOARCH

# Copy the source code into the container.
COPY . /Device-Simulator

# Build the service binary.
WORKDIR /Device-Simulator/app/services/simulator-api

# Run the Go Binary in Alpine.
RUN GOARCH=${GOARCH}  go build -ldflags "-X main.build=${BUILD_REF}"

FROM alpine:latest
RUN apk --no-cache add tzdata

ARG BUILD_DATE
ARG BUILD_REF
ENV BUILD_REF=${BUILD_REF}

COPY --from=builder  /Device-Simulator/app/services/simulator-api/ /services/simulator-api
COPY --from=builder /Device-Simulator/business/template/ /services/simulator-api/business/template/
WORKDIR /services/simulator-api

CMD ["./simulator-api"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="simulator-api" \
      org.opencontainers.image.vendor="Circutor"