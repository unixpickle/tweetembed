# tweetembed

This is a program for training [word embeddings](https://github.com/unixpickle/wordembed) on a large corpus of Tweets collected by [tweetdump](https://github.com/unixpickle/tweetdump). The resulting embeddings can be used for various NLP tasks.

# Results

The following results were from word embeddings trained on a billion-word twitter dump. The twitter dump was started on February 19, 2017 and went on for about a month.

I used 64-dimensional embeddings. I trained these embeddings overnight (about 10 hours), which performed 4.9e9 updates (16 iterations). The final cost was 0.057.

## Neighbors

The easiest way to look inside Twitter's collective mind is with nearest neighbors. Every row in the following table shows the four nearest neighbors of the given word:

|Word        |             |             |             |             |
|------------|-------------|-------------|-------------|-------------|
|cat         |dog          |kitty        |kitten       |cats         |
|video       |youtube      |videos       |liked        |watch        |
|red         |blue         |yellow       |green        |pink         |
|evil        |wicked       |gods         |cult         |fear         |
|nice        |pretty       |good         |cool         |look         |
|trump       |obama        |president    |donald       |trumps       |
|democrat    |republican   |democratic   |candidate    |democrats    |
|smile       |hug          |eyes         |smiling      |smiles       |
|august      |july         |september    |october      |20th         |
|three       |four         |two          |five         |six          |
|3           |2            |4            |1            |5            |
|turd        |slimy        |hag          |turds        |prick        |
|neighbor    |neighbour    |neighbors    |upstairs     |sister       |
|color       |colour       |yellow       |blue         |thin         |
|rap         |hip          |song         |hop          |rappers      |
|loud        |sound        |hear         |noise        |hard         |

## Correlations

Using embeddings, we can measure correlations between arbitrary pairs of words.

Let's look at how delicious various foods are:

|word        |correlation with "delicious" |
|------------|-----------------------------|
|chocolate   |0.706                        |
|pizza       |0.635                        |
|chicken     |0.635                        |
|candy       |0.624                        |
|burger      |0.616                        |
|tacos       |0.584                        |
|poop        |0.372                        |
|soylent     |0.242                        |

Twitter is a harsh critic, especially when it comes to Soylent.

We can also look at correlations between various political figures. Note that these correlations were generated directly from a large corpus of tweets without any of my own political biases.

|word        |correlation with "putin" |
|------------|-------------------------|
|trump       |0.774                    |
|clinton     |0.758                    |
|obama       |0.710                    |
|sanders     |0.511                    |
|bernie      |0.488                    |

With more recent Twitter data, I bet the numbers here would be *very* different.

We can also look at the correlations between days of the week:

|word        |correlation with "saturday" |
|------------|----------------------------|
|friday      |0.936                       |
|thursday    |0.935                       |
|sunday      |0.934                       |
|monday      |0.914                       |
|tuesday     |0.913                       |
|wednesday   |0.904                       |

Looks like Wednesday is as far from Saturday as you can get, at least in spirit.

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

Now you can use the `neighbors`, `correlation` and `analogy` sub-commands:

```
$ go run *.go neighbors -token red
red 0.9999999
blue 0.86534697
yellow 0.7989963
green 0.7967982
pink 0.7936595
$ go run *.go correlation dog cat
0.8497587
$ go run *.go analogy bill gates elon
elon 0.6260362
musk 0.5708943
kady 0.5354746
mavis 0.5264958
davidwalliams 0.49428728
```
