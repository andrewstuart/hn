hn
==

A hackernews ncurses CLI reader written in Go

Right now it's able to view articles, view comments, and open a page in your default browser, all done directly from the site using goquery (jquery-like library for Go) and goncurses.

It's really not much, but I like it because I can view and interact with things a bit faster in the terminal.

It's also got a clunky REST API that can be started via the -s option. I'd intended to do some caching so I could expose an API publicly, but shortly after that, HN came out with their own API. I don't really plan to do much with the API in the future.
