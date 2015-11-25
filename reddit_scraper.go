package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

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
}

var wg sync.WaitGroup

func main() {

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
	var posts []Post

	for i := 0; i < len(articles); i++ {
		wg.Add(1)
		go func(n *html.Node) {
			post := parsepost(n)
			posts = append(posts, post)
			wg.Done()
		}(articles[i])
	}

	wg.Wait()

	for i := 0; i < len(posts); i++ {
		printpost(posts[i])
	}

}

// Basically for debugging, and because go prints structs about as well as a gopher speaks english.
func printpost(post Post) {
	fmt.Println("Title: ", post.title)
	fmt.Println("Author: ", post.author)
	fmt.Println("Subreddit: ", post.subreddit)
	fmt.Println("url: ", post.url)
}

func parsepost(n *html.Node) Post {
	post := Post{}

	// get the title. uses a scrape inbuilt matcher
	title_scrape, _ := scrape.Find(n, scrape.ByClass("title"))
	title := scrape.Text(title_scrape.FirstChild)

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

	// get the author. Uses custom matcher and magic.
	matcher = func(n *html.Node) bool {
		if n.DataAtom == atom.A && n.Parent.DataAtom == atom.P {
			return strings.Contains(scrape.Attr(n, "href"), "/user/")
		}
		return false
	}
	author_scrape, _ := scrape.Find(n, matcher)
	author := scrape.Text(author_scrape)

	post.title = title
	post.subreddit = subreddit
	post.url = url
	post.author = author

	return post
}
