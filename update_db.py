import requests
import time
from datetime import datetime
import pandas as pd
from pymongo import MongoClient


def send_notification(event):
    notification = {
        "date": {"description": event['title'],
                 "date": datetime.timestamp(pd.to_datetime(event['geometry'][0]['date'])),
                 "lat": event['geometry'][0]['coordinates'][0],
                 "lon": event['geometry'][0]['coordinates'][1]}
    }
    print(notification)

def get_events():
    conn = MongoClient()
    print("Connected successfully!!!")
    # create database
    db = conn.disaster
    collection = db.events
    # collection_history = db.historical_events
    while True:
        response = requests.get("https://eonet.sci.gsfc.nasa.gov/api/v3/events").json()['events']

        print(f"Got {len(response)} events from api")
        # save to mongo
        for event in response:
            if collection.find_one({'id':  event['id']}) is None:
                collection.insert_one(event)
                send_notification(event)

        time.sleep(1)  # Sleep for 3 seconds

        response = requests.get("https://eonet.sci.gsfc.nasa.gov/api/v3/events/geojson").json()['features']

        # print(f"Got {len(response)} history events from api")
        # # save to mongo
        # for event in response:
        #     if collection_history.find_one({'properties.id':  event['properties']['id']}) is None:
        #         collection_history.insert_one(event)
        #         #send_notification(event)
        #
        # time.sleep(10)  # Sleep for 3 seconds

if __name__ == "__main__":
    #db_connector()
    get_events()