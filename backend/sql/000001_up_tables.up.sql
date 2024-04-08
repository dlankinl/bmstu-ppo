create schema if not exists ppo;

create table if not exists ppo.users(
    id uuid primary key default gen_random_uuid(),
    username varchar(256) unique,
    full_name varchar(256),
    birthday date,
    gender varchar(1),
    city varchar(128),
    password varchar(256)
);

create table if not exists ppo.skills(
    id uuid primary key default gen_random_uuid(),
    name varchar(64) not null,
    description text not null
);

create table if not exists ppo.user_skills(
    user_id uuid not null,
    skill_id uuid not null
);

alter table ppo.user_skills
add constraint u_s_pk primary key (user_id, skill_id);

create table if not exists ppo.contacts(
    id uuid primary key default gen_random_uuid(),
    owner_id uuid not null,
    name varchar(64),
    value varchar(128)
);

create table if not exists ppo.fin_reports(
    id uuid primary key default gen_random_uuid(),
    company_id uuid not null,
    revenue float4 not null,
    costs float4 not null,
    year int not null,
    quarter int not null
);

create table if not exists ppo.companies(
    id uuid primary key default gen_random_uuid(),
    owner_id uuid not null,
    activity_field_id uuid not null,
    name varchar(128),
    city varchar(128)
);

create table if not exists ppo.activity_fields(
    id uuid primary key default gen_random_uuid(),
    name varchar(128) not null,
    description text not null,
    cost float4 not null
);
