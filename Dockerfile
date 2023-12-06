# Build stage
FROM golang:latest AS build
WORKDIR /app
COPY . .
RUN go build -o bookings cmd/web/*.go

# Production stage
FROM golang:latest
WORKDIR /app
COPY --from=build /app/bookings .
CMD ["./bookings"]
