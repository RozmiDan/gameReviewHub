import json

input_file = "load-test/games.csv"
output_file = "load-test/get_games_urls.csv"

with open(input_file, "r", encoding="utf-8") as fin, open(output_file, "w", encoding="utf-8") as fout:
    for line in fin:
        if not line.strip():
            continue
        entry = json.loads(line)
        fout.write(f"/games/{entry['id']}\n")