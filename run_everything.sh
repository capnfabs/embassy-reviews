#!/bin/bash


go install ./fetch_embassies ./fetch_reviews ./process_cities ./process_reviews
process_cities > out/cities.txt
PLACES_API_KEY=`cat places_api_key.txt` fetch_embassies <out/cities.txt >out/embassies.txt
cat out/embassies.txt | sort | uniq >out/embassies_unique.txt
# Use shuf on linux, gshuf on osx
gshuf out/embassies_unique.txt >out/embassies_unique_shuffled.txt
# Actually fetch the reviews. Don't format them yet, just save them.
cat out/embassies_unique_shuffled.txt | PLACES_API_KEY=`cat places_api_key.txt` fetch_reviews >out/reviews_raw.json

# It might be a good idea to format the reviews now if you're going to manipulate them manually.
# First, make sure `jq` is installed.
jq . <out/reviews_raw.json >out/reviews_formatted.json
mv out/reviews_formatted.json out/reviews_raw.json

# Then, process the reviews.
cat out/reviews_raw.json | process_reviews >out/reviews.json
