CREATE TABLE IF NOT EXISTS phones (
    phone_id SERIAL PRIMARY KEY,
    manufacturer VARCHAR(255) NOT NULL,
    model_tag VARCHAR(255) NOT NULL,
    model_number VARCHAR(255) UNIQUE NOT NULL,
    os_version VARCHAR(255) NOT NULL,
    api_version VARCHAR(255) NOT NULL,
    cpu VARCHAR(255) NOT NULL,
    firmware VARCHAR(255) NOT NULL,
    bootloader VARCHAR(255) NOT NULL,
    supported_archs TEXT[] NOT NULL,
    sim_slots INT NOT NULL DEFAULT 0,
    sd_slots  INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS sim_cards (
    sim_card_id SERIAL PRIMARY KEY,
    phone_id INT REFERENCES phones (phone_id),
    phone_number VARCHAR(255) UNIQUE NOT NULL,
    operator VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS sd_cards (
    sd_card_id SERIAL PRIMARY KEY,
    phone_id INT REFERENCES phones (phone_id),
    sd_manufacturer_id VARCHAR(255) NOT NULL,
    serial_no VARCHAR(255) UNIQUE NOT NULL,
    total_space INT NOT NULL,
    used_space INT NOT NULL,
    free_space INT NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code INT NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS user_phone (
    user_id INT REFERENCES users (user_id),
    phone_id INT REFERENCES phones (phone_id) UNIQUE
);

CREATE TABLE IF NOT EXISTS notifications (
    notification_id SERIAL PRIMARY KEY,
    model_tag VARCHAR(255) NOT NULL,
    source VARCHAR(255) NOT NULL,
    sender VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    timestamp BIGINT DEFAULT 0
);