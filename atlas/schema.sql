-- Add new schema named "public"
CREATE SCHEMA IF NOT EXISTS "public";

-- Set comment to schema: "public"
COMMENT ON SCHEMA "public" IS 'standard public schema';

-- デッキ
CREATE TABLE
    "public"."decks" (
        "id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY,
        "user_id" TEXT NOT NULL,
        "name" VARCHAR(100) NOT NULL,
        "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY ("id")
    );

-- カード
CREATE TABLE
    "public"."cards" (
        "id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY,
        "deck_id" bigint NOT NULL REFERENCES decks (id) ON DELETE RESTRICT,
        "question" TEXT NOT NULL,
        "answer" TEXT NOT NULL,
        "ai_answer" TEXT NOT NULL DEFAULT '',
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY ("id")
    );

-- トリガー関数の作成
CREATE
OR REPLACE FUNCTION update_timestamp () RETURNS TRIGGER AS '
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
' LANGUAGE plpgsql;

-- トリガーの作成
CREATE TRIGGER update_decks_updated_at BEFORE
UPDATE ON decks FOR EACH ROW EXECUTE FUNCTION update_timestamp ();

CREATE TRIGGER update_cards_updated_at BEFORE
UPDATE ON cards FOR EACH ROW EXECUTE FUNCTION update_timestamp ();

--　カード評価結果
CREATE TABLE
    "public"."card_reviews" (
        "id" BIGINT NOT NULL GENERATED BY DEFAULT AS IDENTITY,
        "card_id" BIGINT NOT NULL REFERENCES cards (id) ON DELETE RESTRICT,
        "user_id" TEXT NOT NULL,
        "reviewed_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "grade" INT NOT NULL CHECK (grade BETWEEN 1 AND 5),
        PRIMARY KEY (id)
    );

CREATE INDEX idx_card_reviews_card_id ON "card_reviews" ("card_id");

CREATE INDEX idx_card_reviews_user_id_reviewed_at ON "card_reviews" ("user_id", "reviewed_at");

-- 次のカード学習情報
CREATE TABLE
    "public"."card_schedules" (
        "id" BIGINT NOT NULL GENERATED BY DEFAULT AS IDENTITY,
        "card_id" BIGINT NOT NULL UNIQUE REFERENCES cards (id) ON DELETE RESTRICT,
        "schedule_at" TIMESTAMP NOT NULL,
        "interval" INT NOT NULL, -- 復習間隔（MINUTE）
        "efactor" FLOAT NOT NULL, -- イーファクター
        PRIMARY KEY (card_id)
    );

CREATE INDEX idx_card_schedules_schedule_at ON "card_schedules" ("schedule_at");

-- 権限
CREATE TABLE
    "public"."permissions" (
        "id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY,
        "key" VARCHAR(256) NOT NULL UNIQUE,
        "name" TEXT NOT NULL,
        PRIMARY KEY ("id")
    );

-- ロール
CREATE TABLE
    "public"."roles" (
        "id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY,
        "key" VARCHAR(256) NOT NULL UNIQUE,
        "name" TEXT NOT NULL,
        PRIMARY KEY ("id")
    );

CREATE TABLE
    "public"."roles_permissions" (
        "role_id" bigint NOT NULL REFERENCES roles (id) ON DELETE RESTRICT,
        "permission_id" bigint NOT NULL REFERENCES permissions (id) ON DELETE RESTRICT,
        PRIMARY KEY ("role_id", "permission_id")
    );

-- User
CREATE TABLE
    "public"."users" (
        "id" TEXT NOT NULL UNIQUE,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY ("id")
    );

CREATE TRIGGER update_users_updated_at BEFORE
UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_timestamp ();

-- User Role
CREATE TABLE
    "public"."users_roles" (
        "user_id" TEXT NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
        "role_id" bigint NOT NULL REFERENCES roles (id) ON DELETE RESTRICT,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY ("user_id", "role_id")
    );

CREATE TRIGGER update_users_roles_updated_at BEFORE
UPDATE ON users_roles FOR EACH ROW EXECUTE FUNCTION update_timestamp ();