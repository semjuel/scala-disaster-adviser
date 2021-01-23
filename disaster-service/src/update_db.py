import requests
from datetime import datetime
import pandas as pd
from pymongo import MongoClient
import time
import json
import asyncio
import rx
import rx.operators as ops
from rx.scheduler.eventloop import AsyncIOScheduler
import logging
from confluent_kafka import Producer
import os
import socket


conf = {'bootstrap.servers': os.environ['KAFKA_HOST'],
        'client.id': socket.gethostname()}
producer = Producer(conf)
key_counter = 1


def send_notification(event):
    global key_counter
    if len(event['geometry']) > 0 and len(event['geometry'][0]['coordinates']) > 1:
        notification = {
            "disaster": {"description": event['title'],
                     "date": int(datetime.timestamp(pd.to_datetime(event['geometry'][0]['date'])))*1000,
                     "lat": event['geometry'][0]['coordinates'][0],
                     "lon": event['geometry'][0]['coordinates'][1]}
        }
        print(notification)
        producer.produce(os.environ['KAFKA_TOPIC'], key=str(++key_counter),
                         value=bytes(json.dumps(notification), 'utf-8'))
    return event


async def foo():
    print("Next interval")
    return requests.get("https://eonet.sci.gsfc.nasa.gov/api/v3/events").json()['events']


def intervalRead(rate, fun, collection) -> rx.Observable:
    loop = asyncio.get_event_loop()
    return rx.interval(rate).pipe(
        ops.map(lambda i: rx.from_future(loop.create_task(fun()))),
        ops.merge_all(),
        ops.map(lambda i: list(filter(lambda j: collection.find_one({'id': j['id']}) is None, i))),
        ops.map(lambda i: [send_notification(j) for j in i])
    )


async def main(loop):
    conn = MongoClient(os.environ['MONGO_HOST'], username=os.environ['MONGO_USER'], password=os.environ['MONGO_PASS'])
    print("Connected successfully!!!")

    collection = conn.disaster.events
    obs = intervalRead(3, foo, collection)
    obs.subscribe(
        on_next=lambda i: [collection.insert_one(j) for j in i],
        scheduler=AsyncIOScheduler(loop)
    )


loop = asyncio.get_event_loop()
loop.create_task(main(loop))
loop.run_forever()
