CREATE DATABASE IF NOT EXISTS exchanges;

CREATE TABLE IF NOT EXISTS exchanges
(
    name TEXT NOT NULL,
    fees NUMERIC NOT NULL,
    last_sync_time TIMESTAMP NOT NULL,
    PRIMARY KEY (name)
);

CREATE TABLE IF NOT EXISTS pairs
(
    symbol TEXT NOT NULL,
    PRIMARY KEY (symbol)
);

CREATE TABLE IF NOT EXISTS periods
(
    symbol TEXT NOT NULL,
    PRIMARY KEY (symbol)
);

CREATE TABLE IF NOT EXISTS exchanges_pairs
(
    exchange_name TEXT NOT NULL,
    pair_symbol TEXT NOT NULL,
    PRIMARY KEY (exchange_name, pair_symbol),
    CONSTRAINT fk_exchange_pair_exchange FOREIGN KEY (exchange_name) REFERENCES exchanges(name) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_exchange_pair_pair FOREIGN KEY (pair_symbol) REFERENCES pairs(symbol) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS exchanges_periods
(
    exchange_name TEXT NOT NULL,
    period_symbol TEXT NOT NULL,
    PRIMARY KEY (exchange_name, period_symbol),
    CONSTRAINT fk_exchange_period_exchange FOREIGN KEY (exchange_name) REFERENCES exchanges(name) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_exchange_period_pair FOREIGN KEY (period_symbol) REFERENCES periods(symbol) ON DELETE CASCADE ON UPDATE CASCADE
);