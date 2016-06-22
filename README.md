# Embassy Reviews!

> ★★★ The worst part? You don't have an option whether you like it or not - https://maps.google.com/?cid=15088333349668026791

People review the darndest things. In this case, embassies.

## What / Why?

For all user-generated content on the Internet, there's a signal-to-noise ratio. My theory is that in situations where the abstraction doesn't work, there's very little signal -- it's hard to create meaningful content in a context that's mostly nonsensical. Turns out the noise is often pretty funny.

Reviews of embassies fit the bill for this because reviews were designed to help individuals make comparisons between competitors. It works well for any situation where there's differentiation  - Cafes, restaurants, and software are classic examples, which is why they're rated so frequently. It works less well for pure commodities, where you'd probably only write a review if something really bad happened - for example, at a petrol / gas station, or a chain supermarket. It works terribly for things that are completely non-competitive - for example, [bridges](https://www.google.com.au/maps/place/Anzac+Bridge+Sydney/@-33.8692218,151.1836338,17z/data=!3m1!4b1!4m5!3m4!1s0x6b12afccc0abca15:0xa4b3d2d3d71ca7bb!8m2!3d-33.8692218!4d151.1858225) ("The supporting cables used to vibrate, then they put supporting supporting cables. Good looking bridge, actually quite pleasant to walk across.").

To me, embassies are the epitome of non-competitiveness - you often have to go there to get a visa, but people don't choose countries to go to based on reviews of embassies. You're stuck with it, no matter how you feel about it.

## The code / how it works

### Get started
- Grab a key for the [Google Places API](https://developers.google.com/places/web-service/), and put it in `places_api_key.txt`
- Probably add a credit card to the project if you're doing development - the daily limits get exhausted quickly if you're fetching lots of reviews.

### What's in the repo?
I've designed this in a series of discrete stages that all have output files, because fetching complete data sets is slow, and so I wanted to be able to iterate on each stage. Each stage is contained within a folder (more or less):

- `./data/` - contains information on cities and countries. See [`data/README.md`](data/README.md) for details on sources.
- `./process_cities` - processes the data in `/data`, outputs a list of potentially interesting (lat,lng) pairs to use as inputs to the places API. Initially, this is just based on the capital cities of each country, as well as cities with a population larger than 500 000.
- `./fetch_embassies` - Reads (lat,lng) pairs from `stdin` and fetches as many embassies within 50km of each of these as the Places API will allow. Outputs the Places API placeid for each of these to `stdout`. Note that this might return duplicates, you might want to pipe the output through `| sort | uniq`.
- `./fetch_reviews` - Takes a list of Place IDs and fetches reviews for all of them. Outputs reviews to JSON, along with the place name and a URL that links to Maps.
- `./process_reviews` - Filters and formats the reviews for Twitter. The format is `"★★★ [Review Text] [Maps URL]"`. Input is JSON from `./fetch_reviews`, Output is a pretty-printed JSON array of strings. Does some clever things to shorten what would be > 140 character tweets down to the limit.
- `./tweeter/` - A python thingy that tweets random entries from the included JSON file. The code here is super reusable - just dump new tweets in `tweetsrc.json`, and oauth credentials in `oauth.json` (see `ouath.json.tmpl` for a template), build and deploy to lambda.
- `./build_lambda.sh` - builds `./tweeter/` into a zip file that can be uploaded to AWS lambda for periodic tweets.
