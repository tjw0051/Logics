FROM golang:1.7-alpine

# Copy server to go workspace
ADD . /go/src/github.com/tjw0051/log-go

#Go get needs git (not on alpine image)
RUN apk add --update git

# Download Dependencies
RUN go get github.com/gin-gonic/gin
RUN go get github.com/jinzhu/gorm
RUN go get github.com/jinzhu/gorm/dialects/postgres

# Install Server
RUN go install github.com/tjw0051/log-go
RUN ls /go/src/github.com/tjw0051/log-go
# RUN cp /go/src/github.com/tjw0051/smart-voicemail-server/production.config /go/bin/

ENV PRODUCTION true

WORKDIR /go/bin

# Document that the service listens on port 80.
EXPOSE 80

# Run by default when the container starts.
ENTRYPOINT log-go
