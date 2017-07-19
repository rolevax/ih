--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.3
-- Dumped by pg_dump version 9.6.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

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
-- Name: play(integer[]); Type: FUNCTION; Schema: public; Owner: mako
--

CREATE FUNCTION play(ranks integer[]) RETURNS integer
    LANGUAGE plpgsql
    AS $$
BEGIN
    RETURN ranks[1]+ranks[2]+ranks[3]+ranks[4];
END;
$$;


ALTER FUNCTION public.play(ranks integer[]) OWNER TO mako;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: girls; Type: TABLE; Schema: public; Owner: mako
--

CREATE TABLE girls (
    girl_id integer NOT NULL,
    level smallint DEFAULT 0 NOT NULL,
    pt smallint DEFAULT 0 NOT NULL,
    rating double precision DEFAULT '1500'::double precision NOT NULL
);


ALTER TABLE girls OWNER TO mako;

--
-- Name: replay_of_user; Type: TABLE; Schema: public; Owner: mako
--

CREATE TABLE replay_of_user (
    user_id bigint NOT NULL,
    replay_ids bigint[] NOT NULL
);


ALTER TABLE replay_of_user OWNER TO mako;

--
-- Name: replays; Type: TABLE; Schema: public; Owner: mako
--

CREATE TABLE replays (
    replay_id bigint NOT NULL,
    content text NOT NULL
);


ALTER TABLE replays OWNER TO mako;

--
-- Name: replays_replay_id_seq; Type: SEQUENCE; Schema: public; Owner: mako
--

CREATE SEQUENCE replays_replay_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE replays_replay_id_seq OWNER TO mako;

--
-- Name: replays_replay_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mako
--

ALTER SEQUENCE replays_replay_id_seq OWNED BY replays.replay_id;


--
-- Name: user_girl; Type: TABLE; Schema: public; Owner: mako
--

CREATE TABLE user_girl (
    user_id bigint NOT NULL,
    girl_id integer NOT NULL,
    avg_point double precision DEFAULT '0'::double precision NOT NULL,
    a_top integer DEFAULT 0 NOT NULL,
    a_last integer DEFAULT 0 NOT NULL,
    round integer DEFAULT 0 NOT NULL,
    win integer DEFAULT 0 NOT NULL,
    gun integer DEFAULT 0 NOT NULL,
    bark integer DEFAULT 0 NOT NULL,
    riichi integer DEFAULT 0 NOT NULL,
    win_point double precision DEFAULT '0'::double precision NOT NULL,
    gun_point double precision DEFAULT '0'::double precision NOT NULL,
    bark_point double precision DEFAULT '0'::double precision NOT NULL,
    riichi_point double precision DEFAULT '0'::double precision NOT NULL,
    ready integer DEFAULT 0 NOT NULL,
    ready_turn double precision DEFAULT '0'::double precision NOT NULL,
    win_turn double precision DEFAULT '0'::double precision NOT NULL,
    yaku_rci integer DEFAULT 0 NOT NULL,
    yaku_ipt integer DEFAULT 0 NOT NULL,
    yaku_tmo integer DEFAULT 0 NOT NULL,
    yaku_tny integer DEFAULT 0 NOT NULL,
    yaku_pnf integer DEFAULT 0 NOT NULL,
    yaku_y1y integer DEFAULT 0 NOT NULL,
    yaku_y2y integer DEFAULT 0 NOT NULL,
    yaku_y3y integer DEFAULT 0 NOT NULL,
    yaku_jk1 integer DEFAULT 0 NOT NULL,
    yaku_jk2 integer DEFAULT 0 NOT NULL,
    yaku_jk3 integer DEFAULT 0 NOT NULL,
    yaku_jk4 integer DEFAULT 0 NOT NULL,
    yaku_bk1 integer DEFAULT 0 NOT NULL,
    yaku_bk2 integer DEFAULT 0 NOT NULL,
    yaku_bk3 integer DEFAULT 0 NOT NULL,
    yaku_bk4 integer DEFAULT 0 NOT NULL,
    yaku_ipk integer DEFAULT 0 NOT NULL,
    yaku_rns integer DEFAULT 0 NOT NULL,
    yaku_hai integer DEFAULT 0 NOT NULL,
    yaku_hou integer DEFAULT 0 NOT NULL,
    yaku_ckn integer DEFAULT 0 NOT NULL,
    yaku_ss1 integer DEFAULT 0 NOT NULL,
    yaku_it1 integer DEFAULT 0 NOT NULL,
    yaku_ct1 integer DEFAULT 0 NOT NULL,
    yaku_wri integer DEFAULT 0 NOT NULL,
    yaku_ss2 integer DEFAULT 0 NOT NULL,
    yaku_it2 integer DEFAULT 0 NOT NULL,
    yaku_ct2 integer DEFAULT 0 NOT NULL,
    yaku_toi integer DEFAULT 0 NOT NULL,
    yaku_ctt integer DEFAULT 0 NOT NULL,
    yaku_sak integer DEFAULT 0 NOT NULL,
    yaku_skt integer DEFAULT 0 NOT NULL,
    yaku_stk integer DEFAULT 0 NOT NULL,
    yaku_hrt integer DEFAULT 0 NOT NULL,
    yaku_s3g integer DEFAULT 0 NOT NULL,
    yaku_h1t integer DEFAULT 0 NOT NULL,
    yaku_jc2 integer DEFAULT 0 NOT NULL,
    yaku_mnh integer DEFAULT 0 NOT NULL,
    yaku_jc3 integer DEFAULT 0 NOT NULL,
    yaku_rpk integer DEFAULT 0 NOT NULL,
    yaku_c1t integer DEFAULT 0 NOT NULL,
    yaku_mnc integer DEFAULT 0 NOT NULL,
    yaku_x13 integer DEFAULT 0 NOT NULL,
    yaku_xd3 integer DEFAULT 0 NOT NULL,
    yaku_x4a integer DEFAULT 0 NOT NULL,
    yaku_xt1 integer DEFAULT 0 NOT NULL,
    yaku_xs4 integer DEFAULT 0 NOT NULL,
    yaku_xd4 integer DEFAULT 0 NOT NULL,
    yaku_xcr integer DEFAULT 0 NOT NULL,
    yaku_xr1 integer DEFAULT 0 NOT NULL,
    yaku_xth integer DEFAULT 0 NOT NULL,
    yaku_xch integer DEFAULT 0 NOT NULL,
    yaku_x4k integer DEFAULT 0 NOT NULL,
    yaku_x9r integer DEFAULT 0 NOT NULL,
    yaku_w13 integer DEFAULT 0 NOT NULL,
    yaku_w4a integer DEFAULT 0 NOT NULL,
    yaku_w9r integer DEFAULT 0 NOT NULL,
    kzeykm integer DEFAULT 0 NOT NULL,
    han_rci double precision DEFAULT '0'::double precision NOT NULL,
    han_ipt double precision DEFAULT '0'::double precision NOT NULL,
    han_tmo double precision DEFAULT '0'::double precision NOT NULL,
    han_tny double precision DEFAULT '0'::double precision NOT NULL,
    han_pnf double precision DEFAULT '0'::double precision NOT NULL,
    han_y1y double precision DEFAULT '0'::double precision NOT NULL,
    han_y2y double precision DEFAULT '0'::double precision NOT NULL,
    han_y3y double precision DEFAULT '0'::double precision NOT NULL,
    han_jk1 double precision DEFAULT '0'::double precision NOT NULL,
    han_jk2 double precision DEFAULT '0'::double precision NOT NULL,
    han_jk3 double precision DEFAULT '0'::double precision NOT NULL,
    han_jk4 double precision DEFAULT '0'::double precision NOT NULL,
    han_bk1 double precision DEFAULT '0'::double precision NOT NULL,
    han_bk2 double precision DEFAULT '0'::double precision NOT NULL,
    han_bk3 double precision DEFAULT '0'::double precision NOT NULL,
    han_bk4 double precision DEFAULT '0'::double precision NOT NULL,
    han_ipk double precision DEFAULT '0'::double precision NOT NULL,
    han_rns double precision DEFAULT '0'::double precision NOT NULL,
    han_hai double precision DEFAULT '0'::double precision NOT NULL,
    han_hou double precision DEFAULT '0'::double precision NOT NULL,
    han_ckn double precision DEFAULT '0'::double precision NOT NULL,
    han_ss1 double precision DEFAULT '0'::double precision NOT NULL,
    han_it1 double precision DEFAULT '0'::double precision NOT NULL,
    han_ct1 double precision DEFAULT '0'::double precision NOT NULL,
    han_wri double precision DEFAULT '0'::double precision NOT NULL,
    han_ss2 double precision DEFAULT '0'::double precision NOT NULL,
    han_it2 double precision DEFAULT '0'::double precision NOT NULL,
    han_ct2 double precision DEFAULT '0'::double precision NOT NULL,
    han_toi double precision DEFAULT '0'::double precision NOT NULL,
    han_ctt double precision DEFAULT '0'::double precision NOT NULL,
    han_sak double precision DEFAULT '0'::double precision NOT NULL,
    han_skt double precision DEFAULT '0'::double precision NOT NULL,
    han_stk double precision DEFAULT '0'::double precision NOT NULL,
    han_hrt double precision DEFAULT '0'::double precision NOT NULL,
    han_s3g double precision DEFAULT '0'::double precision NOT NULL,
    han_h1t double precision DEFAULT '0'::double precision NOT NULL,
    han_jc2 double precision DEFAULT '0'::double precision NOT NULL,
    han_mnh double precision DEFAULT '0'::double precision NOT NULL,
    han_jc3 double precision DEFAULT '0'::double precision NOT NULL,
    han_rpk double precision DEFAULT '0'::double precision NOT NULL,
    han_c1t double precision DEFAULT '0'::double precision NOT NULL,
    han_mnc double precision DEFAULT '0'::double precision NOT NULL,
    yaku_dora integer DEFAULT 0 NOT NULL,
    yaku_uradora integer DEFAULT 0 NOT NULL,
    yaku_akadora integer DEFAULT 0 NOT NULL,
    yaku_kandora integer DEFAULT 0 NOT NULL,
    yaku_kanuradora integer DEFAULT 0 NOT NULL,
    ranks integer[] DEFAULT ARRAY[0, 0, 0, 0] NOT NULL
);


ALTER TABLE user_girl OWNER TO mako;

--
-- Name: users; Type: TABLE; Schema: public; Owner: mako
--

CREATE TABLE users (
    user_id bigint NOT NULL,
    username character varying(16) NOT NULL,
    password character(44) NOT NULL,
    level smallint DEFAULT 0 NOT NULL,
    pt smallint DEFAULT 0 NOT NULL,
    rating double precision DEFAULT '1500'::double precision NOT NULL
);


ALTER TABLE users OWNER TO mako;

--
-- Name: users_user_id_seq; Type: SEQUENCE; Schema: public; Owner: mako
--

CREATE SEQUENCE users_user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE users_user_id_seq OWNER TO mako;

--
-- Name: users_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mako
--

ALTER SEQUENCE users_user_id_seq OWNED BY users.user_id;


--
-- Name: replays replay_id; Type: DEFAULT; Schema: public; Owner: mako
--

ALTER TABLE ONLY replays ALTER COLUMN replay_id SET DEFAULT nextval('replays_replay_id_seq'::regclass);


--
-- Name: users user_id; Type: DEFAULT; Schema: public; Owner: mako
--

ALTER TABLE ONLY users ALTER COLUMN user_id SET DEFAULT nextval('users_user_id_seq'::regclass);


--
-- Data for Name: girls; Type: TABLE DATA; Schema: public; Owner: mako
--

COPY girls (girl_id, level, pt, rating) FROM stdin;
0	0	0	1500
710113	9	0	986.125999999999976
710114	11	580	1622.96800000000007
710115	19	4000	2437.85800000000017
712411	11	520	1665.23399999999992
712412	11	310	1602.26999999999998
712413	9	30	1289.4559999999999
712611	13	875	1844.98900000000003
712613	12	600	1717.65000000000009
712714	9	45	1414.42599999999993
712715	12	480	1534.51999999999998
712915	9	0	1291.08899999999994
713301	11	295	1434.73900000000003
713311	11	385	1672.24299999999994
713314	11	160	1524.60300000000007
713811	9	45	1149.84500000000003
713815	11	205	1655.67399999999998
714915	10	110	1356.58500000000004
715212	11	625	1670.96100000000001
990001	11	700	1654.43599999999992
990002	17	1000	2095.77300000000014
990003	10	50	1265.75800000000004
\.


--
-- Data for Name: replay_of_user; Type: TABLE DATA; Schema: public; Owner: mako
--

COPY replay_of_user (user_id, replay_ids) FROM stdin;
\.


--
-- Data for Name: replays; Type: TABLE DATA; Schema: public; Owner: mako
--

COPY replays (replay_id, content) FROM stdin;
\.


--
-- Name: replays_replay_id_seq; Type: SEQUENCE SET; Schema: public; Owner: mako
--

SELECT pg_catalog.setval('replays_replay_id_seq', 2112, true);


--
-- Data for Name: user_girl; Type: TABLE DATA; Schema: public; Owner: mako
--

COPY user_girl (user_id, girl_id, avg_point, a_top, a_last, round, win, gun, bark, riichi, win_point, gun_point, bark_point, riichi_point, ready, ready_turn, win_turn, yaku_rci, yaku_ipt, yaku_tmo, yaku_tny, yaku_pnf, yaku_y1y, yaku_y2y, yaku_y3y, yaku_jk1, yaku_jk2, yaku_jk3, yaku_jk4, yaku_bk1, yaku_bk2, yaku_bk3, yaku_bk4, yaku_ipk, yaku_rns, yaku_hai, yaku_hou, yaku_ckn, yaku_ss1, yaku_it1, yaku_ct1, yaku_wri, yaku_ss2, yaku_it2, yaku_ct2, yaku_toi, yaku_ctt, yaku_sak, yaku_skt, yaku_stk, yaku_hrt, yaku_s3g, yaku_h1t, yaku_jc2, yaku_mnh, yaku_jc3, yaku_rpk, yaku_c1t, yaku_mnc, yaku_x13, yaku_xd3, yaku_x4a, yaku_xt1, yaku_xs4, yaku_xd4, yaku_xcr, yaku_xr1, yaku_xth, yaku_xch, yaku_x4k, yaku_x9r, yaku_w13, yaku_w4a, yaku_w9r, kzeykm, han_rci, han_ipt, han_tmo, han_tny, han_pnf, han_y1y, han_y2y, han_y3y, han_jk1, han_jk2, han_jk3, han_jk4, han_bk1, han_bk2, han_bk3, han_bk4, han_ipk, han_rns, han_hai, han_hou, han_ckn, han_ss1, han_it1, han_ct1, han_wri, han_ss2, han_it2, han_ct2, han_toi, han_ctt, han_sak, han_skt, han_stk, han_hrt, han_s3g, han_h1t, han_jc2, han_mnh, han_jc3, han_rpk, han_c1t, han_mnc, yaku_dora, yaku_uradora, yaku_akadora, yaku_kandora, yaku_kanuradora, ranks) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: mako
--

COPY users (user_id, username, password, level, pt, rating) FROM stdin;
501	ⓝ喵打	iMqwzWh78q4vOGkh7ALsA6ohvh25OQ/VMDNFICkTarc=	10	200	1304.68499999999995
502	ⓝ打喵	iMqwzWh78q4vOGkh7ALsA6ohvh25OQ/VMDNFICkTarc=	9	45	1235.60599999999999
\.


--
-- Name: users_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: mako
--

SELECT pg_catalog.setval('users_user_id_seq', 1912, true);


--
-- Name: girls idx_18446_primary; Type: CONSTRAINT; Schema: public; Owner: mako
--

ALTER TABLE ONLY girls
    ADD CONSTRAINT idx_18446_primary PRIMARY KEY (girl_id);


--
-- Name: replays idx_18454_primary; Type: CONSTRAINT; Schema: public; Owner: mako
--

ALTER TABLE ONLY replays
    ADD CONSTRAINT idx_18454_primary PRIMARY KEY (replay_id);


--
-- Name: users idx_18463_primary; Type: CONSTRAINT; Schema: public; Owner: mako
--

ALTER TABLE ONLY users
    ADD CONSTRAINT idx_18463_primary PRIMARY KEY (user_id);


--
-- Name: user_girl idx_18470_primary; Type: CONSTRAINT; Schema: public; Owner: mako
--

ALTER TABLE ONLY user_girl
    ADD CONSTRAINT idx_18470_primary PRIMARY KEY (user_id, girl_id);


--
-- Name: replay_of_user replay_of_user_pkey; Type: CONSTRAINT; Schema: public; Owner: mako
--

ALTER TABLE ONLY replay_of_user
    ADD CONSTRAINT replay_of_user_pkey PRIMARY KEY (user_id);


--
-- Name: idx_18463_username; Type: INDEX; Schema: public; Owner: mako
--

CREATE UNIQUE INDEX idx_18463_username ON users USING btree (username);


--
-- Name: idx_18470_girl_id; Type: INDEX; Schema: public; Owner: mako
--

CREATE INDEX idx_18470_girl_id ON user_girl USING btree (girl_id);


--
-- Name: replay_of_user replay_of_user_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: mako
--

ALTER TABLE ONLY replay_of_user
    ADD CONSTRAINT replay_of_user_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(user_id);


--
-- Name: user_girl user_girl_ibfk_1; Type: FK CONSTRAINT; Schema: public; Owner: mako
--

ALTER TABLE ONLY user_girl
    ADD CONSTRAINT user_girl_ibfk_1 FOREIGN KEY (user_id) REFERENCES users(user_id) ON UPDATE RESTRICT ON DELETE RESTRICT;


--
-- Name: user_girl user_girl_ibfk_2; Type: FK CONSTRAINT; Schema: public; Owner: mako
--

ALTER TABLE ONLY user_girl
    ADD CONSTRAINT user_girl_ibfk_2 FOREIGN KEY (girl_id) REFERENCES girls(girl_id) ON UPDATE RESTRICT ON DELETE RESTRICT;


--
-- PostgreSQL database dump complete
--

