# PubgStars

A full-stack platform for organizing paid PUBG tournament sessions. Players deposit balance, register for scheduled game rooms, receive room credentials shortly before kick-off, and winners collect prize money — all handled automatically.

Built on a fully serverless AWS stack.

---

## Repository layout

```
pubgstars/
├── pubgstars-web/          # Go backend — AWS Lambda functions + shared library
└── pubgstars-web-static/   # React frontend — deployed to S3 / CloudFront
```

---

## How it works

### Player flow
1. Sign up via AWS Cognito; balance starts at zero.
2. Deposit funds through the Balance page (admin confirms payment manually).
3. Browse active game sessions on the Home page and register; the entry fee is deducted atomically.
4. Within the hour before the game starts, the **room name + password** becomes visible.
5. Play. After the session the admin marks the game complete and selects the top-3 winners.
6. Prize money is credited to winners' balances instantly.
7. Unregistering before the 1-hour cutoff refunds the entry fee as bonus credit.

### Admin flow
- Create / update / delete game sessions (`adminGames`).
- Manage users and manually top up balance or bonus (`adminAddBudget`).
- View registrations per game (`adminGameUsers`).
- Handle support messages (`adminMessages`).
- Complete a game and distribute prizes to 1st / 2nd / 3rd place (`adminCompleteGame`).

---

## Backend — `pubgstars-web`

**Language:** Go 1.22+  
**Runtime:** AWS Lambda (one binary per endpoint)  
**Database:** AWS DynamoDB  
**Auth:** AWS Cognito — JWT tokens validated in each handler  
**Notifications:** Slack API  
**Deployment:** Cross-compiled to `linux/amd64`, zipped, uploaded with `aws lambda update-function-code`

### Lambda functions

| Function | Method(s) | Description |
|---|---|---|
| `games` | GET / POST / PUT / DELETE | List active games; create or update a session |
| `gamePassword` | POST | Return room credentials (only within ±1 h of game time) |
| `gameUsers` | POST | List users registered for a game |
| `registerToGame` | POST | Register the caller; deducts entry fee atomically |
| `unregisterToGame` | POST | Unregister (refunds as bonus if before cutoff) |
| `gamesHistory` | GET | Completed game history |
| `gamesLeaderboard` | GET | Top-3 all-time winner leaderboard |
| `user` | GET | Current user profile |
| `userGames` | GET | Games the caller is registered for |
| `canUserRegister` | — | Cognito pre sign-up trigger: username availability check |
| `depositMoney` | POST | Record a deposit request |
| `withdrawMoney` | POST | Request a withdrawal (requires secret answer) |
| `transactionLog` | GET | Full transaction history for the caller |
| `sendMessage` | POST | Submit a support message |
| `adminGames` | GET / POST / PUT / DELETE | Admin: manage all game sessions |
| `adminCompleteGame` | POST | Admin: mark complete, credit prize money to winners |
| `adminGameUsers` | GET | Admin: view registrations per game |
| `adminUsers` | GET | Admin: list all users |
| `adminMessages` | GET / POST | Admin: view and respond to support messages |
| `adminAddBudget` | POST | Admin: manually credit balance / bonus to a user |
| `userRegistered` | — | Cognito post-confirmation trigger: creates user record |

### DynamoDB tables

| Table | Hash key | Purpose |
|---|---|---|
| `games` | `id` | Tournament sessions |
| `users` | `id` | User profiles (balance, bonus, gain) |
| `gameUsers` | `gameId` | User IDs registered per game |
| `userGames` | `userId` | Game IDs per user |
| `transactionLog` | `id` | Immutable audit trail for all operations |
| `messages` | `id` | Customer support threads |

### Key design decisions

- **Atomic writes everywhere.** Registration, unregistration, and prize distribution all use `TransactWriteItems` — game state, user balance, cross-reference tables, and the transaction log are updated in a single DynamoDB transaction. No partial state is possible.
- **Bonus-first spending.** When a user registers for a game, bonus credit is consumed before real balance.
- **Time-gated room access.** Room passwords are only returned when the current time falls within the [gameTime − 1 h, gameTime + 1 h] window.
- **Cancellation window.** Unregistration is blocked within 1 hour of game start.

### Project structure

```
pubgstars-web/
├── cmd/                        # One main.go per Lambda function
│   ├── games/
│   ├── registerToGame/
│   ├── adminCompleteGame/
│   └── ...
├── internal/
│   ├── AwsUtils.go             # DynamoDB client, JWT parsing, time helpers
│   ├── DataService.go          # DynamoDB read/write operations
│   ├── GameUtils.go            # Time-window logic (passwords, cancellation)
│   ├── SlackService.go         # Slack notifications
│   ├── ModelUtils.go           # ID generation, misc utilities
│   ├── TransactionLogUtils.go  # Transaction log builders
│   ├── store.go                # Store interface (dependency injection)
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
│   └── database.go             # DynamoDB table provisioning
└── pkg/                        # Shared logger / printer utilities
```

### Local development

After cloning, activate the pre-commit hooks (runs `go build` + `go test` before every commit):

```bash
cd pubgstars-web
make setup
```

```bash
make build   # go build ./...
make test    # go test ./...
```

### Testing

```bash
# Unit and handler tests (no external dependencies)
go test ./...

# Integration tests against DynamoDB Local
docker compose -f docker-compose.test.yml up -d
go test -tags integration ./internal/
docker compose -f docker-compose.test.yml down
```

| Test type | Coverage | Dependencies |
|---|---|---|
| Unit | Pure functions — time logic, JWT parsing, transaction builders, model | None |
| Handler unit | Lambda handler routing, input validation, error responses | MockStore |
| Integration | DynamoDB read/write — SaveGame, SaveUser, Register/Unregister, UpdateUserWithTx | DynamoDB Local (Docker) |

Integration tests skip automatically if DynamoDB Local is not running.

### Environment variables

| Variable | Description |
|---|---|
| `SLACK_TOKEN` | Slack bot token (`chat:write` scope) |
| `CHANNEL_NAME` | Default Slack channel for notifications |
| `AWS_PROFILE` | AWS CLI profile name (local development only) |

---

## Frontend — `pubgstars-web-static`

**Framework:** React 18 (Create React App)  
**Auth:** AWS Amplify + Cognito  
**HTTP:** Axios  
**UI:** React-Bootstrap  
**Deployment:** `npm run build` → `aws s3 sync build s3://pubgstars.com/`

### Pages

| Route | Component | Description |
|---|---|---|
| `/` | Home | Active game sessions; register / unregister; reveal room password |
| `/mygames` | MyGames | Games the logged-in user is registered for |
| `/leaderboard` | LeaderBoard | All-time rankings |
| `/balance` | Balance | Deposit and withdraw funds |
| `/transactionlog` | TransactionLog | Full transaction history |
| `/profile` | Profile | Account details |
| `/login` | Login | Cognito sign-in |
| `/signup` | Signup | Registration (username, email, password, secret Q&A) |
| `/login/reset` | ResetPassword | Cognito forgot-password flow |
| `/rules` | Rules | Tournament rules |
| `/sss` | Sss | FAQ |
| `/about` | About | About the platform |
| `/contact` | Contact | Contact / support form |

### Local development

```bash
cd pubgstars-web-static
cp .env.example .env      # fill in Cognito values
npm install
npm start                 # http://localhost:3000
```

Required environment variables (must be prefixed `REACT_APP_`):

| Variable | Description |
|---|---|
| `REACT_APP_STAGE` | Environment name (`dev` / `prod`) |
| `REACT_APP_COGNITO_REGION` | AWS region |
| `REACT_APP_COGNITO_USER_POOL_ID` | Cognito User Pool ID |
| `REACT_APP_COGNITO_APP_CLIENT_ID` | Cognito App Client ID |
| `REACT_APP_COGNITO_IDENTITY_POOL_ID` | Cognito Identity Pool ID |

---

## Infrastructure overview

```
Browser
  │
  ├─► S3 / CloudFront          (React SPA — pubgstars.com)
  │
  └─► API Gateway (api.pubgstars.com/v1)
        │
        ├─► Lambda (Go) ──► DynamoDB
        ├─► Lambda (Go) ──► Cognito (JWT validation)
        └─► Lambda (Go) ──► Slack API

Cognito User Pool ──► pre sign-up trigger    ──► canUserRegister Lambda
                  └─► post-confirmation trigger ──► userRegistered Lambda
```

All backend compute is Lambda; there are no long-running servers.

---

## Deploying to AWS

### Prerequisites

| Tool | Version | Notes |
|---|---|---|
| Go | 1.22+ | `brew install go` |
| Node.js | 18–20 | `react-scripts 5` does not support Node 21+; use `nvm use` (`.nvmrc` provided) |
| AWS CLI | v2 | `brew install awscli` |
| Docker | any | Integration tests only |

Configure an AWS CLI profile named `pg`:

```bash
aws configure --profile pg
# Default region: eu-central-1
```

---

### 1. Cognito User Pool

Create a Cognito User Pool with the following custom attributes:
- `custom:secretQuestion` (String)
- `custom:secretAnswer` (String)

Create an **App Client** (no client secret) and note the:
- User Pool ID
- App Client ID
- Identity Pool ID

Set these in `pubgstars-web-static/.env`.

---

### 2. DynamoDB Tables

```bash
cd pubgstars-web/scripts
AWS_PROFILE=pg go run database.go
```

Tables created (all in `eu-central-1`, on-demand billing):

| Table | Hash key | Range key |
|---|---|---|
| `games` | `id` | — |
| `users` | `id` | — |
| `gameUsers` | `gameId` | — |
| `userGames` | `userId` | — |
| `transactionLog` | `id` | `userId` |
| `messages` | `id` | `from` |

---

### 3. IAM Role for Lambda

Create an IAM role named `lambda-service-role` with:
- Trust policy: `lambda.amazonaws.com`
- Policies: `AWSLambdaBasicExecutionRole`, `AmazonDynamoDBFullAccess`

---

### 4. Lambda Functions (Backend)

Set environment variables on each Lambda function:

| Variable | Description |
|---|---|
| `SLACK_TOKEN` | Slack bot token (`chat:write` scope) |
| `CHANNEL_NAME` | Default Slack channel for notifications |

Deploy all functions:

```bash
cd pubgstars-web/scripts
./buildAndUploadAll.sh
```

Or deploy a single function:

```bash
./buildAndUpload.sh games
```

The scripts cross-compile to `linux/amd64`, zip the binary, and call `aws lambda update-function-code`. Functions must already exist — create them once:

```bash
aws lambda create-function \
  --function-name games \
  --runtime provided.al2023 \
  --role arn:aws:iam::<account-id>:role/lambda-service-role \
  --handler bootstrap \
  --zip-file fileb://main.zip \
  --region eu-central-1 \
  --profile pg
```

Attach Cognito triggers:
- `canUserRegister` → pre sign-up trigger
- `userRegistered` → post-confirmation trigger

---

### 5. API Gateway

Create a REST API (`api.pubgstars.com`) with Lambda Proxy integration:

| Path | Methods | Lambda |
|---|---|---|
| `/games` | GET, POST, PUT, DELETE | `games` |
| `/games/password` | POST | `gamePassword` |
| `/games/users` | POST | `gameUsers` |
| `/games/register` | POST | `registerToGame` |
| `/games/unregister` | POST | `unregisterToGame` |
| `/games/history` | GET | `gamesHistory` |
| `/games/leaderboard` | GET | `gamesLeaderboard` |
| `/user` | GET | `user` |
| `/user/games` | GET | `userGames` |
| `/user/transactionlog` | GET | `transactionLog` |
| `/user/depositmoney` | POST | `depositMoney` |
| `/user/withdraw` | POST | `withdrawMoney` |
| `/user/sendmessage` | POST | `sendMessage` |
| `/admin/games` | GET, POST, PUT, DELETE | `adminGames` |
| `/admin/games/users` | GET | `adminGameUsers` |
| `/admin/users` | GET | `adminUsers` |
| `/admin/messages` | GET, POST | `adminMessages` |
| `/admin/completegame` | POST | `adminCompleteGame` |
| `/admin/addbudget` | POST | `adminAddBudget` |

Enable CORS on all resources. Deploy to a stage named `v1`.

The API Gateway passes events to Lambda in this format:

```json
{
  "body-json": { ... },
  "params": { "header": { "Authorization": "..." }, "querystring": {} },
  "context": { "http-method": "GET" }
}
```

---

### 6. Frontend (React SPA)

```bash
cd pubgstars-web-static
cp .env.example .env
npm install
npm run build
aws s3 sync build s3://pubgstars.com/ --profile pg
```

Configure the S3 bucket for static website hosting and point CloudFront to it.

---

## License

[MIT](LICENSE)
