import requests
import pandas as pd
from flask import Flask
from datetime import datetime
from pymongo import MongoClient

app = Flask(__name__)


@app.route("/new")#, methods=['POST']
def echo():
    #data = request.json
    data = {
      "user" : "username",
      "date": {
        "timestamp": "1601840340",
        "range": 100000
      },
      "location": {
        "lat": 39.5,
        "lon": 48.5,
        "range": 0.5
      }
    }
    print(f"Got {data['user']} info")

    conn = MongoClient()
    print("Connected successfully!!!")
    # create database
    db = conn.disaster
    collection = db.events

    lat_min = data['location']['lat']-data['location']['range']
    lat_max = data['location']['lat']+data['location']['range']
    lon_min = data['location']['lon']-data['location']['range']
    lon_max = data['location']['lon']+data['location']['range']

    date_min = datetime.fromtimestamp(int(data['date']['timestamp'])-data['date']['range'])
    date_max = datetime.fromtimestamp(int(data['date']['timestamp'])+data['date']['range'])

    answer = list(collection.find({"$and": [{'geometry.0.coordinates.0': {'$gt': lat_min, '$lt': lat_max}},
                                         {'geometry.0.coordinates.1': {'$gt': lon_min, '$lt': lon_max}}]}))
    #,{'geometry.0.date': {'$gt': date_min, '$lt': date_max}}
    answer = {"events": [form_answer(i) for i in answer]}
    print(answer)
    return "ok"

# 2020-10-04T22:39:00Z
def form_answer(item):
    return {"description": item['title'],
                 "date": datetime.timestamp(pd.to_datetime(item['geometry'][0]['date'])),
                 "lat": item['geometry'][0]['coordinates'][0],
                 "lon": item['geometry'][0]['coordinates'][1]}


if __name__ == "__main__":

    app.run()

