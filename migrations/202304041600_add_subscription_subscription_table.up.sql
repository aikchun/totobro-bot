BEGIN;
CREATE TABLE IF NOT EXISTS subscription_subscription (
	"id" bigserial PRIMARY KEY,
	"chat_id" bigint NOT NULL UNIQUE,
	"threshold" bigint NOT NULL,
	"is_active" BOOLEAN NOT NULL DEFAULT '1',
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL
);
COMMIT;