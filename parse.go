package main

import (
    "io"

    "github.com/antchfx/xmlquery"
)

func parseArticles(body io.Reader) ([]article, error) {

    root, err := xmlquery.Parse(body)
    if err != nil {
        return nil, err
    }

    articles := make([]article, 0)

    for _, item := range xmlquery.Find(root, "//item") {
        var article article

        article.title = xmlquery.FindOne(item, "/title").InnerText()
        article.link = xmlquery.FindOne(item, "/link").InnerText()

        articles = append(articles, article)
    }

    return articles, nil
}