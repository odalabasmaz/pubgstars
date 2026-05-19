# PubgStars

A full-stack platform for organizing paid PUBG tournament sessions. Players deposit balance, register for scheduled game rooms, receive room credentials shortly before kick-off, and winners collect prize money — all handled automatically.

Built on a fully serverless AWS stack. The UI targets Turkish-speaking players; backend logic and code are in English.

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

**Language:** Go  
**Runtime:** AWS Lambda (one binary per endpoint)  
**Database:** AWS DynamoDB  
**Auth:** AWS Cognito — JWT tokens validated in each handler  
**Notifications:** Slack API (new registrations posted to a channel)  
**Deployment:** Cross-compiled to `linux/amd64`, zipped, and uploaded with `aws lambda update-function-code`

### Lambda functions

| Function | Method(s) | Description |
|---|---|---|
| `games` | GET / POST / PUT / DELETE | List active games; create or update a session |
| `gamePassword` | POST | Return room credentials (only within ±1 h of game time) |
| `gameUsers` | POST | List users registered for a game |
| `registerToGame` | POST | Register the caller; deducts entry fee atomically |
| `unregisterToGame` | POST | Unregister (refunds as bonus if before cutoff) |
| `gamesHistory` | GET | Completed game history |
| `gamesLeaderboard` | GET | All-time winner leaderboard |
| `user` | GET | Current user profile |
| `userGames` | GET | Games the caller is registered for |
| `canUserRegister` | GET | Username availability check |
| `depositMoney` | POST | Record a deposit request |
| `withdrawMoney` | POST | Request a withdrawal (requires secret answer) |
| `transactionLog` | GET | Full balance + game history for the caller |
| `sendMessage` | POST | Submit a support message |
| `adminGames` | GET / POST | Admin: list all games, create new sessions |
| `adminCompleteGame` | POST | Admin: mark complete, credit prize money to winners |
| `adminGameUsers` | GET | Admin: view registrations per game |
| `adminUsers` | GET | Admin: list all users |
| `adminMessages` | GET | Admin: view support messages |
| `adminAddBudget` | POST | Admin: manually credit balance / bonus to a user |
| `userRegistered` | — | Cognito post-confirmation trigger: creates user record |

### DynamoDB tables

| Table | Key | Purpose |
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
├── cmd/                  # One main.go per Lambda function
│   ├── games/
│   ├── registerToGame/
│   ├── adminCompleteGame/
│   └── ...
├── internal/
│   ├── AwsUtils.go       # DynamoDB client, JWT parsing, time helpers
│   ├── DataService.go    # All DynamoDB read/write operations
│   ├── GameUtils.go      # Time-window logic for passwords / cancellation
│   ├── SlackService.go   # Slack notifications
│   ├── ModelUtils.go     # ID generation, misc utilities
│   ├── TransactionLogUtils.go
│   └── LoggerService.go
├── model/
│   ├── Model.go          # Game, User, TransactionLog, Message structs
│   └── tables/Tables.go  # DynamoDB table name constants
├── scripts/
│   ├── buildAndUpload.sh     # Build + deploy a single Lambda
│   └── buildAndUploadAll.sh  # Build + deploy all Lambdas
└── pkg/                  # Shared logger / printer utilities
```

### Local development

```bash
# Build a single function
cd scripts
./buildAndUpload.sh games

# Build and deploy all functions
./buildAndUploadAll.sh

# Run with local DynamoDB (requires AWS CLI profile named 'pg')
isForTest=true go run cmd/games/main.go
```

Required environment variables:

```
SLACK_TOKEN=        # Slack bot token
CHANNEL_NAME=       # Target Slack channel (e.g. #customer-requests)
```

See `pubgstars-web/.env.example` for the full list.

---

## Frontend — `pubgstars-web-static`

**Framework:** React (Create React App)  
**Auth:** AWS Amplify + Cognito  
**HTTP:** Axios — all calls go to `https://api.pubgstars.com/v1`  
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
| `/profile` | Profile | Account details and password change |
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
npm start
```

Required environment variables (Create React App — must be prefixed `REACT_APP_`):

```
REACT_APP_STAGE=dev
REACT_APP_COGNITO_REGION=eu-central-1
REACT_APP_COGNITO_USER_POOL_ID=
REACT_APP_COGNITO_APP_CLIENT_ID=
REACT_APP_COGNITO_IDENTITY_POOL_ID=
```

See `pubgstars-web-static/.env.example`.

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

Cognito User Pool ──► post-confirmation trigger ──► userRegistered Lambda
```

All backend compute is Lambda; there are no long-running servers.

---

## License

[MIT](LICENSE)
