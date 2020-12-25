FROM golang:1.15-buster as builder

WORKDIR /app
RUN apt-get update
RUN apt-get install -y -qq libtesseract-dev libleptonica-dev
ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/4.00/tessdata/
RUN apt-get install -y -qq \
  tesseract-ocr-eng \
  tesseract-ocr-tha
COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -mod=readonly -v -o lineshark

FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*
RUN apt-get update
RUN apt-get install -y -qq libtesseract-dev libleptonica-dev
ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/4.00/tessdata/
RUN apt-get install -y -qq \
  tesseract-ocr-eng \
  tesseract-ocr-tha

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/lineshark /app/lineshark
COPY --from=builder /app/template /app/template
COPY --from=builder /app/evidence /app/evidence
WORKDIR /app

# Run the web service on container startup.
CMD ["./lineshark"]
