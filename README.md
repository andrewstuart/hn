##About

A hackernews ncurses CLI reader written in Go

Right now it's able to view articles, view comments, and open a page in your default browser, all done directly from the site using goquery (jquery-like library for Go), goncurses, and xdg-open for opening pages.

It's really not much, but I like it because I can view and interact with things a bit faster in the terminal.

It's also got a clunky REST API that can be started via the -s option. I'd intended to do some caching so I could expose an API publicly, but shortly after that, HN came out with their own API. I don't really plan to do much with the API in the future.


##Usage

###Story view
- n) Go to next page
- p) Go to previous page
- <num>c) View comments for story <num>
- <num>o) Open story <num> in default browser
- q) Quit hn

###Comments view
- d) Go down 30 lines
- u) Go up 30 lines
- j) Go down 1 line
- k) Go up 1 line
- n) Go down 1 page
- p) Go up 1 page
- q) Go back to story view
