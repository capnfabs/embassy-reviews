# Embassy Reviews!

People review the darndest things.

## Get started
- Grab a key for the [Google Places API](https://developers.google.com/places/web-service/), and put it in `places_api_key.txt`
- Probably add a credit card to the project if you're doing development - the daily limits get exhausted quickly if you're fetching lots of reviews.


## What's in the repo?
I've designed this so you should be able to precompute a list of embassies, and then fetch reviews for them periodically. I'm not expecting the list of embassies to change much, so it's probably ok to save it to a file and recompute it every now and again.

- `./data/` - contains information on cities and countries. See [`data/README.md`](data/README.md) for details on sources.
- `./process_cities` - processes the data in `/data`, outputs a list of potentially interesting (lat,lng) pairs to use as inputs to the places API. Initially, this is just based on the capital city of each country.
- `./fetch_embassies` - Reads (lat,lng) pairs from the stdin and Fetches as many embassies within 50km of each of these as the Places API will allow. Outputs the Places API placeid for each of these. Note that this might return duplicates, you might want to pipe the output through `| sort | uniq`.
- `./fetch_review` - Takes a list of Place IDs and fetches reviews for all of them.
