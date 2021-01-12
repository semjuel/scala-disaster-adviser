from pymongo import MongoClient


def delete_one(id_):
    conn = MongoClient()
    print("Connected successfully!!!")
    # database
    db = conn.disaster
    # switched to collection
    collection = db.events

    i = 0
    cursor = collection.find({})
    for document in cursor:
        i+=1
    print("Events before: ", i)

    collection.remove({'id': id_})

    i = 0
    cursor = collection.find({})
    for document in cursor:
        i+=1
    print("Events after", i)


def delete_all():
    conn = MongoClient()
    print("Connected successfully!!!")
    # database
    db = conn.disaster
    # switched to collection
    collection = db.events

    i = 0
    cursor = collection.find({})
    for document in cursor:
        i+=1
    print("Events before: ", i)

    collection.remove()

    i = 0
    cursor = collection.find({})
    for document in cursor:
        i+=1
    print("Events after", i)

if __name__ == "__main__":
    delete_one('EONET_5198')
    #delete_all()