# bofhwits-feed [![Build Status](https://travis-ci.org/amauragis/bofhwits-feed.svg?branch=master)](https://travis-ci.org/amauragis/bofhwits-feed)

Supposedly go makes it really easy to make web services with just net.http, so I decided to try to *completely rip off*
[Comradephate's Unicorn/Sinatra version](https://github.com/Comradephate/bofhwits_blotter).  Honestly, there isn't a whole lot to it.
This takes posts in a mysql database created by [bofhwits](https://github.com/amauragis/bofhwits) and serves them to web clients.

It's a goal to add some filtering ability, a search, and potentially a "worst poster" type dealy with stats or something.  It's
probably pretty easy because it's all in a database.

Also, I don't really know enough about databases to say anything for sure, but issuing a query for every row on every request seems
wasteful and I should maybe consider some sort of "caching".  Maybe create an index?  Maybe google will tell me.  :shrug:

## Installing
1. [Install Go (>= 1.6)](https://golang.org/doc/install)
  - If you haven't installed go before, do read that page a bit.  $GOPATH and $GOROOT and such are kind of
    non intuitive.
1. `go get github.com/amauragis/bofhwits-feed`
1. Set up the configuration file to point to your bofhwits database.
