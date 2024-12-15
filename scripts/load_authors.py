from pymongo import MongoClient
import json

with open(f'./database/songci/author.song.json', 'r', encoding='utf-8') as f:
    data = json.load(f)
    
new_data = []

for item in data: 
    if item['description'] == '--':
        continue
    new_item = {}
    for key in item:
        if   key == 'description':
            new_item['bio'] = item[key]
        elif key == 'name':
            new_item['name'] = item[key]
    new_data.append(new_item)
    
client = MongoClient('mongodb://localhost:27017/')

database   = client['serenesong']
collection = database['Author']
collection.delete_many({})
result     = collection.insert_many(new_data)

print("Data inserted successfully.")
client.close()