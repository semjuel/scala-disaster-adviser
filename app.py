import requests
import pandas as pd
from flask import Flask, request, jsonify
from datetime import datetime
from pymongo import MongoClient

app = Flask(__name__)


@app.route("/new", methods=['POST'])#
def echo():
    data = request.get_json()
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
    print(date_min, date_max)

    answer = list(collection.find({"$and": [{'geometry.0.coordinates.0': {'$gt': lat_min, '$lt': lat_max}},
                                         {'geometry.0.coordinates.1': {'$gt': lon_min, '$lt': lon_max}},
                                            {'geometry.0.date': {'$gt': str(date_min), '$lt': str(date_max)}}]}))

    answer = {"events": [form_answer(i) for i in answer]}
    print(answer)
    return jsonify(answer)


def form_answer(item):
    return {"description": item['title'],
                 "date": datetime.timestamp(pd.to_datetime(item['geometry'][0]['date'])),
                 "lat": item['geometry'][0]['coordinates'][0],
                 "lon": item['geometry'][0]['coordinates'][1]}


if __name__ == "__main__":
    app.run()

