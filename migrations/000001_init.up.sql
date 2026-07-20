create table if not exists users (
    id uuid primary key,
    username varchar(255) unique not null,
    password varchar(255) not null,
    created_at timestamp default now(),
    updated_at timestamp default now()
);

create table if not exists crypto_coins (
    coin_id varchar(255) primary key,
    symbol varchar(50) not null
);

create table if not exists crypto_positions (
    id uuid primary key,
    user_id uuid not null references users(id),
    coin_id varchar(255) not null references crypto_coins(coin_id),
    amount numeric(38, 18) not null,
    avg_price_usd_cents bigint
);

create table if not exists bank_deposits (
    id uuid primary key,
    user_id uuid not null references users(id),
    name varchar(255) not null,
    currency varchar(3) not null,
    amount_cents bigint not null,
    interest_rate numeric(6, 3) not null,
    opened_at timestamp not null,
    period_months int not null,
    total_income_cents bigint not null
);

create table if not exists real_estates (
    id uuid primary key,
    user_id uuid not null references users(id),
    name varchar(255) not null,
    currency varchar(3) not null,
    purchase_price_cents bigint not null,
    monthly_income_cents bigint,
    purchased timestamp
);

create table if not exists stocks (
    id uuid primary key,
    user_id uuid not null references users(id),
    ticker varchar(10) not null,
    amount bigint not null,
    avg_price_cents bigint not null,
    currency varchar(3) not null
);

create table if not exists bonds (
    id uuid primary key,
    user_id uuid not null references users(id),
    ticker varchar(100) not null,
    amount bigint not null,
    avg_price_cents bigint not null,
    coupon_cents bigint not null,
    coupon_period_months int not null,
    currency varchar(3) not null
);