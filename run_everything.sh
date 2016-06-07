#!/bin/bash


go install ./fetch_embassies ./fetch_review ./process_cities
process_cities > out/cities.txt
PLACES_API_KEY=`cat places_api_key.txt` fetch_embassies <out/cities.txt >out/embassies.txt
cat out/embassies.txt | sort | uniq >out/embassies_unique.txt
# Use shuf on linux, gshuf on osx
gshuf out/embassies_unique.txt >out/embassies_unique_shuffled.txt
# Don't worry, we also print reviews to Stdout.
 cat out/embassies_unique_shuffled.txt | PLACES_API_KEY=`cat places_api_key.txt` fetch_review --all >out/reviews.txt
