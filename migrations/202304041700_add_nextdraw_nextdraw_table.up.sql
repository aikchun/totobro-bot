BEGIN;
CREATE TABLE IF NOT EXISTS nextdraw_nextdraw (
	"id" bigserial PRIMARY KEY,
    "date_string" VARCHAR(64) UNIQUE NOT NULL,
    "prize" VARCHAR(64) NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL
);
COMMIT;