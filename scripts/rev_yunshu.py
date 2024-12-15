import json

with open('./database/rhymes/Pingshui_Rhyme.json', 'r', encoding='utf-8') as f:
    data = json.load(f)

new_data = {}

for tone, rhymes in data.items():
    for rhyme, characters in rhymes.items():
        characters = list(set(characters))
        for char in characters:
            if char not in new_data:
                new_data[char] = []
            new_data[char].append({
                "Tone": tone,
                "Rhyme": rhyme
            })

with open('./database/rhymes/Reversed_Pingshui_Rhyme.json', 'w', encoding='utf-8') as f:
    json.dump(new_data, f, ensure_ascii=False, indent=4)

