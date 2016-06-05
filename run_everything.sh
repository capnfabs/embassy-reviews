#!/bin/bash


go install ./fetch_embassies ./fetch_review ./process_cities
process_cities > out/cities.txt
PLACES_API_KEY=`cat places_api_key.txt` fetch_embassies <out/cities.txt >out/embassies.txt
cat out/embassies.txt | sort | uniq >out/embassies_unique.txt
# Use shuf on linux, gshuf on osx
gshuf -n 100 out/embassies_unique.txt | PLACES_API_KEY=`cat places_api_key.txt` fetch_review
