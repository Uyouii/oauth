use oauth_db;

drop table partner_tab;

create table if not exists partner_tab(
    id bigint primary key auto_increment,
    partner_name varchar(64),
    partner_key varchar(64) unique,
    partner_secret varchar(64),
    expire bigint default 3600,
    create_timestamp bigint,
    update_timestamp bigint
);


drop table token_tab;

create table if not exists token_tab(
    id bigint primary key auto_increment,
    partner_key varchar(64),
    token varchar(128) unique,
    expire bigint,
    expire_timestamp bigint,
    create_timestamp bigint,
    update_timestamp bigint
);

create index pratner_token_index on token_tab (partner_key,token);
create index expire_timestamp on token_tab(expire_timestamp);