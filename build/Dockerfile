FROM golang as base
WORKDIR /app

RUN addgroup gouser && adduser --ingroup gouser --uid 19998 --shell /bin/false gouser

ARG TARGETARCH
ARG TARGETOS

COPY go.mod go.sum ./
COPY ./vendor ./vendor
COPY ./internal ./internal
COPY ./cmd ./cmd
COPY ./views ./views
COPY main.go ./

RUN --mount=type=cache,target="/root/.cache/go-build" \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GOCACHE=/root/.cache/go-build \ 
    go build -mod=vendor -o /app/fiber-htmx main.go


FROM scratch

COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /app/fiber-htmx /app/fiber-htmx
USER gouser

ENTRYPOINT [ "/app/fiber-htmx" ]