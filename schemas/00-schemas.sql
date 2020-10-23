-- remove existing table if exist
DO
$$
    DECLARE
        r RECORD;
    BEGIN
        FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public')
            LOOP
                EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
            END LOOP;
    END
$$;

-- remove existing sequence if exist
DO
$$
    DECLARE
        r RECORD;
    BEGIN
        FOR r IN (SELECT relname FROM pg_class WHERE relkind = 'S')
            LOOP
                EXECUTE 'DROP SEQUENCE IF EXISTS ' || quote_ident(r.relname) || ' CASCADE';
            END LOOP;
    END
$$;

create table audience
(
    id     serial,
    title  varchar(50),
    emails varchar(100)[]
);

create type invitation_type as enum ('multi-link','single-link');

create table groups
(
    id              serial,
    title           varchar(100),
    invitation_type invitation_type,
    audience        integer[], /** audience collection **/
    disabled        boolean default false,
    published       boolean,
    created_date    timestamp,
    updated_date    timestamp,
    published_date  timestamp
);

create type question_type as enum ('essay','choice');

create table questions
(
    id           serial,
    question     varchar(500),
    expect       question_type,
    options      jsonb   default '[]',
    group_id     integer
        constraint questions_groups_id references groups (id),
    disabled     boolean default false,
    created_date timestamp,
    updated_date timestamp
);

create table links
(
    id           serial,
    hash         varchar(128) unique,
    pin          varchar(10),
    group_id     integer
        constraint links_group_id references groups (id),
    disabled     boolean default false,
    usage_limit  integer default 0,
    created_date timestamp,
    updated_date timestamp
);

create table link_visits
(
    id      serial,
    link_id integer,
    at      timestamp,
    agent   text,
    ref     varchar(255)
);

create table submission
(
    id              serial,
    hash            varchar(128),
    question_id     integer,
    question        varchar(500),
    group_id        integer,
    group_title     varchar(100),
    invitation_type invitation_type,
    expect          question_type,
    options         jsonb default '[]',
    answer_choice   varchar(50),
    answer_essay    text
);

create index on submission (hash);

create table submission_audience
(
    id                  serial,
    submission_group_id integer,
    audiences           varchar(100)[],
    audience_title      varchar(50)
);


/** for admin dashboard **/
create type user_role as enum ('sysadmin','admin','member','custom');

create table users
(
    id           serial,
    username     varchar(25),
    display_name varchar(50),
    email        varchar(50),
    phone        varchar(15),
    password     varchar(128),
    role         user_role,
    created_date timestamp,
    updated_date timestamp
);

create type scope as enum ('custom', 'all');
create type access_level as enum ('ro','wo','rw');
create type request_method as enum ('get', 'post', 'put', 'delete', 'head', 'options');

create table role_access
(
    id             serial,
    role           user_role,
    path           varchar(255),
    method_allowed request_method[],
    disabled       boolean default true,
    access_level   access_level
);

create table user_access
(
    id             serial,
    user_id        integer,
    scope          scope   default 'custom',
    path           varchar(255),
    method_allowed request_method[],
    access_level   access_level,
    disabled       boolean default false
);