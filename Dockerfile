FROM golang:1.20.4-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o main.exe main.go

LABEL traefik.port=80
LABEL traefik.http.routers.tt.rule="Host(`timetable.mlsctiet.com`)"
LABEL traefik.http.routers.tt.tls=true
LABEL traefik.http.routers.tt.tls.certresolver="lets-encrypt"
LABEL org.opencontainers.image.source="https://github.com/utkarsh-1905/time-table"

EXPOSE 3000

CMD [ "./main.exe" ]

# FROM scratch
# WORKDIR /app
# COPY --from=build /app/main.exe ./
# COPY --from=build /app/data.json ./
# COPY --from=build /app/timetable.xlsx ./
# # if deployed on personal server
# EXPOSE 3000
# CMD [ "./main.exe" ]
