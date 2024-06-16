create schema if not exists ppo;

create table if not exists ppo.users(
    id uuid primary key default gen_random_uuid(),
    username varchar(256) unique,
    full_name varchar(256),
    birthday date,
    gender varchar(1),
    city varchar(128),
    password varchar(256),
    role varchar(32)
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
    name varchar(64) not null,
    value varchar(128) not null
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
    name varchar(128) not null ,
    city varchar(128) not null
);

create table if not exists ppo.activity_fields(
    id uuid primary key default gen_random_uuid(),
    name varchar(128) not null,
    description text not null,
    cost float4 not null
);

create table if not exists ppo.reviews(
    id uuid primary key default gen_random_uuid(),
    target_id uuid not null,
    reviewer_id uuid not null,
    pros text not null,
    cons text not null,
    description text,
    rating int not null
);

alter table ppo.activity_fields add constraint chck_cost check ( cost >= 0.0 );

alter table ppo.user_skills add constraint fk_user foreign key (user_id) references ppo.users(id);
alter table ppo.user_skills add constraint fk_skill foreign key (skill_id) references ppo.skills(id);

alter table ppo.contacts add constraint fk_user foreign key (owner_id) references ppo.users(id);

alter table ppo.fin_reports add constraint fk_company foreign key (company_id) references ppo.companies(id);
alter table ppo.fin_reports add constraint chk_revenue check ( revenue >= 0.0 );
alter table ppo.fin_reports add constraint chk_costs check ( costs >= 0.0 );
alter table ppo.fin_reports add constraint chk_year check ( year > 0 );
alter table ppo.fin_reports add constraint chk_quarter check ( quarter >= 1 and quarter <= 4 );

alter table ppo.companies add constraint fk_owner foreign key (owner_id) references ppo.users(id);
alter table ppo.companies add constraint fk_activity_field foreign key (activity_field_id) references ppo.activity_fields(id);

alter table ppo.reviews add constraint fk_target foreign key (target_id) references ppo.users(id);
alter table ppo.reviews add constraint fk_reviewer foreign key (reviewer_id) references ppo.users(id);
alter table ppo.reviews add constraint chk_rating check ( rating >= 1 and rating <= 5 );


insert into ppo.users(username, password, role)
values ('admin', '$2a$10$4MYWtRfOlgU9smD01vZCFel4WmfsXc2RHuQm6Wq.uUezTeYb3HrNm', 'admin');


create user user_accountant;
grant accountant to user_accountant;

create user user_admin;
grant admin to user_admin;

create user user_editor;
grant editor to user_editor;