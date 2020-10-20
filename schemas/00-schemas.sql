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

create table groups (
    id serial,
    title varchar(50),
    audience jsonb default '[]',
    disabled boolean default false,
    created_date timestamp,
    updated_date timestamp
);

create type question_type as enum ('essay','choice');

create table questions (
    id serial,
    question varchar(500),
    expect question_type,
    options jsonb default '[]',
    group_id int constraint questions_groups_id references groups(id),
    disabled boolean default false,
    created_date timestamp,
    updated_date timestamp
);

create table links (
    id serial,
    hash varchar(128) unique,
    pin varchar(10),
    question_group_id int constraint links_questions_group_id references groups(id),
    disabled boolean default false,
    created_date timestamp,
    updated_date timestamp
);