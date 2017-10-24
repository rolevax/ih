--
-- PostgreSQL database dump
--

-- Dumped from database version 10.0
-- Dumped by pg_dump version 10.0

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: postgres; Type: COMMENT; Schema: -; Owner: postgres
--

COMMENT ON DATABASE postgres IS 'default administrative connection database';


--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

--
-- Name: play(integer[]); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION play(ranks integer[]) RETURNS integer
    LANGUAGE plpgsql
    AS $$
BEGIN
  RETURN ranks[1]+ranks[2]+ranks[3]+ranks[4];
END;
$$;


ALTER FUNCTION public.play(ranks integer[]) OWNER TO postgres;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: replay_of_user; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE replay_of_user (
    user_id bigint NOT NULL,
    replay_ids bigint[] NOT NULL
);


ALTER TABLE replay_of_user OWNER TO postgres;

--
-- Name: replays; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE replays (
    replay_id bigint NOT NULL,
    content text NOT NULL
);


ALTER TABLE replays OWNER TO postgres;

--
-- Name: replays_replay_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE replays_replay_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE replays_replay_id_seq OWNER TO postgres;

--
-- Name: replays_replay_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE replays_replay_id_seq OWNED BY replays.replay_id;


--
-- Name: tasks; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE tasks (
    task_id integer NOT NULL,
    title character varying(64) NOT NULL,
    content text NOT NULL,
    state integer DEFAULT 0 NOT NULL,
    assignee_id bigint,
    c_point integer DEFAULT 0 NOT NULL
);


ALTER TABLE tasks OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE users (
    user_id bigint NOT NULL,
    username character varying(64) NOT NULL,
    password character(44) NOT NULL,
    c_point integer DEFAULT 0 NOT NULL
);


ALTER TABLE users OWNER TO postgres;

--
-- Name: users_user_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE users_user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE users_user_id_seq OWNER TO postgres;

--
-- Name: users_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE users_user_id_seq OWNED BY users.user_id;


--
-- Name: replays replay_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY replays ALTER COLUMN replay_id SET DEFAULT nextval('replays_replay_id_seq'::regclass);


--
-- Name: users user_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY users ALTER COLUMN user_id SET DEFAULT nextval('users_user_id_seq'::regclass);


--
-- Data for Name: replay_of_user; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY replay_of_user (user_id, replay_ids) FROM stdin;
\.


--
-- Data for Name: replays; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY replays (replay_id, content) FROM stdin;
\.


--
-- Data for Name: tasks; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY tasks (task_id, title, content, state, assignee_id, c_point) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY users (user_id, username, password, c_point) FROM stdin;
501	ⓝ喵打	iMqwzWh78q4vOGkh7ALsA6ohvh25OQ/VMDNFICkTarc=	0
502	ⓝ打喵	iMqwzWh78q4vOGkh7ALsA6ohvh25OQ/VMDNFICkTarc=	0
1000	rolevax	Ddddddddddddddddddddddddddddddddddddddddddd=	0
\.


--
-- Name: replays_replay_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('replays_replay_id_seq', 1, true);


--
-- Name: users_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('users_user_id_seq', 1976, true);


--
-- Name: replays idx_17506_primary; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY replays
    ADD CONSTRAINT idx_17506_primary PRIMARY KEY (replay_id);


--
-- Name: users idx_17515_primary; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY users
    ADD CONSTRAINT idx_17515_primary PRIMARY KEY (user_id);


--
-- Name: replay_of_user replay_of_user_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY replay_of_user
    ADD CONSTRAINT replay_of_user_pkey PRIMARY KEY (user_id);


--
-- Name: tasks tasks_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY tasks
    ADD CONSTRAINT tasks_pk PRIMARY KEY (task_id);


--
-- Name: idx_17515_username; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_17515_username ON users USING btree (username);


--
-- Name: replay_of_user replay_of_user_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY replay_of_user
    ADD CONSTRAINT replay_of_user_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(user_id);


--
-- Name: tasks tasks_assignee_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY tasks
    ADD CONSTRAINT tasks_assignee_fkey FOREIGN KEY (assignee_id) REFERENCES users(user_id);


--
-- PostgreSQL database dump complete
--

