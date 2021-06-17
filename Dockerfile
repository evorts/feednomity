FROM golang:1.15.2-alpine as builder

LABEL Maintainer="Evorts Technology"

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR /apps/

COPY go.* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /go/bin/app .

FROM alpine:latest

ARG USER_ID
ARG GROUP_ID

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
#COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/app /go/bin/app
COPY --from=builder /apps/tmpl /go/bin/tmpl
COPY --from=builder /apps/tmpl_mail /go/bin/tmpl_mail
COPY --from=builder /apps/assets /go/bin/assets
COPY --from=builder /apps/forms /go/bin/forms
COPY --from=builder /apps/config.docker.yml /go/bin/config.yml

WORKDIR /go/bin/

ENV USER=appuser
ENV UID=1001
ENV GID=1000
RUN addgroup --gid "$GID" "$USER"
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "$(pwd)" \
    --ingroup "$USER" \
    --no-create-home \
    --uid "$UID" \
    "$USER"
USER "$USER"

ENV TZ=Asia/Jakarta


