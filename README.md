# Wii/DS DLS Server

OpenPak's Wii/DS DLS Server.

| | |
| --- | --- |
| Network | Wii, DS |
| Role | DLS1 downloadable-content request codec and native download counts |
| Depends on | WFC database; no remote WFC service |
| Image | `ghcr.io/openpak/wfc-dls` |
| Based on | [WiiLink WFC](https://github.com/WiiLink24/wfc-server) (AGPL-3.0) |
| Licence | AGPL-3.0-only |

## What it serves

| Endpoint | What |
| --- | --- |
| `POST /download` | DLS1 request codec, `count` action (native download counts) |
| `GET /healthz` | Unauthenticated health check |

## Build

```sh
GOWORK=off go test -mod=vendor ./...
docker build -t wfc-dls .
```

## Docs

- [Protocol](docs/protocol.md)
- [Evidence](docs/evidence.md)
- [Deploy](docs/deploy.md)
- [Notes](docs/notes.md)

## Licence

AGPL-3.0-only. See [LICENSE](LICENSE).
