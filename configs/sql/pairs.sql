CREATE DATABASE IF NOT EXISTS pairs;

CREATE TABLE IF NOT EXISTS pairs
(
    symbol TEXT NOT NULL,
    base_asset_symbol TEXT NOT NULL,
    quote_asset_symbol TEXT NOT NULL,
    PRIMARY KEY (symbol)
);