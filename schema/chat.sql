-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2023-09-19T06:23:58.823Z

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "avatar" varchar NOT NULL,
  "nickname" varchar NOT NULL,
  "role" varchar NOT NULL,
  "room_id" bigint NOT NULL,
  "deleted" boolean NOT NULL DEFAULT false,
  "create_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "refresh_token" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "expire_at" timestamptz NOT NULL,
  "create_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "friendships" (
  "user_id" bigint,
  "friend_id" bigint,
  "room_id" bigint NOT NULL,
  "status" varchar NOT NULL,
  "create_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY ("user_id", "friend_id")
);

CREATE TABLE "rooms" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "cover" varchar NOT NULL,
  "category" varchar NOT NULL,
  "create_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "room_members" (
  "room_id" bigint NOT NULL,
  "member_id" bigint NOT NULL,
  "rank" varchar NOT NULL,
  "join_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY ("room_id", "member_id")
);

CREATE TABLE "messages" (
  "id" bigserial PRIMARY KEY,
  "room_id" bigint NOT NULL,
  "sender_id" bigint NOT NULL,
  "content" varchar NOT NULL,
  "kind" varchar NOT NULL,
  "send_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "users" ADD FOREIGN KEY ("room_id") REFERENCES "rooms" ("id");

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "friendships" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "friendships" ADD FOREIGN KEY ("friend_id") REFERENCES "users" ("id");

ALTER TABLE "friendships" ADD FOREIGN KEY ("room_id") REFERENCES "rooms" ("id");

ALTER TABLE "room_members" ADD FOREIGN KEY ("room_id") REFERENCES "rooms" ("id");

ALTER TABLE "room_members" ADD FOREIGN KEY ("member_id") REFERENCES "users" ("id");

ALTER TABLE "messages" ADD FOREIGN KEY ("room_id") REFERENCES "rooms" ("id");

ALTER TABLE "messages" ADD FOREIGN KEY ("sender_id") REFERENCES "users" ("id");
