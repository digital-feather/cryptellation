# Cryptellation

Cryptellation is a **scalable cryptocurrency investment system**.

This system allows developers to create bots to manage their investments on 
different cryptographic markets, featuring **backtesting**, **livetesting** and 
**live running**.

## Supported clients

* Python (documentation incoming...)

## Services 

| Service          | Description                             |
| ---------------- | --------------------------------------- |
| **Backtests**    | Execute backtests                       |
| **Candlesticks** | Get cached informations on candlesticks |
| **Exchanges**    | Get cached informations on exchanges    |
| **Livetests**    | Execute livetests                       |
| **Ticks**        | Get ticks from exchanges                |


## Running Python example

### Requirements

* docker
* docker-compose
* pip

### How to

First launch the cryptellation system:

    docker-compose up -d

Then you can use the client to execute an example:

    cd clients/python
    pip install -r requirements.txt
    pip install -e .
    python examples/graph.py # Or any other from examples/ directory
