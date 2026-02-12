# Go Latency Patterns

A collection of practical Go examples demonstrating patterns and techniques for reducing latency in backend services. Each pattern lives in its own directory with a fully working demo, Docker Compose stack, and monitoring via Prometheus + Grafana.

## Patterns

| Pattern | Description | Directory |
|---------|-------------|-----------|
| **Singleflight** | Deduplicate concurrent requests to prevent cache stampedes and protect the database from redundant load | [`singleflight/`](singleflight/) |
