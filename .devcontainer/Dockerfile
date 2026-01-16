# Note: You can use any Debian/Ubuntu based image you want. 
FROM golang:1.26rc2

RUN apt update
RUN apt install -y jq
RUN go install github.com/goreleaser/goreleaser/v2@latest
RUN go install github.com/spf13/cobra-cli@latest
