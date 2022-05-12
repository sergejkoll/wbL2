package main

import "fmt"

type builder interface {
	makeTitle(title string) string
	makeString(str string) string
	makeItems(items []string) string
	close() string
}

type Director struct {
	builder builder
}

func (d *Director) Construct() string {
	result := d.builder.makeTitle("Title")
	result += d.builder.makeString("String")
	result += d.builder.makeItems([]string{
		"Item1",
		"Item2",
	})
	result += d.builder.close()
	return result
}

type TextBuilder struct {
}

func (*TextBuilder) makeTitle(title string) string {
	return "# " + title + "\n"
}

func (*TextBuilder) makeString(str string) string {
	return "## " + str + "\n"
}

func (*TextBuilder) makeItems(items []string) string {
	var result string
	for _, item := range items {
		result += "- " + item + "\n"
	}
	return result
}

func (*TextBuilder) close() string {
	return "\n"
}

func main() {
	d := Director{&TextBuilder{}}
	result := d.Construct()
	fmt.Println(result)
}
