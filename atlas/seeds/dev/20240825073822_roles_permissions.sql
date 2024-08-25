insert into
    "public"."permissions" ("id", "key", "name")
values
    (1, 'unlimitedCardCreation', '無制限のカード作成'),
    (2, 'unlimitedAIAnswerGeneration', '無制限のAI解答生成'),
    (3, 'limitedAIAnswer100rpd', '100リクエスト/day AI解答生成') on conflict (id) do
update
set
    "key" = excluded."key",
    "name" = excluded."name";

insert into
    "public"."roles" ("id", "key", "name")
values
    (1, 'planBasic', 'ベーシックプラン'),
    (2, 'planPro', 'プロプラン') on conflict (id) do
update
set
    "key" = excluded."key",
    "name" = excluded."name";

insert into
    "public"."roles_permissions" ("role_id", "permission_id")
values
    -- planBasic
    (1, 1),
    (1, 3),
    -- planPro
    (2, 1),
    (2, 2) on conflict (role_id, permission_id) do
update
set
    "role_id" = excluded."role_id",
    "permission_id" = excluded."permission_id";