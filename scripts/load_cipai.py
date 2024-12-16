from pymongo import MongoClient
import json

with open('./database/rhymes/Ci_Tunes.json', 'r', encoding='utf-8') as f:
    data = json.load(f)
    
new_data = {}
for item, info in data.items():
    new_data[item] = info
    
client = MongoClient('mongodb://root:example@mongodb:27017/mydb?authSource=admin')

database   = client['serenesong']
collection = database['CipaiList']
collection.delete_many({})
result     = collection.insert_one(new_data)

print("Data inserted successfully.")
client.close()