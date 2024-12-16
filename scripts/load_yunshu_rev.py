from pymongo import MongoClient
import json

with open('./database/rhymes/Reversed_Pingshui_Rhyme.json', 'r', encoding='utf-8') as f:
    data = json.load(f)
    
new_data = {}
for tone, info in data.items():
    new_data[tone] = info
    
# client = MongoClient('mongodb://root:example@mongodb:27017/mydb?authSource=admin')
client = MongoClient('mongodb://localhost:27017/')

database   = client['serenesong']
collection = database['Characters']
if collection.count_documents({}) == 0:
    result = collection.insert_one(new_data)
    print("Data inserted successfully.")
else:
    print("Collection already contains data. No insertion performed.")

client.close()