package main

import (
	"bufio"
	"bytes"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"log"
	"regexp"
	"time"
)

var (
	leftComment  = []byte("<!--")
	rightComment = []byte("-->")

	summaryRegexp = regexp.MustCompile(`^[\t ]*more[\t ]*$`)
	commentRegexp = regexp.MustCompile(`^[\t ]*([0-9A-Za-z_]+?)[\t ]*[:=][\t ]*(.*?)[\t ]*$`)
)

type Article struct {
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

func NewArticle(name string) *Article {
	return &Article{Name: name}
}

func (article *Article) ParseFile(filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	article.parse(b)
	return nil
}

func (article *Article) parse(text []byte) {
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
			article.Summary = template.HTML(blackfriday.MarkdownBasic(content.Bytes()))
		} else {
			scanner := bufio.NewScanner(bytes.NewReader(comment))
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				if slice := commentRegexp.FindStringSubmatch(scanner.Text()); slice != nil {
					switch slice[1] {
					case "author":
						article.Author = slice[2]
					case "date":
						if date, err := time.Parse("2006-01-02 15:04:05", slice[2]); err != nil {
							log.Println(err.Error())
						} else {
							article.Date = date
						}
					case "title":
						article.Title = slice[2]
					default:
					}
				}
			}
		}
		text = text[leftPos+len(leftComment)+rightPos+len(rightComment):]
	}
	article.Content = template.HTML(blackfriday.MarkdownBasic(content.Bytes()))
}
