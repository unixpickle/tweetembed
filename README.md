# tweetembed

This is a program for training [word embeddings](https://github.com/unixpickle/wordembed) on a large corpus of Tweets collected by [tweetdump](https://github.com/unixpickle/tweetdump). The resulting embeddings could be used for various NLP tasks.

# Usage

First, install Go and tweetembed. Then run the following commands in whatever directory you want to save the results in. Replace `tweets.csv` with your Twitter data:

```shell
cat tweets.csv | tweetembed tokens
cat tweets.csv | tweetembed matrix
tweetembed train
tweetembed embed
```

You will need to stop `tweetembed train` manually by pressing Ctrl+C exactly once. Here is an example of the output you should see:

```shell
$ cat tweets.csv | tweetembed tokens
reading tweet CSV from standard input.
processed 76288 tweets (134332 tokens)
$ cat tweets.csv | tweetembed matrix
reading tweet CSV from standard input.
processed 76288 tweets
$ tweetembed train
Creating a new trainer...
2017/06/30 18:53:37 done 520 updates: cost=0.440643
2017/06/30 18:53:37 done 1040 updates: cost=0.368217
...
^C
Caught interrupt. Ctrl+C again to terminate.
2017/06/30 18:53:39 Saving result...
$ tweetembed embed
2017/06/30 18:54:13 Creating embedding...
2017/06/30 18:54:13 Saving result...
```
