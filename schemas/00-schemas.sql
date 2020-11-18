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

-- remove existing enum type if exist
DO
$$
    DECLARE
        r RECORD;
    BEGIN
        FOR r IN (
            SELECT pe.enumtypid, pe.enumlabel, pt.typname
            FROM pg_enum pe
                     JOIN pg_type pt ON pt.oid = pe.enumtypid)
            LOOP
                EXECUTE 'DROP TYPE IF EXISTS ' || quote_ident(r.typname) || ' CASCADE';
            END LOOP;
    END
$$;

create table recipients
(
    id          serial primary key,
    name        varchar(100),
    attributes  jsonb default '{}',
    emails      varchar(100)[],
    phones      varchar(20)[],
    disabled    bool  default false,
    created_at  timestamp,
    updated_at  timestamp,
    disabled_at timestamp
);

create unique index recipients_name_unique on recipients (name);

create table audience
(
    id          serial primary key,
    title       varchar(50),
    emails      varchar(100)[],
    disabled    bool default false,
    created_at  timestamp,
    updated_at  timestamp,
    disabled_at timestamp
);

create unique index audience_title_unique on audience (title);

create type invitation_type as enum ('multi-link','single-link');

create table groups
(
    id              serial primary key,
    title           varchar(50),
    invitation_type invitation_type,
    audiences       integer[], /** audience collection **/
    disabled        boolean default false,
    published       boolean,
    created_at      timestamp,
    updated_at      timestamp,
    disabled_at     timestamp,
    published_at    timestamp
);

create unique index groups_title_unique on groups (title);

create type question_type as enum ('essay','choice');

create table questions
(
    id          serial primary key,
    sequence    integer, -- question number/sequence
    question    varchar(500),
    expect      question_type,
    options     varchar(150)[],
    group_id    integer
        constraint questions_groups_id references groups (id),
    mandatory   bool,
    disabled    boolean default false,
    created_at  timestamp,
    updated_at  timestamp,
    disabled_at timestamp
);

create table links
(
    id           serial primary key,
    hash         varchar(512),
    pin          varchar(10),
    group_id     integer
        constraint links_group_id references groups (id),
    disabled     boolean default false,
    published    bool,
    usage_limit  integer default 0,
    created_at   timestamp,
    updated_at   timestamp,
    disabled_at  timestamp,
    published_at timestamp
);

create unique index links_hash_unique on links (hash);

create table link_visits
(
    id      serial primary key,
    link_id integer,
    at      timestamp,
    agent   text,
    ref     jsonb default '{}'
);

create type mark_as_type as enum ('favorite');

create table submission
(
    id              serial primary key,
    hash            varchar(128),
    question_id     integer,
    question_number integer,
    question        varchar(500),
    group_id        integer,
    group_title     varchar(100),
    invitation_type invitation_type,
    expect          question_type,
    options         jsonb default '[]',
    answer_choice   smallint,
    answer_essay    text,
    marked_as       mark_as_type[],
    created_at      timestamp,
    updated_at      timestamp
);

create index on submission (hash);

create table submission_audience
(
    id                  serial primary key,
    submission_group_id integer,
    audiences           varchar(100)[],
    audience_title      varchar(50)
);


/** for admin dashboard **/
create type user_role as enum ('sysadmin','admin','member','invitation','custom');

create table users
(
    id           serial primary key,
    username     varchar(25) unique,
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
    id             serial primary key,
    role           user_role,
    path           varchar(255),
    method_allowed request_method[],
    disabled       boolean default true,
    access_level   access_level
);

create table user_access
(
    id             serial primary key,
    user_id        integer,
    scope          scope   default 'custom',
    path           varchar(255),
    method_allowed request_method[],
    access_level   access_level,
    disabled       boolean default false
);