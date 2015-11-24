package main

import (
	"fmt"
	"net/http"
	_ "strconv"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// TODO: Implement a score field for Post
type Post struct {
	author    string
	title     string
	subreddit string
	url       string
	comments  []Comment
}

type Comment struct {
	author string
	text   string
}

func main() {
	fmt.Println("starting...")

	resp, err := http.Get("https://www.reddit.com")
	if err != nil {
		panic(err)
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}

	matcher := func(n *html.Node) bool {
		if n.DataAtom == atom.Div && n.Parent != nil {
			return scrape.Attr(n, "id") == "siteTable"
		}
		return false
	}
	table, ok := scrape.Find(root, matcher)
	if !ok {
		panic(ok)
	}
	matcher = func(n *html.Node) bool {
		if n.DataAtom == atom.Div && n.Parent != nil {
			return scrape.Attr(n, "data-type") == "link"
		}
		return false
	}

	articles := scrape.FindAll(table, matcher)
	if ok {
		//var posts []Post
		for i := 0; i < len(articles); i++ {
			post := parsepost(articles[i])
			fmt.Println(post)
		}
	}

}

func getcomments(n *html.Node, post *Post) {
	// pass
}

func parsepost(n *html.Node) Post {
	post := Post{}

	// get the author. uses a scrape inbuilt matcher
	auth, _ := scrape.Find(n, scrape.ByClass("title"))
	author := scrape.Text(auth.FirstChild)

	// get the subreddit. This requires a custom matcher.
	matcher := func(n *html.Node) bool {
		if n.DataAtom == atom.A && n.Parent != nil {
			return scrape.Attr(n, "class") == "subreddit hover may-blank"
		}
		return false
	}
	sub, _ := scrape.Find(n, matcher)
	subreddit := scrape.Text(sub)

	// get the url to the comments. requires custom matcher.
	matcher = func(n *html.Node) bool {
		if n.DataAtom == atom.Ul && n.FirstChild != nil {
			return scrape.Attr(n, "class") == "flat-list buttons" && scrape.Attr(n.FirstChild, "class") == "first"
		}
		return false
	}
	ul, _ := scrape.Find(n, matcher)          // ul is a list of two buttons: one that links to a post's comments page, one a "share" function
	li := ul.FirstChild                       // the first list item of ul -- this will always be the comments page link.
	url := scrape.Attr(li.FirstChild, "href") // finally, the url found in the list item.

	post.author = author
	post.subreddit = subreddit
	post.url = url
	return post
}
func parsecomment(n *html.Node) Comment {
	// pass
	return Comment{}
}
