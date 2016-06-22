#!/bin/bash
# Panic if can't complete a single instruction
set -e

go install ./fetch_embassies ./fetch_reviews ./process_cities
process_cities > out/cities.txt
PLACES_API_KEY=`cat places_api_key.txt` fetch_embassies <out/cities.txt >out/embassies.txt
cat out/embassies.txt | sort | uniq >out/embassies_unique.txt
# Use shuf on linux, gshuf on osx
gshuf out/embassies_unique.txt >out/embassies_unique_shuffled.txt || \
shuf out/embassies_unique.txt >out/embassies_unique_shuffled.txt
# Actually fetch the reviews. Don't format them yet, just save them.
cat out/embassies_unique_shuffled.txt | PLACES_API_KEY=`cat places_api_key.txt` fetch_reviews >out/reviews_raw.json

# Then, process the reviews.
( cd process_reviews; ./gradlew run )
# Output is in reviews_processed.json
cp out/reviews_processed.json tweeter/tweetsrc.json

./build_lambda.sh
