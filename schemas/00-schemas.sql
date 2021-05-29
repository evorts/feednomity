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

/* enabling crypt on psql */
create
    extension if not exists pgcrypto;

/* create sequence for users role id */
create sequence if not exists users_role_id_seq
    start with 1
    increment by 1
    minvalue 1
    no maxvalue
    cache 1;

alter sequence users_role_id_seq restart with 1;

/* create sequence for users id */
create sequence if not exists users_id_seq
    start with 1
    increment by 1
    minvalue 1
    no maxvalue
    cache 1;

alter sequence users_id_seq restart with 1;

/* create sequence for links id */
create sequence if not exists links_id_seq
    start with 1
    increment by 1
    minvalue 1
    no maxvalue
    cache 1;

alter sequence links_id_seq restart with 1;

/* create sequence for links hash */
create sequence if not exists links_helper_seq
    start with 1
    increment by 1
    minvalue 1
    no maxvalue
    cache 1;

alter sequence links_helper_seq restart with 1;

/* create sequence for users id */
create sequence if not exists dist_object_id_seq
    start with 1
    increment by 1
    minvalue 1
    no maxvalue
    cache 1;

alter sequence dist_object_id_seq restart with 1;

/** for admin dashboard **/
/** could be utilise as company **/
create table users_organization
(
    id          serial primary key,
    name        varchar(40) unique,
    address     varchar(150) default '',
    phone       varchar(15)  default '',
    disabled    bool         default false,
    created_at  timestamp,
    updated_at  timestamp,
    disabled_at timestamp
);

create table users_group
(
    id          serial primary key,
    name        varchar(40) unique,
    disabled    bool default false,
    org_id      int
        constraint users_group_org_id references users_organization (id),
    created_at  timestamp,
    updated_at  timestamp,
    disabled_at timestamp
);

create type user_role as enum ('sysadmin','admin-org','admin-group','member','guest','custom');

create table users
(
    id           serial primary key,
    username     varchar(25),
    display_name varchar(50),
    attributes   jsonb        default '{}',
    email        varchar(50) not null,
    phone        varchar(15),
    password     varchar(128) default null,
    pin          varchar(6)   default null,
    access_role  user_role,
    job_role     varchar(50),
    assignment   varchar(50),
    group_id     int
        constraint users_group_id references users_group (id),
    disabled     bool         default false,
    verified     bool         default false,
    created_at   timestamp,
    updated_at   timestamp,
    disabled_at  timestamp
);

create unique index users_username on users (username);
create unique index users_email_index on users (email);
create index users_attributes_index on users using gin (attributes);

/** as audit log trail **/
create table user_activities
(
    id          bigserial primary key,
    user_id     int
        constraint activities_users_id references users (id),
    action      varchar(50),
    values      jsonb default '{}',
    values_prev jsonb default '{}',
    notes       varchar(100),
    at          timestamp
);

create type access_level as enum ('get', 'post', 'put', 'delete', 'head', 'options');
create type access_scope as enum ('self', 'group', 'org', 'global', 'none');

/** default role access **/
create table role_access
(
    id                serial primary key,
    role              user_role,
    path              varchar(100), /** should consistent pattern such as <module>.<method> **/
    regex             bool    default false,
    access_allowed    access_level[],
    access_disallowed access_level[],
    access_scope      access_scope,
    disabled          boolean default false,
    created_at        timestamp,
    updated_at        timestamp,
    disabled_at       timestamp
);

create table role_users_limit
(
    id         serial primary key,
    role       user_role unique,
    max_user   int,
    created_at timestamp,
    updated_at timestamp
);

create table user_access
(
    id                serial primary key,
    user_id           integer,
    path              varchar(100), /** should consistent pattern such as <module>.<method> **/
    regex             bool    default false,
    access_allowed    access_level[],
    access_disallowed access_level[],
    access_scope      access_scope,
    disabled          boolean default false,
    created_at        timestamp,
    updated_at        timestamp,
    disabled_at       timestamp
);

create
    unique index idx_user_access_id_path ON user_access (user_id, path);

create table links
(
    id           serial primary key,
    hash         varchar(128),
    pin          varchar(64),
    disabled     boolean   default false,
    published    bool      default false,
    usage_limit  integer   default 0,
    attributes   jsonb     default '{}',
    created_by   int
        constraint links_created_by references users (id),
    updated_by   int
        constraint links_updated_by references users (id),
    expired_at   timestamp default now(),
    created_at   timestamp,
    updated_at   timestamp,
    disabled_at  timestamp,
    published_at timestamp
);

create
    unique index links_hash_unique on links (hash);

create table link_visits
(
    id      serial primary key,
    link_id integer,
    at      timestamp,
    agent   text,
    ref     jsonb default '{}'
);

create table distributions
(
    id                 serial primary key,
    topic              varchar(100),
    disabled           bool default false,
    archived           bool default false,
    distributed        bool default false,
    distribution_limit int, /** max limit distribution **/
    distribution_count int, /** how many times its distributed **/
    range_start        timestamp, /** review start **/
    range_end          timestamp, /** review end **/
    created_by         int
        constraint distributions_created_by_users_id references users (id),
    updated_by         int
        constraint distributions_updated_by_users_id references users (id),
    for_group_id       int
        constraint distributions_users_group_id references users_group (id),
    created_at         timestamp,
    updated_at         timestamp,
    disabled_at        timestamp,
    archived_at        timestamp,
    distributed_at     timestamp
);

create type distribution_object_status as enum ('none', 'sent', 'failed');

create table distribution_objects
(
    id                serial primary key,
    distribution_id   int
        constraint distribution_recipients_distributions_id references distributions (id),
    recipient_id      int
        constraint distribution_recipients_recipient_id_recipients_id references users (id),
    respondent_id     int
        constraint distribution_respondents_respondent_id references users (id),
    link_id           int
        constraint distribution_objects_links_id references links (id),
    publishing_status distribution_object_status default 'none', /** when its published -- sent to respondent **/
    publishing_log    jsonb                      default '[]',
    retry_count       int,
    created_by        int
        constraint distribution_objects_created_by references users (id),
    updated_by        int
        constraint distribution_objects_updated_by references users (id),
    created_at        timestamp,
    updated_at        timestamp,
    published_at      timestamp
);

create unique index idx_unique_distribution_item on distribution_objects (distribution_id, recipient_id, respondent_id);
create index idx_distribution_objects_publishing_status on distribution_objects (publishing_status);

create table distribution_mail_queue
(
    id                     bigserial primary key,
    distribution_object_id int,
    recipient_id           int,
    respondent_id          int,
    from_email             varchar(100),
    to_email               varchar(100),
    subject                varchar(200),
    template               varchar(50),
    arguments              jsonb default '{}'
);

create unique index idx_unique_distributions_queue on distribution_mail_queue (distribution_object_id, recipient_id, respondent_id);

create table distribution_log
(
    id          bigserial primary key,
    action      varchar(50),
    values      jsonb default '{}',
    values_prev jsonb default '{}',
    notes       varchar(100),
    at          timestamp
);

/**
  draft => recipient already submit their feedback but still in draft
  final => recipient already finalize their submission -- cannot be change
 */
create type feedback_status as enum ('not-started','draft', 'final');

create table feedbacks
(
    id                     serial primary key,
    distribution_id        int,
    distribution_topic     varchar(100),
    distribution_object_id int,
    range_start            timestamp,
    range_end              timestamp,
    respondent_id          int,
    respondent_username    varchar(100),
    respondent_name        varchar(100),
    respondent_email       varchar(100),
    respondent_group_id    int,
    respondent_group_name  varchar(50),
    respondent_org_id      int,
    respondent_org_name    varchar(50),
    respondent_role        varchar(50),
    respondent_assigment   varchar(50),
    recipient_id           int,
    recipient_username     varchar(100),
    recipient_name         varchar(100),
    recipient_email        varchar(100),
    recipient_group_id     int,
    recipient_group_name   varchar(50),
    recipient_org_id       int,
    recipient_org_name     varchar(50),
    recipient_role         varchar(50),
    recipient_assignment   varchar(50),
    link_id                int,
    hash                   varchar(128),
    status                 feedback_status,
    content                jsonb default '{}',
    created_at             timestamp,
    updated_at             timestamp
);

create unique index idx_unique_feedbacks on feedbacks (distribution_id, distribution_object_id, respondent_id, recipient_id);

create table feedback_log
(
    id          bigserial primary key,
    feedback_id int
        constraint feedback_log_feedbacks_id references feedbacks (id),
    action      varchar(50),
    values      jsonb default '{}',
    values_prev jsonb default '{}',
    notes       text,
    at          timestamp
);

create type question_type as enum ('essay','choice');

create table questions
(
    id          serial primary key,
    sequence    integer, -- question number/sequence
    question    varchar(500),
    expect      question_type,
    options     varchar(150)[],
    mandatory   bool    default true,
    disabled    boolean default false,
    created_at  timestamp,
    updated_at  timestamp,
    disabled_at timestamp
);

/* for page based on template */
create table pages
(
    id          serial primary key,
    name        varchar(50),
    template    varchar(250),
    /* format:
        {
            "name": "",
            "validation": []
        }
     */
    validations jsonb default '[]'
);

/* for dynamic forms scaffolding */
create table forms
(
    id          serial primary key,
    template    varchar(250),         -- path of the template
    /* format:
       {
            "type": "text|checklist|dropdown|choice",
            "name": "",
            "attributes": [],
            "values": "",
            "default": "",
            "validation": ["notEmpty", "r'regex-pattern"],
            "mandatory": true
       }
     */
    forms       jsonb   default '[]', -- scaffolding dynamic form
    disabled    boolean default false,
    created_at  timestamp,
    updated_at  timestamp,
    disabled_at timestamp
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
    expect          question_type,
    options         jsonb default '[]',
    answer_choice   smallint,
    answer_essay    text,
    marked_as       mark_as_type[],
    created_at      timestamp,
    updated_at      timestamp
);

create
    index on submission (hash);

create table submission_audience
(
    id                  serial primary key,
    submission_group_id integer,
    audiences           varchar(100)[],
    audience_title      varchar(50)
);
