package main

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

// sort
type ArticleSlice []*Article

func (a ArticleSlice) Len() int {
	return len(a)
}

func (a ArticleSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ArticleSlice) Less(i, j int) bool {
	return a[i].Date.After(a[j].Date)
}

var (
	articles       = make(map[string]*Article)
	sortedArticles = make(ArticleSlice, 0)
)

func loadBlogs() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				loadArticles() // reload articles
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("blog")
	if err != nil {
		log.Fatal(err)
	}
	loadArticles()
}

func loadArticles() {
	filePaths, err := filepath.Glob("blog/*.md")
	if err != nil {
		log.Fatalln(err.Error())
	}
	newArticles := make(map[string]*Article)
	newSortedArticles := make(ArticleSlice, 0)
	for _, filePath := range filePaths {
		fileName := filepath.Base(filePath)
		name := strings.TrimSuffix(fileName, ".md")
		article := NewArticle(name)
		if err = article.ParseFile(filePath); err != nil {
			log.Fatalln(err.Error())
		}
		newArticles[name] = article
		newSortedArticles = append(newSortedArticles, article)
	}
	sort.Sort(newSortedArticles)
	articles = newArticles
	sortedArticles = newSortedArticles
}
