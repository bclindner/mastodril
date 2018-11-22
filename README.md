# mastodril

A hacky Go script that syncs a Twitter user's Tweets to a Mastodon account.
This is primarily used to mirror [@dril](https://twitter.com/dril)'s tweets 
to a [Mastodon account](https://botsin.space/@mastodril), hence the name.

## Configuration

A valid config file is stored in `mastodril.json` in the same folder as the
executable, and looks like this. You'll need app credentials for Twitter and
Mastodon.
```json
{
	"handle": "<handle_of_user_to_sync>",
	"twitter": {
		"consumerKey": "<your_twitter_consumer_key>",
		"consumerSecret": "<your_twitter_consumer_secret>",
		"accessToken": "<your_twitter_access_token>",
		"accessSecret": "<your_twitter_access_secret>"
	},
	"mastodon": {
		"server": "<url_to_your_mastodon_instance>",
		"clientID": "<your_mastodon_client_id>",
		"clientSecret": "<your_mastodon_client_secret>",
		"accessToken": "<your_mastodon_access_token>"
	}
}
```

## Author's Note

The bot does not run as a daemon - it's a oneshot process that syncs the
last 5 tweets, then takes the last tweet in that list and stores its ID
to a file. It then pulls up that file on the next run and tries to find at
most 5 tweets starting from that saved point.

As such, if the user you're trying to mirror posts a lot, this may not be 
a feasible script to use. Theoretically Twitter's streaming API would work
well there, but that comes with its own set of pitfalls (ratelimiting,
filtering problems, etc). This is a simple enough solution for syncing 
slower timelines.
