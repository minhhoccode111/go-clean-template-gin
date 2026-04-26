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
-- Create "schema_migrations" table
CREATE TABLE "public"."schema_migrations" (
  "version" bigint NOT NULL,
  "dirty" boolean NOT NULL,
  PRIMARY KEY ("version")
);
