package icons

import (
	"github.com/a-h/templ"
	"strings"
)

func cn(class string, attrs []templ.Attributes) string {
	var classes []string
	classes = append(classes, class)
	for _, attr := range attrs {
		cl, ok := attr["class"]
		if !ok {
			continue
		}
		c, ok := cl.(string)
		if !ok {
			continue
		}
		classes = append(classes, c)
	}
	return strings.Join(classes, " ")
}

func at(attrs []templ.Attributes) templ.Attributes {
	attr := templ.Attributes{}
	for _, attr := range attrs {
		for key, value := range attr {
			attr[key] = value
		}
	}
	return attr
}
