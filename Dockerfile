# builder image
FROM golang:1.11.5 as builder
# set working DIR
WORKDIR /go/src/github.com/dhyaniarun1993/foody-catalog-service
# COPY the codebase in the working DIR
COPY . .
# Install dep package for dependency management
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
# install all the dependency
RUN dep ensure -v
# Build catalog server
RUN go build -o ./cmd/catalog-server/main -v ./cmd/catalog-server

# generate clean, final image for end users
FROM golang:1.11.5
COPY --from=builder /go/src/github.com/dhyaniarun1993/foody-catalog-service/cmd/catalog-server/main /catalog-server
ENTRYPOINT [ "/catalog-server" ]