from pymongo import MongoClient
import json

with open('./database/rhymes/Word_Tune.json', 'r', encoding='utf-8') as f:
    data = json.load(f)
    
new_data = {}
for item, info in data.items():
    new_data[item] = info
    
# client = MongoClient('mongodb://root:example@mongodb:27017/mydb?authSource=admin')
client = MongoClient('mongodb://localhost:27017/')

database   = client['serenesong']
collection = database['CharacterTune']
if collection.count_documents({}) == 0:
    result = collection.insert_one(new_data)
    print("Data inserted successfully.")
else:
    print("Collection already contains data. No insertion performed.")
    
client.close()