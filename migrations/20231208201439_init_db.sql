-- +goose Up
CREATE TABLE IF NOT EXISTS user_stats
(
    id SERIAL PRIMARY KEY NOT NULL UNIQUE,
    message_count INTEGER NOT NULL default 0,
    dick_plus_count INTEGER NOT NULL default 0,
    dick_minus_count INTEGER NOT NULL default 0,
    yes_count INTEGER NOT NULL default 0,
    no_count INTEGER NOT NULL default 0,
    duels_count INTEGER NOT NULL default 0,
    duels_win_count INTEGER NOT NULL default 0,
    duels_lose_count INTEGER NOT NULL default 0,
    kill_count INTEGER NOT NULL default 0,
    die_count INTEGER NOT NULL default 0,
    gay_count INTEGER NOT NULL default 0
);

CREATE TABLE IF NOT EXISTS users
(
    id SERIAL PRIMARY KEY NOT NULL UNIQUE,
    tg_id BIGINT NOT NULL,
    chat_id BIGINT NOT NULL,
    is_bot BOOLEAN NOT NULL,
    is_premium BOOLEAN NOT NULL,
    first_name VARCHAR,
    last_name VARCHAR,
    username VARCHAR,
    dick_size INTEGER NOT NULL DEFAULT 1,
    change_dick_at TIMESTAMP WITH TIME ZONE NOT NULL,
    user_stat_id SERIAL,
    FOREIGN KEY (user_stat_id) REFERENCES user_stats,
    health_points INTEGER NOT NULL DEFAULT 3,
    hp_taked_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_gay BOOLEAN NOT NULL,
    gay_at TIMESTAMP WITH TIME ZONE NOT NULL,
    points INTEGER NOT NULL DEFAULT 0,
    cur_dick_change_count INTEGER NOT NULL default 0,
    max_dick_change_count INTEGER NOT NULL default 0
);

CREATE TABLE IF NOT EXISTS calendars
(
    chat_id BIGINT PRIMARY KEY NOT NULL UNIQUE,
    calendar_id VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS gays
(
    id serial PRIMARY KEY NOT NULL UNIQUE,
    chat_id BIGINT NOT NULL,
    tg_id INTEGER NOT NULL,
    username VARCHAR,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS homeworks
(
    id SERIAL PRIMARY KEY NOT NULL UNIQUE,
    chat_id BIGINT NOT NULL,
    subject VARCHAR,
    task VARCHAR,
    created_at TIMESTAMP WITH TIME ZONE
);


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
