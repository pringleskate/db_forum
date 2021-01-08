SET client_encoding = 'UTF8';

DROP TABLE IF EXISTS forum CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS thread CASCADE;
DROP TABLE IF EXISTS post CASCADE;
DROP TABLE IF EXISTS vote CASCADE;
DROP TABLE IF EXISTS forum_users CASCADE;

CREATE TABLE forum
(
    ID      serial            NOT NULL PRIMARY KEY,
    slug    text              NOT NULL,
    threads integer DEFAULT 0 NOT NULL,
    posts   integer DEFAULT 0 NOT NULL,
    title   text              NOT NULL,
    author  text              NOT NULL
);

CREATE TABLE users
(
    ID        serial NOT NULL PRIMARY KEY,
    nick_name text   NOT NULL,
    email     text   NOT NULL,
    full_name text   NOT NULL,
    about     text
);

CREATE TABLE forum_users
(
    forum     text NOT NULL,
    user_nick text NOT NULL
);
ALTER TABLE IF EXISTS forum_users ADD CONSTRAINT uniq PRIMARY KEY (forum, user_nick);

CREATE TABLE thread
(
    ID      serial                                 NOT NULL PRIMARY KEY,
    author  text                                   NOT NULL,
    created timestamp with time zone DEFAULT now() NOT NULL,
    forum   text                                   NOT NULL,
    message text                                   NOT NULL,
    slug    text,
    title   text                                   NOT NULL,
    votes   integer                  DEFAULT 0
);

CREATE TABLE post
(
    id        integer                            NOT NULL PRIMARY KEY,
    author    text                               NOT NULL,
    created   text                               NOT NULL,
    forum     text                               NOT NULL,
    is_edited boolean   DEFAULT false            NOT NULL,
    message   text                               NOT NULL,
    parent    integer   DEFAULT 0                NOT NULL,
    thread    integer                            NOT NULL,
    path      INTEGER[] DEFAULT '{0}'::INTEGER[] NOT NULL
);

CREATE SEQUENCE post_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE post_id_seq OWNED BY post.id;
ALTER TABLE ONLY post ALTER COLUMN id SET DEFAULT nextval('post_id_seq'::regclass);
SELECT pg_catalog.setval('post_id_seq', 1, false);

CREATE TABLE vote (
                             user_nick text NOT NULL,
                             voice integer NOT NULL,
                             thread_id integer NOT NULL
);

CREATE UNIQUE INDEX forum_slug_uindex ON forum USING btree (lower(slug));

CREATE INDEX post_author_forum_index ON post USING btree (lower(author), lower(forum));

CREATE INDEX post_forum_index ON post USING btree (lower(forum));

CREATE INDEX post_parent_index ON post USING btree (parent);

CREATE INDEX post_path_index ON post USING gin (path);

CREATE INDEX post_thread_index ON post USING btree (thread);

CREATE INDEX thread_forum_index ON thread USING btree (lower(forum));

CREATE UNIQUE INDEX thread_id_uindex ON thread USING btree (id);

CREATE INDEX thread_slug_index ON thread USING btree (lower(slug));

CREATE UNIQUE INDEX user_email_uindex ON users USING btree (lower(email));

CREATE UNIQUE INDEX user_nick_name_uindex ON users USING btree (lower(nick_name));
