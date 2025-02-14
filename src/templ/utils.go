package icons

import (
	"sort"
	"strconv"
	"strings"

	"github.com/a-h/templ"
)

// Returns true if the value is "true" or true
func asBool(value any) bool {
	if value == nil {
		return false
	}
	bStr, ok := value.(string)
	if ok {
		return bStr == "true"
	}
	b, ok := value.(bool)
	if !ok {
		return false
	}
	return b
}

func cn(class string, attrs []templ.Attributes) string {
	classes := make(map[string]struct{})
	for _, cl := range strings.Split(class, " ") {
		classes[cl] = struct{}{}
	}
	for _, attr := range attrs {
		if classAttrib, ok := attr["class"]; ok {
			if c, ok := classAttrib.(string); ok {
				for _, cl := range strings.Split(c, " ") {
					classes[cl] = struct{}{}
				}
			}
		}
	}
	uniqClasses := []string{}
	for cl := range classes {
		uniqClasses = append(uniqClasses, cl)
	}
	sort.Strings(uniqClasses)
	return strings.Join(uniqClasses, " ")
}

func at(attrs []templ.Attributes) templ.Attributes {
	attr := templ.Attributes{}
	for _, attrParam := range attrs {
		for key, value := range attrParam {
			if key == "class" {
				continue
			}
			attr[key] = value
		}
	}
	size := "24"
	if _, ok := attr["size"]; ok {
		size = attr["size"].(string)
		if _, ok := attr["width"]; !ok {
			attr["width"] = attr["size"]
		}
		if _, ok := attr["height"]; !ok {
			attr["height"] = attr["size"]
		}
		delete(attr, "size")
	}
	if _, ok := attr["color"]; ok {
		if _, ok := attr["stroke"]; !ok {
			attr["stroke"] = attr["color"]
		}
		delete(attr, "color")
	}

	strokeWidth := "2"
	if _, ok := attr["stroke-width"]; ok {
		strokeWidth = attr["stroke-width"].(string)
	}
	absoluteStrokeWidth := false
	if val, ok := attr["absoluteStrokeWidth"]; ok {
		absoluteStrokeWidth = asBool(val)
		delete(attr, "absoluteStrokeWidth")
	} else if val, ok := attr["absolute-stroke-width"]; ok {
		absoluteStrokeWidth = asBool(val)
		delete(attr, "absolute-stroke-width")
	}
	strokeWidthFloat, errStrokeWidth := strconv.ParseFloat(strokeWidth, 64)
	sizeFloat, errSize := strconv.ParseFloat(size, 64)
	if absoluteStrokeWidth && errStrokeWidth == nil && errSize == nil {
		strokeWidth = strconv.FormatFloat(strokeWidthFloat*24/sizeFloat, 'f', -1, 64)
	}
	attr["stroke-width"] = strokeWidth
	for key, value := range defaultAttributes {
		if _, ok := attr[key]; !ok {
			attr[key] = value
		}
	}
	return attr
}
