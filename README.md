# wiwi
Twitterの指定ハッシュタグのツイートと、Twitchの指定チャンネルのコメントを同時に流すくん

![image](https://user-images.githubusercontent.com/5152601/80907805-d182c900-8d54-11ea-817a-f1fd5cca2af0.png)

## Usage
```
$ wiwi -c #mogra -h #MU2020
```

## Installation
```
$ go get github.com/alitaso345/wiwi
```

### Required envs

```
export TWITCH_PASSWORD=YOUR_TWITCH_PASSORD
export TWITCH_NICK=YOUR_TWITCH_NICK

export TWITTER_CONSUMER_KEY=YOUR_TWITTER_CONSUMER_KEY
export TWITTER_CONSUMER_SECRET=YOUR_TWITTER_CONSUMER_SECRET_KEY
export TWITTER_ACCESS_TOKEN=YOUR_TWITTER_ACCESS_TOKEN
export TWITTER_ACCESS_TOKEN_SECRET=YOUR_TWITTER_ACCESS_TOKEN_SECRET
```

Also see https://dev.twitch.tv/docs/ and https://developer.twitter.com/en/docs

### Options

| Option | Description | Usage |
|--------|-------------|-------|
| -c     | twitch channel name | -c #mogra (ex: https://www.twitch.tv/mogra) |
| -h     | twitter hash tag name | -h #MU2020 |

## Credits
inspired by @9c5s and @yuki_sonogenic
