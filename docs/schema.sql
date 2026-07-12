-- Add new schema named "public"
CREATE SCHEMA IF NOT EXISTS "public";
-- Set comment to schema: "public"
COMMENT ON SCHEMA "public" IS 'standard public schema';
-- Create "history" table
CREATE TABLE "public"."history" (
  "id" serial NOT NULL,
  "source" character varying(255) NOT NULL,
  "destination" character varying(255) NOT NULL,
  "original" character varying(255) NOT NULL,
  "translation" character varying(255) NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "users" table
CREATE TABLE "public"."users" (
  "id" character varying(255) NOT NULL,
  "username" character varying(255) NOT NULL,
  "email" character varying(255) NOT NULL,
  "password_hash" character varying(255) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  UNIQUE ("username"),
  UNIQUE ("email")
);
-- Create "tasks" table
CREATE TABLE "public"."tasks" (
  "id" character varying(255) NOT NULL,
  "user_id" character varying(255) NOT NULL,
  "title" character varying(255) NOT NULL,
  "description" character varying(1000) NOT NULL DEFAULT '',
  "status" character varying(50) NOT NULL DEFAULT 'todo',
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  FOREIGN KEY ("user_id") REFERENCES users ("id") ON DELETE CASCADE
);
-- Create "schema_migrations" table
CREATE TABLE "public"."schema_migrations" (
  "version" bigint NOT NULL,
  "dirty" boolean NOT NULL,
  PRIMARY KEY ("version")
);
