CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "queue" (
	"id" uuid DEFAULT uuid_generate_v4(),
    "name" text not null,
    "priority" INT default 1,
    "status" INT default 1,
    "service" INT default 1,
    "request_payload" JSONB NULL,
    "response_payload" JSONB NULL,
    "execute_at" timestamp default null, 
    "create_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    "update_at" timestamp default NULL,
    CONSTRAINT queue_id PRIMARY KEY (id)
);