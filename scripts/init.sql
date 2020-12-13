DROP TABLE IF EXISTS forum CASCADE;
DROP TABLE IF EXISTS user CASCADE;
DROP TABLE IF EXISTS thread CASCADE;
DROP TABLE IF EXISTS post CASCADE;
DROP TABLE IF EXISTS vote CASCADE;
DROP TABLE IF EXISTS forum_users CASCADE;

CREATE TABLE forum (
    ID serial NOT NULL,
    slug text NOT NULL,
    threads integer DEFAULT 0 NOT NULL,
    posts integer DEFAULT 0 NOT NULL,
    title text NOT NULL,
    user text NOT NULL
);

CREATE TABLE user (
    ID serial NOT NULL,
    nick_name text NOT NULL,
    email text NOT NULL,
    full_name text NOT NULL,
    about text
);

CREATE TABLE forum_user (
    forum_ID integer NOT NULL,
    user_ID integer NOT NULL
);

CREATE TABLE thread
(
    ID      serial                                 NOT NULL,
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
    ID        serial                NOT NULL,
    author    text                  NOT NULL,
    created   text                  NOT NULL,
    forum     text                  NOT NULL,
    is_edited boolean DEFAULT false NOT NULL,
    message   text                  NOT NULL,
    parent    integer DEFAULT 0     NOT NULL,
    thread    integer               NOT NULL,
    path      integer[]:: integer [] NOT NULL
);