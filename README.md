# Courier Slot Coordinator

Courier Slot Coordinator is a small Go HTTP service for local delivery teams. Dispatchers create a shipment for a delivery zone and time window, assign a courier, record collection, and then confirm delivery. The service also provides a per-zone operational summary.

## Project structure

- `cmd/server`: HTTP server entry point.
- `internal/domain`: shipment types, lifecycle rules, and domain errors.
- `internal/repository`: concurrency-safe in-memory persistence.
- `internal/service`: dispatch workflows and reporting rules.
- `internal/api`: JSON HTTP handlers.

## Run

```bash
go run ./cmd/server
```

The server listens on `:8080` by default. Set `PORT` to choose another port.

## API

Create a shipment:

```bash
curl -X POST http://localhost:8080/shipments \
  -H 'Content-Type: application/json' \
  -d '{"recipient":"Ava","zone":"north","window":"09:00-11:00","packages":[{"sku":"book","units":2}]}'
```

Assign a courier:

```bash
curl -X POST http://localhost:8080/shipments/<shipment-id>/courier \
  -H 'Content-Type: application/json' \
  -d '{"courier_id":"courier-7"}'
```

Record collection or delivery with `POST /shipments/<shipment-id>/collect` and `POST /shipments/<shipment-id>/deliver`.

View a zone summary:

```bash
curl http://localhost:8080/zones/north/summary
```

## Build and test

```bash
go build ./...
go test ./...
```
