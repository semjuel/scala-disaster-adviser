import requests
import time
import json
from datetime import datetime
import pandas as pd
from pymongo import MongoClient
from confluent_kafka import Producer
import socket


def send_notification(event):
    notification = {
        "disaster": {"description": event['title'],
                 "date": datetime.timestamp(pd.to_datetime(event['geometry'][0]['date'])),
                 "lat": event['geometry'][0]['coordinates'][0],
                 "lon": event['geometry'][0]['coordinates'][1]}
    }
    print(notification)
    # TO DO change url
    # r = requests.post('http://127.0.0.1:5000/new', json=json.dumps(notification))
    # print(r.status_code)

    conf = {'bootstrap.servers': "host1:9092,host2:9092",
            'client.id': socket.gethostname()}
    producer = Producer(conf)
    producer.produce(topic, key=event['id'], value=notification)

def get_events():
    conn = MongoClient()
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

        time.sleep(3)  # Sleep for 3 seconds


if __name__ == "__main__":
    get_events()
