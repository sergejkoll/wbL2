package main

import "fmt"

var db = map[string]string{
	"a@a.com": "a",
	"b@b.com": "b",
}

type database struct {
}

func (database) getNameByMail(mail string) string {
	return db[mail]
}

type mdWriter struct {
}

func (mdWriter) title(title string) string {
	return "# Welcome to " + title + "'s page!"
}

type PageMaker struct {
}

func (PageMaker) MakeWelcomePage(mail string) string {
	database := database{}
	writer := mdWriter{}

	name := database.getNameByMail(mail)
	page := writer.title(name)

	return page
}

func main() {
	pageMaker := PageMaker{}
	page := pageMaker.MakeWelcomePage("a@a.com")
	fmt.Println(page)

	page = pageMaker.MakeWelcomePage("b@b.com")
	fmt.Println(page)
}
