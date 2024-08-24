insert into
    "public"."decks" ("id", "user_id", "name")
values
    (1, 1, 'sample') on conflict (id) do
update
set
    "user_id" = excluded."user_id",
    "name" = excluded."name"