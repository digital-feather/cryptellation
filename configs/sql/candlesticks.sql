CREATE DATABASE IF NOT EXISTS candlesticks;

CREATE TABLE IF NOT EXISTS candlesticks
(
    exchange_name TEXT NOT NULL,
    pair_symbol TEXT NOT NULL,
    period_symbol TEXT NOT NULL,
    time TIMESTAMP NOT NULL,
    open NUMERIC NOT NULL,
    high NUMERIC NOT NULL,
    low NUMERIC NOT NULL,
    close NUMERIC NOT NULL,
    volume NUMERIC NOT NULL,
    uncomplete BOOLEAN NOT NULL,
    PRIMARY KEY (exchange_name, pair_symbol, period_symbol, time)
)