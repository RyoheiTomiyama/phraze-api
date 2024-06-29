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
│   ├── router
│   ├── graph/
│   │   └── repository.go
│   └── db/
│       ├── client.go
│       └── get_user.go
└── util/
    └── logger/
        └── logger.go

# マイグレーション管理
atlas/
└── schema.sql

# シード管理
seeds/
└── dev/
    └── 20240628144146_users.sql
```

なんちゃってクリーンアーキテクチャを意識

層が厚くなり開発スピードが落ちないようにインターフェイス、レポジトリ層を省略。
ただし、DI や関心の分離は意識して作成している。

他の層と依存する場合は、interface を参照すること。

## 🔧 Development

環境変数は、`.env` で管理

開発環境の起動

```bash
make up
```

マイグレーション・シード管理に関しては、`atlas/` の README を参照してください。
