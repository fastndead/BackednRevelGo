--
-- PostgreSQL database dump
--

-- Dumped from database version 9.5.14
-- Dumped by pg_dump version 9.5.14

-- Started on 2018-08-27 09:46:53

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 7 (class 2615 OID 16400)
-- Name: airport; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA airport;


ALTER SCHEMA airport OWNER TO postgres;

--
-- TOC entry 1 (class 3079 OID 12355)
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- TOC entry 2138 (class 0 OID 0)
-- Dependencies: 1
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


--
-- TOC entry 187 (class 1259 OID 16615)
-- Name: flights_seq; Type: SEQUENCE; Schema: airport; Owner: postgres
--

CREATE SEQUENCE airport.flights_seq
    START WITH 1
    INCREMENT BY 1
    MINVALUE 0
    MAXVALUE 9999
    CACHE 1;


ALTER TABLE airport.flights_seq OWNER TO postgres;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- TOC entry 181 (class 1259 OID 16401)
-- Name: flights; Type: TABLE; Schema: airport; Owner: postgres
--

CREATE TABLE airport.flights (
    c_id integer DEFAULT nextval('airport.flights_seq'::regclass) NOT NULL,
    c_arrival_point character varying(30) NOT NULL,
    c_departure_point character varying(30) NOT NULL,
    c_fk_planes integer NOT NULL
);


ALTER TABLE airport.flights OWNER TO postgres;

--
-- TOC entry 185 (class 1259 OID 16451)
-- Name: pilots_seq; Type: SEQUENCE; Schema: airport; Owner: postgres
--

CREATE SEQUENCE airport.pilots_seq
    START WITH 1
    INCREMENT BY 1
    MINVALUE 0
    MAXVALUE 9999
    CACHE 1;


ALTER TABLE airport.pilots_seq OWNER TO postgres;

--
-- TOC entry 183 (class 1259 OID 16423)
-- Name: pilots; Type: TABLE; Schema: airport; Owner: postgres
--

CREATE TABLE airport.pilots (
    c_id integer DEFAULT nextval('airport.pilots_seq'::regclass) NOT NULL,
    c_first_name character varying,
    c_last_name character varying
);


ALTER TABLE airport.pilots OWNER TO postgres;

--
-- TOC entry 186 (class 1259 OID 16454)
-- Name: planes_seq; Type: SEQUENCE; Schema: airport; Owner: postgres
--

CREATE SEQUENCE airport.planes_seq
    START WITH 1
    INCREMENT BY 1
    MINVALUE 0
    MAXVALUE 9999
    CACHE 1;


ALTER TABLE airport.planes_seq OWNER TO postgres;

--
-- TOC entry 182 (class 1259 OID 16410)
-- Name: planes; Type: TABLE; Schema: airport; Owner: postgres
--

CREATE TABLE airport.planes (
    c_id integer DEFAULT nextval('airport.planes_seq'::regclass) NOT NULL,
    c_name character varying(30) NOT NULL
);


ALTER TABLE airport.planes OWNER TO postgres;

--
-- TOC entry 188 (class 1259 OID 16633)
-- Name: toc_flights_pilots_seq; Type: SEQUENCE; Schema: airport; Owner: postgres
--

CREATE SEQUENCE airport.toc_flights_pilots_seq
    START WITH 1
    INCREMENT BY 1
    MINVALUE 0
    MAXVALUE 9999
    CACHE 1;


ALTER TABLE airport.toc_flights_pilots_seq OWNER TO postgres;

--
-- TOC entry 184 (class 1259 OID 16434)
-- Name: toc_flights_pilots; Type: TABLE; Schema: airport; Owner: postgres
--

CREATE TABLE airport.toc_flights_pilots (
    c_id integer DEFAULT nextval('airport.toc_flights_pilots_seq'::regclass) NOT NULL,
    c_fk_flight integer NOT NULL,
    c_fk_pilot integer NOT NULL
);


ALTER TABLE airport.toc_flights_pilots OWNER TO postgres;

--
-- TOC entry 2007 (class 2606 OID 16414)
-- Name: flights_pkey; Type: CONSTRAINT; Schema: airport; Owner: postgres
--

ALTER TABLE ONLY airport.flights
    ADD CONSTRAINT flights_pkey PRIMARY KEY (c_id);


--
-- TOC entry 2011 (class 2606 OID 16427)
-- Name: pilots_pkey; Type: CONSTRAINT; Schema: airport; Owner: postgres
--

ALTER TABLE ONLY airport.pilots
    ADD CONSTRAINT pilots_pkey PRIMARY KEY (c_id);


--
-- TOC entry 2009 (class 2606 OID 16416)
-- Name: planes_pkey; Type: CONSTRAINT; Schema: airport; Owner: postgres
--

ALTER TABLE ONLY airport.planes
    ADD CONSTRAINT planes_pkey PRIMARY KEY (c_id);


--
-- TOC entry 2015 (class 2606 OID 16438)
-- Name: toc_flights_pilots_pkey; Type: CONSTRAINT; Schema: airport; Owner: postgres
--

ALTER TABLE ONLY airport.toc_flights_pilots
    ADD CONSTRAINT toc_flights_pilots_pkey PRIMARY KEY (c_id);


--
-- TOC entry 2005 (class 1259 OID 16578)
-- Name: fki_fk; Type: INDEX; Schema: airport; Owner: postgres
--

CREATE INDEX fki_fk ON airport.flights USING btree (c_fk_planes);


--
-- TOC entry 2012 (class 1259 OID 16444)
-- Name: fki_flight_toc; Type: INDEX; Schema: airport; Owner: postgres
--

CREATE INDEX fki_flight_toc ON airport.toc_flights_pilots USING btree (c_fk_flight);


--
-- TOC entry 2013 (class 1259 OID 16450)
-- Name: fki_pilot_toc; Type: INDEX; Schema: airport; Owner: postgres
--

CREATE INDEX fki_pilot_toc ON airport.toc_flights_pilots USING btree (c_fk_pilot);


--
-- TOC entry 2016 (class 2606 OID 16646)
-- Name: flights_c_fk_planes_fkey; Type: FK CONSTRAINT; Schema: airport; Owner: postgres
--

ALTER TABLE ONLY airport.flights
    ADD CONSTRAINT flights_c_fk_planes_fkey FOREIGN KEY (c_fk_planes) REFERENCES airport.planes(c_id);


--
-- TOC entry 2017 (class 2606 OID 16651)
-- Name: toc_flights_pilots_c_fk_flight_fkey; Type: FK CONSTRAINT; Schema: airport; Owner: postgres
--

ALTER TABLE ONLY airport.toc_flights_pilots
    ADD CONSTRAINT toc_flights_pilots_c_fk_flight_fkey FOREIGN KEY (c_fk_flight) REFERENCES airport.flights(c_id);


-- Completed on 2018-08-27 09:46:54

--
-- PostgreSQL database dump complete
--

