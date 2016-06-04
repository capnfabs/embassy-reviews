## data sources

## cities15000.txt, countryInfo.txt

These files were downloaded from http://download.geonames.org/export/dump/ on 3 Jun 2016.

 - cities15000.txt is the extracted version of [cities15000.zip](http://download.geonames.org/export/dump/cities15000.zip)
 - [countryInfo.txt](http://download.geonames.org/export/dump/countryInfo.txt) was downloaded as-is.

## sample_places_response.json

Sample response from the Places API.

## sample_fetch_embassies_response.txt

stdout from running

```sh
go install ./fetch_embassies && echo "-35.3018436,149.1241899" | PLACES_API_KEY=`cat places_api_key.txt` fetch_embassies
```
