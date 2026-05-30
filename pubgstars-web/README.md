# pubgstars-web

Go backend for the PubgStars platform. Each Lambda function lives in its own `cmd/` directory and is deployed independently.

For full architecture, deployment, and infrastructure documentation see the [root README](../README.md).

---

## Requirements

| Tool | Version |
|---|---|
| Go | 1.22+ |
| AWS CLI | v2 |
| Docker | any recent version (integration tests only) |

---

## Project structure

```
pubgstars-web/
├── cmd/                        # One main.go per Lambda function
│   ├── games/
│   ├── registerToGame/
│   ├── unregisterToGame/
│   ├── adminCompleteGame/
│   └── ...
├── internal/
│   ├── AwsUtils.go             # DynamoDB client, JWT parsing, time helpers
│   ├── DataService.go          # DynamoDB read/write operations
│   ├── GameUtils.go            # Time-window logic (password access, cancellation)
│   ├── SlackService.go         # Slack notifications
│   ├── ModelUtils.go           # ID generation, misc utilities
│   ├── TransactionLogUtils.go  # Transaction log builders
│   ├── store.go                # Store interface (for dependency injection)
│   └── dynamo_store.go         # DynamoDB implementation of Store
├── model/
│   ├── Model.go                # Game, User, TransactionLog, Message structs
│   └── tables/Tables.go        # DynamoDB table name constants
├── testutil/
│   └── mock_store.go           # MockStore for handler unit tests
├── test/
│   └── game_test.go            # Model serialisation tests
├── scripts/
│   ├── buildAndUpload.sh       # Build + deploy a single Lambda
│   ├── buildAndUploadAll.sh    # Build + deploy all Lambdas
│   └── database.go             # DynamoDB table provisioning script
└── pkg/                        # Shared logger / printer utilities
```

---

## Local development

Activate pre-commit hooks once after cloning (runs `go build` + `go test` before every commit):

```bash
make setup
```

Build and test:

```bash
make build   # go build ./...
make test    # go test ./...
```

---

## Running tests

### Unit + handler tests

```bash
go test ./...
```

### Integration tests (requires Docker)

Start DynamoDB Local, run the integration suite, then tear down:

```bash
docker compose -f ../docker-compose.test.yml up -d
go test -tags integration ./internal/
docker compose -f ../docker-compose.test.yml down
```

Integration tests skip automatically if DynamoDB Local is not reachable.

---

## Deploying

Build and upload a single Lambda:

```bash
cd scripts
./buildAndUpload.sh games
```

Build and upload all Lambdas:

```bash
cd scripts
./buildAndUploadAll.sh
```

The scripts cross-compile to `linux/amd64`, zip the binary, and call `aws lambda update-function-code`.

---

## Environment variables

| Variable | Description |
|---|---|
| `SLACK_TOKEN` | Slack bot token (`chat:write` scope) |
| `CHANNEL_NAME` | Default Slack channel for notifications |
| `AWS_PROFILE` | AWS CLI profile name (local development only) |

Copy `.env.example` and fill in values for local runs.
