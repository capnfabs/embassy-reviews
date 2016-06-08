#!/usr/bin/env python
import json
import random
import tweepy

def tweet(event, context):
    # Read credentials
    f = open('oauth.json', 'r')
    auth_info = json.loads(f.read())
    f.close()

    # Set up auth
    auth = tweepy.OAuthHandler(auth_info['consumer_key'], auth_info['consumer_secret'])
    auth.secure = True
    auth.set_access_token(auth_info['access_token'], auth_info['access_token_secret'])

    # Actually authenticate
    api = tweepy.API(auth)

    # Pick a tweet.
    f = open('tweetsrc.json')
    tweets = json.loads(f.read())
    f.close()
    tweet = random.choice(tweets)

    # Actually tweet.
    api.update_status(status=tweet)


if __name__ == "__main__":
    tweet(None, None)
