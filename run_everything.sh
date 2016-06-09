#!/bin/bash


go install ./fetch_embassies ./fetch_reviews ./process_cities
process_cities > out/cities.txt
PLACES_API_KEY=`cat places_api_key.txt` fetch_embassies <out/cities.txt >out/embassies.txt
cat out/embassies.txt | sort | uniq >out/embassies_unique.txt
# Use shuf on linux, gshuf on osx
gshuf out/embassies_unique.txt >out/embassies_unique_shuffled.txt
# Don't worry, we also print reviews to Stdout.
 cat out/embassies_unique_shuffled.txt | PLACES_API_KEY=`cat places_api_key.txt` fetch_reviews --all >out/reviews.txt

# If you want to manually curate, install 'jq' and then
jq . <out/reviews.json >out/reviews_formatted.json
mv out/reviews_formatted.json out/reviews.json
