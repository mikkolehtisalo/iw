--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = off;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET escape_string_warning = off;

SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: activities; Type: TABLE; Schema: public; Owner: wiki; Tablespace: 
--

CREATE TABLE activities (
    activity_id uuid NOT NULL,
    "timestamp" timestamp with time zone,
    user_id character varying(32),
    user_name character varying(32),
    activity_type character varying(32),
    target_type character varying(16),
    target_title character varying(128),
    target_id character varying(128),
    readacl character varying(1024),
    writeacl character varying(1024),
    adminacl character varying(1024)
);


ALTER TABLE public.activities OWNER TO wiki;

--
-- Name: attachments; Type: TABLE; Schema: public; Owner: wiki; Tablespace: 
--

CREATE TABLE attachments (
    attachment_id uuid NOT NULL,
    wiki_id uuid NOT NULL,
    attachment bytea,
    mime character varying,
    filename character varying,
    modified timestamp with time zone NOT NULL,
    status character varying,
    create_user character varying(32)
);


ALTER TABLE public.attachments OWNER TO wiki;

--
-- Name: contentfields; Type: TABLE; Schema: public; Owner: wiki; Tablespace: 
--

CREATE TABLE contentfields (
    contentfield_id uuid NOT NULL,
    wiki_id uuid NOT NULL,
    content text,
    modified timestamp with time zone NOT NULL,
    status character varying,
    create_user character varying(32)
);


ALTER TABLE public.contentfields OWNER TO wiki;

--
-- Name: favoritewikis; Type: TABLE; Schema: public; Owner: wiki; Tablespace: 
--

CREATE TABLE favoritewikis (
    username character varying NOT NULL,
    wiki_id uuid NOT NULL,
    modified timestamp with time zone NOT NULL,
    status character varying
);


ALTER TABLE public.favoritewikis OWNER TO wiki;

--
-- Name: locks; Type: TABLE; Schema: public; Owner: wiki; Tablespace: 
--

CREATE TABLE locks (
    target_id uuid NOT NULL,
    wiki_id uuid NOT NULL,
    username character varying(32),
    realname character varying(32),
    modified timestamp with time zone
);


ALTER TABLE public.locks OWNER TO wiki;

--
-- Name: pages; Type: TABLE; Schema: public; Owner: wiki; Tablespace: 
--

CREATE TABLE pages (
    page_id uuid NOT NULL,
    wiki_id uuid NOT NULL,
    path character varying(256),
    title character varying(256),
    create_user character varying(32),
    readacl character varying(32),
    writeacl character varying(32),
    adminacl character varying(32),
    stopinheritation boolean,
    index integer,
    depth integer,
    modified timestamp with time zone NOT NULL,
    status character varying
);


ALTER TABLE public.pages OWNER TO wiki;

--
-- Name: wikis; Type: TABLE; Schema: public; Owner: wiki; Tablespace: 
--

CREATE TABLE wikis (
    wiki_id uuid NOT NULL,
    title character varying(128),
    description text,
    create_user character varying(32),
    readacl character varying(1024),
    writeacl character varying(1024),
    adminacl character varying(1024),
    modified timestamp with time zone NOT NULL,
    status character varying
);


ALTER TABLE public.wikis OWNER TO wiki;

--
-- Name: activities_pkey; Type: CONSTRAINT; Schema: public; Owner: wiki; Tablespace: 
--

ALTER TABLE ONLY activities
    ADD CONSTRAINT activities_pkey PRIMARY KEY (activity_id);


--
-- Name: attachments_pkey; Type: CONSTRAINT; Schema: public; Owner: wiki; Tablespace: 
--

ALTER TABLE ONLY attachments
    ADD CONSTRAINT attachments_pkey PRIMARY KEY (attachment_id, wiki_id, modified);


--
-- Name: contentfields_pkey; Type: CONSTRAINT; Schema: public; Owner: wiki; Tablespace: 
--

ALTER TABLE ONLY contentfields
    ADD CONSTRAINT contentfields_pkey PRIMARY KEY (contentfield_id, wiki_id, modified);


--
-- Name: favoritewikis_pkey; Type: CONSTRAINT; Schema: public; Owner: wiki; Tablespace: 
--

ALTER TABLE ONLY favoritewikis
    ADD CONSTRAINT favoritewikis_pkey PRIMARY KEY (username, wiki_id, modified);


--
-- Name: locks_pkey; Type: CONSTRAINT; Schema: public; Owner: wiki; Tablespace: 
--

ALTER TABLE ONLY locks
    ADD CONSTRAINT locks_pkey PRIMARY KEY (target_id, wiki_id);


--
-- Name: pages_pkey; Type: CONSTRAINT; Schema: public; Owner: wiki; Tablespace: 
--

ALTER TABLE ONLY pages
    ADD CONSTRAINT pages_pkey PRIMARY KEY (page_id, wiki_id, modified);


--
-- Name: wikis_pkey; Type: CONSTRAINT; Schema: public; Owner: wiki; Tablespace: 
--

ALTER TABLE ONLY wikis
    ADD CONSTRAINT wikis_pkey PRIMARY KEY (wiki_id, modified);


--
-- Name: ix_activities_activity_id; Type: INDEX; Schema: public; Owner: wiki; Tablespace: 
--

CREATE INDEX ix_activities_activity_id ON activities USING btree (activity_id);


--
-- Name: ix_activities_timestamp; Type: INDEX; Schema: public; Owner: wiki; Tablespace: 
--

CREATE INDEX ix_activities_timestamp ON activities USING btree ("timestamp");


--
-- Name: ix_activities_user_id; Type: INDEX; Schema: public; Owner: wiki; Tablespace: 
--

CREATE INDEX ix_activities_user_id ON activities USING btree (user_id);


--
-- Name: ix_activities_user_name; Type: INDEX; Schema: public; Owner: wiki; Tablespace: 
--

CREATE INDEX ix_activities_user_name ON activities USING btree (user_name);


--
-- Name: ix_attachments_attachment_id; Type: INDEX; Schema: public; Owner: wiki; Tablespace: 
--

CREATE INDEX ix_attachments_attachment_id ON attachments USING btree (attachment_id);


--
-- Name: ix_attachments_modified; Type: INDEX; Schema: public; Owner: wiki; Tablespace: 
--

CREATE INDEX ix_attachments_modified ON attachments USING btree (modified);


--
-- Name: ix_attachments_status; Type: INDEX; Schema: public; Owner: wiki; Tablespace: 
--

CREATE INDEX ix_attachments_status ON attachments USING btree (status);


--
-- Name: ix_attachments_wiki_id; Type: INDEX; Schema: public; Owner: wiki; Tablespace: 
--

CREATE INDEX ix_attachments_wiki_id ON attachments USING btree (wiki_id);


--
-- Name: ix_contentfields_contentfield_id; Type: INDEX; Schema: public; Owner: wiki; Tablespace: 
--

CREATE INDEX ix_contentfields_contentfield_id ON contentfields USING btree (contentfield_id);


--
-- Name: ix_contentfields_modified; Type: INDEX; Schema: public; Owner: wiki; Tablespace: 
--

CREATE INDEX ix_contentfields_modified ON contentfields USING btree (modified);


--
-- Name: ix_contentfields_status; Type: INDEX; Schema: public; Owner: wiki; Tablespace: 
--

CREATE INDEX ix_contentfields_status ON contentfields USING btree (status);


--
-- Name: ix_contentfields_wiki_id; Type: INDEX; Schema: public; Owner: wiki; Tablespace: 
--

CREATE INDEX ix_contentfields_wiki_id ON contentfields USING btree (wiki_id);


--
-- Name: ix_favoritewikis_modified; Type: INDEX; Schema: public; Owner: wiki; Tablespace: 
--

CREATE INDEX ix_favoritewikis_modified ON favoritewikis USING btree (modified);


--
-- Name: ix_favoritewikis_status; Type: INDEX; Schema: public; Owner: wiki; Tablespace: 
--

CREATE INDEX ix_favoritewikis_status ON favoritewikis USING btree (status);


--
-- Name: ix_favoritewikis_username; Type: INDEX; Schema: public; Owner: wiki; Tablespace: 
--

CREATE INDEX ix_favoritewikis_username ON favoritewikis USING btree (username);


--
-- Name: ix_favoritewikis_wiki_id; Type: INDEX; Schema: public; Owner: wiki; Tablespace: 
--

CREATE INDEX ix_favoritewikis_wiki_id ON favoritewikis USING btree (wiki_id);


--
-- Name: ix_locks_page_id; Type: INDEX; Schema: public; Owner: wiki; Tablespace: 
--

CREATE INDEX ix_locks_page_id ON locks USING btree (target_id);


--
-- Name: ix_locks_wiki_id; Type: INDEX; Schema: public; Owner: wiki; Tablespace: 
--

CREATE INDEX ix_locks_wiki_id ON locks USING btree (wiki_id);


--
-- Name: ix_pages_modified; Type: INDEX; Schema: public; Owner: wiki; Tablespace: 
--

CREATE INDEX ix_pages_modified ON pages USING btree (modified);


--
-- Name: ix_pages_page_id; Type: INDEX; Schema: public; Owner: wiki; Tablespace: 
--

CREATE INDEX ix_pages_page_id ON pages USING btree (page_id);


--
-- Name: ix_pages_path; Type: INDEX; Schema: public; Owner: wiki; Tablespace: 
--

CREATE INDEX ix_pages_path ON pages USING btree (path);


--
-- Name: ix_pages_status; Type: INDEX; Schema: public; Owner: wiki; Tablespace: 
--

CREATE INDEX ix_pages_status ON pages USING btree (status);


--
-- Name: ix_pages_wiki_id; Type: INDEX; Schema: public; Owner: wiki; Tablespace: 
--

CREATE INDEX ix_pages_wiki_id ON pages USING btree (wiki_id);


--
-- Name: ix_wikis_modified; Type: INDEX; Schema: public; Owner: wiki; Tablespace: 
--

CREATE INDEX ix_wikis_modified ON wikis USING btree (modified);


--
-- Name: ix_wikis_status; Type: INDEX; Schema: public; Owner: wiki; Tablespace: 
--

CREATE INDEX ix_wikis_status ON wikis USING btree (status);


--
-- Name: ix_wikis_wiki_id; Type: INDEX; Schema: public; Owner: wiki; Tablespace: 
--

CREATE INDEX ix_wikis_wiki_id ON wikis USING btree (wiki_id);


--
-- Name: public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

