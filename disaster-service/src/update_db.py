import requests
import time
import json
from datetime import datetime
import pandas as pd
from pymongo import MongoClient
from confluent_kafka import Producer
import os
import socket


def send_notification(event):
    print(event)
    notification = {
        "disaster": {"description": event['title'],
                 "date": datetime.timestamp(pd.to_datetime(event['geometry'][0]['date'])),
                 "lat": event['geometry'][0]['coordinates'][0],
                 "lon": event['geometry'][0]['coordinates'][1]}
    }
    print(notification)

    conf = {'bootstrap.servers': os.environ['KAFKA_HOST'],
            'client.id': socket.gethostname()}
    producer = Producer(conf)
    producer.produce(os.environ['KAFKA_TOPIC'], key=event['id'], value=bytes(str(notification), 'utf-8'))


def get_events():
    conn = MongoClient(os.environ['MONGO_HOST'], username=os.environ['MONGO_USER'], password=os.environ['MONGO_PASS'])
    print("Connected successfully!!!")

    db = conn.disaster # create database
    collection = db.events # create connection

    while True:
        response = requests.get("https://eonet.sci.gsfc.nasa.gov/api/v3/events").json()['events']

        print(f"Got {len(response)} events from api")
        # save to mongo
        for event in response:
            if collection.find_one({'id':  event['id']}) is None:
                collection.insert_one(event)
                send_notification(event)

        time.sleep(60)  # Sleep for 60 seconds


if __name__ == "__main__":
    get_events()
