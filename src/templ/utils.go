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

func hasAttr(attrs templ.Attributes, key string) bool {
	if attrs == nil {
		return false
	}
	_, ok := attrs[key]
	return ok
}

func getAttrStrOr(attrs templ.Attributes, key string, defaultValue string) string {
	if hasAttr(attrs, key) {
		if val, ok := attrs[key].(string); ok {
			return val
		}
	}
	return defaultValue
}

func mergeRightAttrs(attrs ...templ.Attributes) templ.Attributes {
	attr := templ.Attributes{}
	for _, attrParam := range attrs {
		for key, value := range attrParam {
			attr[key] = value
		}
	}
	return attr
}
func at(attrs []templ.Attributes) templ.Attributes {
	attr := mergeRightAttrs(attrs...)
	if hasAttr(attr, "class") {
		delete(attr, "class")
	}
	size := getAttrStrOr(attr, "size", "")
	if hasAttr(attr, "size") {
		delete(attr, "size")
	}
	strokeWidth := getAttrStrOr(attr, "stroke-width", defaultStrokeWidth)
	if hasAttr(attr, "stroke-width") {
		delete(attr, "stroke-width")
	}
	absoluteStrokeWidth := false
	if hasAttr(attr, "absoluteStrokeWidth") {
		absoluteStrokeWidth = asBool(attr["absoluteStrokeWidth"])
		delete(attr, "absoluteStrokeWidth")
	} else if hasAttr(attr, "absolute-stroke-width") {
		absoluteStrokeWidth = asBool(attr["absolute-stroke-width"])
		delete(attr, "absolute-stroke-width")
	}
	color := getAttrStrOr(attr, "color", "")
	if hasAttr(attr, "color") {
		delete(attr, "color")
	}

	finalAttr := templ.Attributes{
		"width":        defaultWidth,
		"height":       defaultHeight,
		"stroke":       defaultStroke,
		"stroke-width": strokeWidth,
	}
	if size != "" {
		finalAttr["width"] = size
		finalAttr["height"] = size
	}
	if color != "" {
		finalAttr["stroke"] = color
	}
	strokeWidthFloat, errStrokeWidth := strconv.ParseFloat(strokeWidth, 64)
	sizeFloat, errSize := strconv.ParseFloat(size, 64)
	if absoluteStrokeWidth && errStrokeWidth == nil && errSize == nil {
		finalAttr["stroke-width"] = strconv.FormatFloat(strokeWidthFloat*24/sizeFloat, 'f', -1, 64)
	}

	return mergeRightAttrs(
		defaultAttributes,
		finalAttr,
		attr,
	)
}
