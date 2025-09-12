<img src="https://imgs.xkcd.com/comics/lisp.jpg" />

My implementation of the Exercise 4.12 from the https://gopl.io: 

The popular web comic *xkcd* has a JSON interface. For example, a request to https://xkcd.com/571/info.0.json produces a detailed description of comic 571, one of many favorites. Download each URL (once!) and build an offline index. Write a tool `xkcd` that, using this index, prints the URL and transcript of each comic that matches a search term provided on the command line.

```
$ go install ./cmd/xkcd.go

$ xkcd -h
Usage of xkcd:
  -c int
        max number of concurrent http requests when building offline index (default 20)
  -i string
        file holding offline index of comics (default "xkcd.json")
  -t    print also the transcript

$ xkcd perl
# 181 (2006) Interblag                           https://xkcd.com/181/
# 208 (2007) Regular Expressions                 https://xkcd.com/208/
# 224 (2007) Lisp                                https://xkcd.com/224/
# 312 (2007) With Apologies to Robert Frost      https://xkcd.com/312/
# 353 (2007) Python                              https://xkcd.com/353/
# 407 (2008) Cheap GPS                           https://xkcd.com/407/
# 519 (2008) 11th Grade                          https://xkcd.com/519/
# 621 (2009) Superlative                         https://xkcd.com/621/
#1149 (2012) Broomstick                          https://xkcd.com/1149/
#1171 (2013) Perl Problems                       https://xkcd.com/1171/
#1286 (2013) Encryptic                           https://xkcd.com/1286/
#1306 (2013) Sigil Cycle                         https://xkcd.com/1306/
#1599 (2015) Water Delivery                      https://xkcd.com/1599/
```
