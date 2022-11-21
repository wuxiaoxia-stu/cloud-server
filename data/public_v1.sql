/*
 Navicat Premium Data Transfer

 Source Server         : local
 Source Server Type    : PostgreSQL
 Source Server Version : 120003
 Source Host           : localhost:5432
 Source Catalog        : aiyun_cloud
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 120003
 File Encoding         : 65001

 Date: 24/05/2022 09:43:46
*/


-- ----------------------------
-- Sequence structure for lpm_auth_ukey_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."lpm_auth_ukey_id_seq";
CREATE SEQUENCE "public"."lpm_auth_ukey_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 32767
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for lpm_hospital_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."lpm_hospital_id_seq";
CREATE SEQUENCE "public"."lpm_hospital_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for lpm_region_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."lpm_region_id_seq";
CREATE SEQUENCE "public"."lpm_region_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for lpm_sys_admin_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."lpm_sys_admin_id_seq";
CREATE SEQUENCE "public"."lpm_sys_admin_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 32767
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for lpm_sys_role_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."lpm_sys_role_id_seq";
CREATE SEQUENCE "public"."lpm_sys_role_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 32767
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for lpm_user_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."lpm_user_id_seq";
CREATE SEQUENCE "public"."lpm_user_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Table structure for lpm_auth_ukey
-- ----------------------------
DROP TABLE IF EXISTS "public"."lpm_auth_ukey";
CREATE TABLE "public"."lpm_auth_ukey" (
  "id" int2 NOT NULL GENERATED BY DEFAULT AS IDENTITY (
INCREMENT 1
MINVALUE  1
MAXVALUE 32767
START 1
CACHE 1
),
  "name" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "create_at" timestamptz(6),
  "update_at" timestamptz(6),
  "status" int2 NOT NULL DEFAULT 1
)
;

-- ----------------------------
-- Records of lpm_auth_ukey
-- ----------------------------

-- ----------------------------
-- Table structure for lpm_hospital
-- ----------------------------
DROP TABLE IF EXISTS "public"."lpm_hospital";
CREATE TABLE "public"."lpm_hospital" (
  "id" int4 NOT NULL GENERATED BY DEFAULT AS IDENTITY (
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1
),
  "name" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "license_code" varchar(20) COLLATE "pg_catalog"."default",
  "create_at" timestamptz(6),
  "update_at" timestamptz(6),
  "status" int2 DEFAULT 1
)
;
COMMENT ON COLUMN "public"."lpm_hospital"."name" IS '医院名称';
COMMENT ON COLUMN "public"."lpm_hospital"."license_code" IS '营业执照';

-- ----------------------------
-- Records of lpm_hospital
-- ----------------------------

-- ----------------------------
-- Table structure for lpm_region
-- ----------------------------
DROP TABLE IF EXISTS "public"."lpm_region";
CREATE TABLE "public"."lpm_region" (
  "id" int4 NOT NULL GENERATED BY DEFAULT AS IDENTITY (
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1
),
  "pid" int4,
  "name" varchar(20) COLLATE "pg_catalog"."default",
  "status" int2 NOT NULL DEFAULT 1
)
;

-- ----------------------------
-- Records of lpm_region
-- ----------------------------

-- ----------------------------
-- Table structure for lpm_sys_admin
-- ----------------------------
DROP TABLE IF EXISTS "public"."lpm_sys_admin";
CREATE TABLE "public"."lpm_sys_admin" (
  "id" int2 NOT NULL GENERATED BY DEFAULT AS IDENTITY (
INCREMENT 1
MINVALUE  1
MAXVALUE 32767
START 1
CACHE 1
),
  "role_id" int2 NOT NULL DEFAULT 1,
  "username" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "password" varchar(50) COLLATE "pg_catalog"."default",
  "salt" varchar(5) COLLATE "pg_catalog"."default",
  "phone" varchar(20) COLLATE "pg_catalog"."default",
  "email" varchar(30) COLLATE "pg_catalog"."default",
  "avatar" varchar(150) COLLATE "pg_catalog"."default",
  "create_at" timestamptz(6),
  "update_at" timestamptz(6),
  "status" int2 DEFAULT 1
)
;

-- ----------------------------
-- Records of lpm_sys_admin
-- ----------------------------
INSERT INTO "public"."lpm_sys_admin" VALUES (2, 1, 'admin', '315193380d2fedffa677b7fe236fadb5', 'sinb', NULL, '778774780@qq.com', NULL, '2022-05-23 15:58:23+08', '2022-05-23 15:58:23+08', 1);

-- ----------------------------
-- Table structure for lpm_sys_role
-- ----------------------------
DROP TABLE IF EXISTS "public"."lpm_sys_role";
CREATE TABLE "public"."lpm_sys_role" (
  "id" int2 NOT NULL GENERATED BY DEFAULT AS IDENTITY (
INCREMENT 1
MINVALUE  1
MAXVALUE 32767
START 1
CACHE 1
),
  "name" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "create_at" timestamptz(6),
  "update_at" timestamptz(6),
  "status" int2 NOT NULL DEFAULT 1
)
;

-- ----------------------------
-- Records of lpm_sys_role
-- ----------------------------

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."lpm_auth_ukey_id_seq"
OWNED BY "public"."lpm_auth_ukey"."id";
SELECT setval('"public"."lpm_auth_ukey_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."lpm_hospital_id_seq"
OWNED BY "public"."lpm_hospital"."id";
SELECT setval('"public"."lpm_hospital_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."lpm_region_id_seq"
OWNED BY "public"."lpm_region"."id";
SELECT setval('"public"."lpm_region_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."lpm_sys_admin_id_seq"
OWNED BY "public"."lpm_sys_admin"."id";
SELECT setval('"public"."lpm_sys_admin_id_seq"', 2, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."lpm_sys_role_id_seq"
OWNED BY "public"."lpm_sys_role"."id";
SELECT setval('"public"."lpm_sys_role_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."lpm_user_id_seq"', 3, true);

-- ----------------------------
-- Auto increment value for lpm_auth_ukey
-- ----------------------------
SELECT setval('"public"."lpm_auth_ukey_id_seq"', 1, false);

-- ----------------------------
-- Auto increment value for lpm_hospital
-- ----------------------------
SELECT setval('"public"."lpm_hospital_id_seq"', 1, false);

-- ----------------------------
-- Primary Key structure for table lpm_hospital
-- ----------------------------
ALTER TABLE "public"."lpm_hospital" ADD CONSTRAINT "lpm_hospital_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Auto increment value for lpm_region
-- ----------------------------
SELECT setval('"public"."lpm_region_id_seq"', 1, false);

-- ----------------------------
-- Primary Key structure for table lpm_region
-- ----------------------------
ALTER TABLE "public"."lpm_region" ADD CONSTRAINT "lpm_region_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Auto increment value for lpm_sys_admin
-- ----------------------------
SELECT setval('"public"."lpm_sys_admin_id_seq"', 2, true);

-- ----------------------------
-- Auto increment value for lpm_sys_role
-- ----------------------------
SELECT setval('"public"."lpm_sys_role_id_seq"', 1, false);
