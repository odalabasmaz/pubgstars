# pubgstars-web-static

React frontend for the PubgStars platform. Deployed as a static SPA to S3 / CloudFront.

For full architecture and infrastructure documentation see the [root README](../README.md).

---

## Requirements

| Tool | Version |
|---|---|
| Node.js | 18–20 (use `.nvmrc`: `nvm use`) |
| npm | 9+ |

> `react-scripts 5` does not support Node 21+.

---

## Local development

```bash
cp .env.example .env   # fill in Cognito values
npm install
npm start              # starts dev server at http://localhost:3000
```

---

## Environment variables

All variables must be prefixed with `REACT_APP_` (Create React App requirement):

| Variable | Description |
|---|---|
| `REACT_APP_STAGE` | Environment name (`dev` / `prod`) |
| `REACT_APP_COGNITO_REGION` | AWS region of the Cognito User Pool |
| `REACT_APP_COGNITO_USER_POOL_ID` | Cognito User Pool ID |
| `REACT_APP_COGNITO_APP_CLIENT_ID` | Cognito App Client ID |
| `REACT_APP_COGNITO_IDENTITY_POOL_ID` | Cognito Identity Pool ID |

---

## Available scripts

| Command | Description |
|---|---|
| `npm start` | Start development server |
| `npm run build` | Production build into `build/` |
| `npm test` | Run React test suite |

---

## Deploying

```bash
npm run build
aws s3 sync build s3://pubgstars.com/ --profile pg
```

Requires an S3 bucket configured for static website hosting with a CloudFront distribution in front of it.
