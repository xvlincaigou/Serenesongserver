from pymongo import MongoClient
import json

# client = MongoClient('mongodb://root:example@mongodb:27017/mydb?authSource=admin')
client = MongoClient('mongodb://localhost:27017/')

database   = client['serenesong']
collection = database['Ci']
if collection.count_documents({}) != 0:
    print("Collection already contains data. No insertion performed.")
    exit()

for iter in range(0, 22):
    with open(f'./database/songci/ci.song.{iter*1000}.json', 'r', encoding='utf-8') as f:
        data = json.load(f)
        
    new_data = []
    for item in data:
        new_item = {}
        for key in item:
            if   key == 'paragraphs':
                new_content = []
                for sentence in item[key]:
                    if sentence == " >> " or sentence == "词牌介绍":
                        continue
                    new_content.append(sentence)
                new_item['content'] = new_content
            elif key == 'rhythmic':
                if '・' in item[key]:
                    cipai_list = []
                    for cipai in item[key].split('・'):
                        cipai_list.append(cipai)
                    new_item['cipai'] = cipai_list
                else:
                    new_item['cipai'] = [item[key]]
            elif key == 'prologue':
                new_item['xiaoxu'] = item[key]
            else:
                new_item[key] = item[key]
        new_data.append(new_item)
        
    result = collection.insert_many(new_data)
    print(f"Data ci.song.{iter*1000}.json inserted successfully.")
    
with open(f'./database/songci/ci.song.2019y.json', 'r', encoding='utf-8') as f:
    data = json.load(f)
    
new_data = []
for item in data:
    new_item = {}
    for key in item:
        if key == 'paragraphs':
            new_item['content'] = item[key]
        elif key == 'rhythmic':
            if '・' in item[key]:
                cipai_list = []
                for cipai in item[key].split('・'):
                    cipai_list.append(cipai)
                new_item['cipai'] = cipai_list
            else:
                new_item['cipai'] = [item[key]]
        elif key == 'prologue':
            new_item['xiaoxu'] = item[key]
        else:
            new_item[key] = item[key]
    new_data.append(new_item)
    
result = collection.insert_many(new_data)
print(f"Data ci.song.2019y.json inserted successfully.")

client.close()