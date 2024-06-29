# GraphQL 管理

[gqlgen](https://gqlgen.com/)を利用して、GraphQL のコード生成を行っている。

## schema の追加

`schema/`に GraphQL のスキーマを定義する

user のスキーマを作る場合は、`user.graphqls` を作成

```gql
type User {
  id: ID!
  name: String!
}

extend type Query {
  users: [User!]!
  user: User!
}
```

gqlgen でファイルを生成する

```bash
make gqlgen
```
