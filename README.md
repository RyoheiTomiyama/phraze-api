# Phraze API

## 📂 ディレクトリ構成

```tree
# APIアプリケーション
app/
├── domain/
│   └── domain.go
├── application/
│   └── usecase/
│       ├── auth.go
│       └── deck.go
├── infra/
│   └── db/
│       ├── client.go
│       └── get_user.go
├── router/
│   ├── graph/
│   │   └── resolver/
│   └── schema/
└── util/
    └── logger/
        └── logger.go


# マイグレーション・シード管理
atlas/
├── schema.sql
└── seeds/
    └── development/
        └── 20240628144146_users.sql
```

なんちゃってクリーンアーキテクチャを意識

層が厚くなり開発スピードが落ちないようにインターフェイス、レポジトリ層を省略。
ただし、DI や関心の分離は意識して作成している。

他の層と依存する場合は、interface を参照すること。

router ← usecase ← infra

## 🔧 Development

### 環境変数

`.env` で管理

### Google サービスアカウントのクレデンシャル情報

クレデンシャル情報が必要になります。  
管理者に問い合わせてください。

### 開発環境の起動

```bash
make up
```

- [マイグレーション・シード管理について](./atlas/)
- [GraphQL の管理について](./app/infra/graph/)

## Test

### 環境変数

`.env.test` で管理

### 環境設定

#### VSCode の Run Test に環境設定を読み込ませる

settings.json に以下を追加する

```json
  "go.testEnvVars": {
    "TZ": "UTC"
  },
  "go.testEnvFile": "${workspaceFolder}/.env.test"
```

#### VSCode Golangci integration

VSCode で linter を有効にする設定

settings.json に以下を追加する

```json
  "go.lintTool": "golangci-lint",
  "go.lintFlags": ["--fast"]
```

## Production

### 本番環境用の Docker image をビルドする

API

```sh
docker build --platform linux/amd64 -f docker/Dockerfile.production -t phraze-app-prd .
```

Migration

```sh
docker build --platform linux/amd64 -f docker/Dockerfile.migration -t phraze-migration-prd .
```

本番サーバーで動かすために platform の指定が必要
