FROM golang:latest as builder

# For placing build source and artifacts.
RUN mkdir -p /go/{src,bin,pkg}

ADD . /go/src/github.com/clickyotomy/network-check
WORKDIR /go/src/github.com/clickyotomy/network-check

# Setup `go' stuff.
ENV GOBIN="/go/bin" GO111MODULE="on" CGO_ENABLED="0"

# Set any extra build tags.
ARG TAG

# Build the binary.
RUN make ship TAG=${TAG} && rm -rf /go/{src,bin,pkg}

# Create a image (lightweight AF).
FROM scratch
COPY --from=builder /tmp/pod-network-check /bin/pod-network-check

ENTRYPOINT ["/bin/pod-network-check"]
CMD ["-h"]
