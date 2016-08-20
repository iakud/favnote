package main

import (
	"bufio"
	"bytes"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

import (
	"github.com/fsnotify/fsnotify"
)

// FIXME : watcher done
var done = make(chan bool)

func loadAndWatchingBlogs() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}
	go startWatchingBlogs(watcher)
	loadBlogs()
}

func onEvent(event fsnotify.Event) {
	log.Println("event:", event)
	loadBlogs() // reload blogs
}

func startWatchingBlogs(watcher *fsnotify.Watcher) {
	err := watcher.Add("blog")
	if err != nil {
		log.Fatalln(err)
	}
	for {
		select {
		case event := <-watcher.Events:
			onEvent(event)
		case err := <-watcher.Errors:
			log.Println(err.Error())
		case <-done:
			break
		}
	}
}

// sort
type BlogSlice []*Blog

func (bs BlogSlice) Len() int {
	return len(bs)
}

func (bs BlogSlice) Swap(i, j int) {
	bs[i], bs[j] = bs[j], bs[i]
}

func (bs BlogSlice) Less(i, j int) bool {
	return bs[i].Date.After(bs[j].Date)
}

var (
	allBlogs  = make(map[string]*Blog)
	listBlogs = make(BlogSlice, 0)
)

func loadBlogs() {
	filePaths, err := filepath.Glob("blog/*.md")
	if err != nil {
		log.Fatalln(err.Error())
	}
	newAllBlogs := make(map[string]*Blog)
	newListBlogs := make(BlogSlice, 0)
	for _, filePath := range filePaths {
		fileName := filepath.Base(filePath)
		name := strings.TrimSuffix(fileName, ".md")
		blog := NewBlog(name)
		if err = blog.ParseFile(filePath); err != nil {
			log.Fatalln(err.Error())
		}
		newAllBlogs[name] = blog
		newListBlogs = append(newListBlogs, blog)
	}
	sort.Sort(newListBlogs)
	allBlogs = newAllBlogs
	listBlogs = newListBlogs
}

var (
	leftComment  = []byte("<!--")
	rightComment = []byte("-->")

	summaryRegexp = regexp.MustCompile(`^[\t ]*more[\t ]*$`)
	commentRegexp = regexp.MustCompile(`^[\t ]*([0-9A-Za-z_]+?)[\t ]*[:=][\t ]*(.*?)[\t ]*$`)
)

type Blog struct {
	Name   string
	Author string
	Date   time.Time
	Title  string
	//	tags     string
	//	category string
	//	status   string
	Summary template.HTML
	Content template.HTML
}

func NewBlog(name string) *Blog {
	return &Blog{Name: name}
}

func (b *Blog) ParseFile(filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	b.parse(buf)
	return nil
}

func (b *Blog) parse(text []byte) {
	var content bytes.Buffer
	for {
		leftPos := bytes.Index(text, leftComment)
		if leftPos < 0 {
			content.Write(text)
			break
		}
		rightPos := bytes.Index(text[leftPos+len(leftComment):], rightComment)
		if rightPos < 0 {
			content.Write(text)
			break
		}
		content.Write(text[:leftPos])
		comment := text[leftPos+len(leftComment) : leftPos+len(leftComment)+rightPos]
		if summaryRegexp.Match(comment) {
			b.Summary = template.HTML(blackfriday.MarkdownBasic(content.Bytes()))
		} else {
			b.parseComment(comment)
		}
		text = text[leftPos+len(leftComment)+rightPos+len(rightComment):]
	}
	b.Content = template.HTML(blackfriday.MarkdownBasic(content.Bytes()))
}

func (b *Blog) parseComment(comment []byte) {
	scanner := bufio.NewScanner(bytes.NewReader(comment))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if slice := commentRegexp.FindStringSubmatch(scanner.Text()); slice != nil {
			switch slice[1] {
			case "author":
				b.Author = slice[2]
			case "date":
				if date, err := time.Parse("2006-01-02 15:04:05", slice[2]); err != nil {
					log.Println(err.Error())
				} else {
					b.Date = date
				}
			case "title":
				b.Title = slice[2]
			default:
			}
		}
	}
}
