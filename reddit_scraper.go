package main

import (
	"fmt"
	"net/http"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Post struct {
	author    string
	title     string
	text      string
	subreddit string
	comments  []Comment
}

type Comment struct {
	author string
	text   string
}

func main() {
	fmt.Println("Hello there. I'm a placeholder for when the main function serves no purpose.")

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
	author, ok := scrape.Find(n, scrape.ByClass("title"))
	if ok {
		post.author = scrape.Text(author.FirstChild)
	}

	return post
}
func parsecomment(n *html.Node) Comment {
	// pass
	return Comment{}
}
