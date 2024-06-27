schema "public" {
  comment = "standard public schema"

}

table "users" {
  schema = schema.public
  column "id" {
    type = bigint
    null = false
    identity {
      generated = BY_DEFAULT
    }
  }
  column "name" {
    type = varchar(255)
    null = false
  }
  primary_key {
    columns = [
      column.id
    ]
  }
}