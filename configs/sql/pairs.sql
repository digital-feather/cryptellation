CREATE DATABASE IF NOT EXISTS pairs;

CREATE TABLE IF NOT EXISTS pairs
(
    base_symbol TEXT NOT NULL,
    quote_symbol TEXT NOT NULL,
    PRIMARY KEY (base_symbol, quote_symbol)
);