# GitHub:       https://github.com/paper2code-bot
FROM golang:1.15-alpine AS build

ARG VERSION
ARG GIT_COMMIT
ARG BUILD_DATE

ARG CGO=1
ENV CGO_ENABLED=${CGO}
ENV GOOS=linux
ENV GO111MODULE=on

WORKDIR /go/src/github.com/paper2code-bot/stargazer

COPY . /go/src/github.com/paper2code-bot/stargazer/

# gcc/g++ are required to build SASS libraries for extended version
RUN apk update && \
    apk add --no-cache gcc musl-dev git ca-certificates make

RUN go build -ldflags "-extldflags=-static -extldflags=-lm" -o /go/bin/stargazer


FROM alpine:3.12

COPY --from=build /go/bin/stargazer /usr/bin/stargazer

# libc6-compat & libstdc++ are required for extended SASS libraries
# ca-certificates are required to fetch outside resources (like Twitter oEmbeds)
RUN apk update && \
    apk add --no-cache ca-certificates nano bash

VOLUME /data
WORKDIR /data

# Expose port for live server
EXPOSE 8898

CMD ["/usr/bin/stargazer", "--admin"]
