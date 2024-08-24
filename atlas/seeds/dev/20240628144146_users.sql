insert into
    "public"."users" ("id", "name")
values
    (1, 'Admin') on conflict (id) do
update
set
    name = excluded.name;