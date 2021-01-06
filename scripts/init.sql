SET client_encoding = 'UTF8';
SELECT pg_catalog.set_config('search_path', '', false);

DROP TABLE IF EXISTS public.forum CASCADE;
DROP TABLE IF EXISTS public.users CASCADE;
DROP TABLE IF EXISTS public.thread CASCADE;
DROP TABLE IF EXISTS public.post CASCADE;
DROP TABLE IF EXISTS public.vote CASCADE;
DROP TABLE IF EXISTS public.forum_users CASCADE;

CREATE TABLE public.forum (
    ID serial NOT NULL PRIMARY KEY,
    slug text NOT NULL,
    threads integer DEFAULT 0 NOT NULL,
    posts integer DEFAULT 0 NOT NULL,
    title text NOT NULL,
    author text NOT NULL
);
ALTER TABLE public.forum OWNER TO forum_user;

CREATE TABLE public.users (
    ID serial NOT NULL PRIMARY KEY,
    nick_name text NOT NULL,
    email text NOT NULL,
    full_name text NOT NULL,
    about text
);
ALTER TABLE public.users OWNER TO forum_user;

CREATE TABLE public.forum_users (
    forum text NOT NULL,
    user_nick text NOT NULL
);
ALTER TABLE IF EXISTS public.forum_users ADD CONSTRAINT uniq PRIMARY KEY (forum, user_nick);
ALTER TABLE public.forum_users OWNER TO forum_user;

CREATE TABLE public.thread
(
    ID      serial                                 NOT NULL PRIMARY KEY,
    author  text                                   NOT NULL,
--    created text NOT NULL,
    created timestamp with time zone DEFAULT now() NOT NULL,
    forum   text                                   NOT NULL,
    message text                                   NOT NULL,
    slug    text,
    title   text                                   NOT NULL,
    votes   integer                  DEFAULT 0
);
ALTER TABLE public.thread OWNER TO forum_user;


/*CREATE TABLE public.post
(
    ID        serial                NOT NULL PRIMARY KEY,
    author    text                  NOT NULL,
    created text NOT NULL,
 --   created TIMESTAMP WITH TIME ZONE,
    forum     text                  NOT NULL,
    is_edited boolean DEFAULT false NOT NULL,
    message   text                  NOT NULL,
    parent    integer DEFAULT 0     NOT NULL,
    thread    integer               NOT NULL,
    path    INTEGER[] DEFAULT '{0}':: INTEGER [] NOT NULL
);*/

CREATE TABLE public.post (
                             id integer NOT NULL PRIMARY KEY,
                             author text NOT NULL,
                             created text NOT NULL,
                             --created TIMESTAMP WITH TIME ZONE,

                             forum text NOT NULL,
                             is_edited boolean DEFAULT false NOT NULL,
                             message text NOT NULL,
                             parent integer DEFAULT 0 NOT NULL,
                             thread integer NOT NULL,
                             path INTEGER[] DEFAULT '{0}'::INTEGER[] NOT NULL
);

CREATE SEQUENCE public.post_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.post_id_seq OWNED BY public.post.id;
ALTER TABLE ONLY public.post ALTER COLUMN id SET DEFAULT nextval('public.post_id_seq'::regclass);
SELECT pg_catalog.setval('public.post_id_seq', 1, false);
ALTER TABLE public.post OWNER TO forum_user;


CREATE TABLE public.vote (
    user_nick text NOT NULL,
    voice integer NOT NULL,
    thread_id integer NOT NULL
);
ALTER TABLE public.vote OWNER TO forum_user;

--
-- Name: forum_slug_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX forum_slug_uindex ON public.forum USING btree (lower(slug));


--
-- Name: post_author_forum_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX post_author_forum_index ON public.post USING btree (lower(author), lower(forum));


--
-- Name: post_forum_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX post_forum_index ON public.post USING btree (lower(forum));


--
-- Name: post_parent_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX post_parent_index ON public.post USING btree (parent);


--
-- Name: post_path_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX post_path_index ON public.post USING gin (path);


--
-- Name: post_thread_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX post_thread_index ON public.post USING btree (thread);


--
-- Name: thread_forum_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX thread_forum_index ON public.thread USING btree (lower(forum));


--
-- Name: thread_id_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX thread_id_uindex ON public.thread USING btree (id);


--
-- Name: thread_slug_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX thread_slug_index ON public.thread USING btree (lower(slug));


--
-- Name: user_email_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX user_email_uindex ON public.users USING btree (lower(email));


--
-- Name: user_nick_name_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX user_nick_name_uindex ON public.users USING btree (lower(nick_name));