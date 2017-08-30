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

SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: repositories; Type: TABLE; Schema: public; Owner: flavio
--

CREATE TABLE repositories (
    id integer NOT NULL,
    id_of_repository_on_github integer,
    repository_name character varying(191),
    repository_owner character varying(191),
    default_branch character varying(100),
    created_at character varying(20),
    added_at timestamp(0) without time zone,
    enabled boolean DEFAULT true NOT NULL,
    private boolean DEFAULT false NOT NULL,
    fork boolean DEFAULT false NOT NULL,
    description text DEFAULT ''::text,
    initialized boolean DEFAULT false NOT NULL,
    has_static_public_page boolean DEFAULT false NOT NULL,
    repository_created_months_ago text,
    total_stars integer DEFAULT 0,
    total_issues_opened integer DEFAULT 0,
    total_commits integer DEFAULT 0,
    commits_count_last_12_months integer DEFAULT 0,
    commits_count_last_4_weeks integer DEFAULT 0,
    commits_count_last_week integer DEFAULT 0,
    stars_count_last_12_months integer DEFAULT 0,
    stars_count_last_4_weeks integer DEFAULT 0,
    stars_count_last_week integer DEFAULT 0,
    stars_per_month text
);


ALTER TABLE repositories OWNER TO flavio;

--
-- Name: repositories_id_seq; Type: SEQUENCE; Schema: public; Owner: flavio
--

CREATE SEQUENCE repositories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE repositories_id_seq OWNER TO flavio;

--
-- Name: repositories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: flavio
--

ALTER SEQUENCE repositories_id_seq OWNED BY repositories.id;


--
-- Name: repositories id; Type: DEFAULT; Schema: public; Owner: flavio
--

ALTER TABLE ONLY repositories ALTER COLUMN id SET DEFAULT nextval('repositories_id_seq'::regclass);


--
-- Name: repositories repositories_pkey; Type: CONSTRAINT; Schema: public; Owner: flavio
--

ALTER TABLE ONLY repositories
    ADD CONSTRAINT repositories_pkey PRIMARY KEY (id);


--
-- Name: repositories repositories_repository_id_unique; Type: CONSTRAINT; Schema: public; Owner: flavio
--

ALTER TABLE ONLY repositories
    ADD CONSTRAINT repositories_repository_id_unique UNIQUE (id_of_repository_on_github);


--
-- PostgreSQL database dump complete
--