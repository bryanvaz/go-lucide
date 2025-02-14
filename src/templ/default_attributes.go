package icons

import "github.com/a-h/templ"

const (
	defaultXmlns          = "http://www.w3.org/2000/svg"
	defaultWidth          = "24"
	defaultHeight         = "24"
	defaultViewBox        = "0 0 24 24"
	defaultFill           = "none"
	defaultStroke         = "currentColor"
	defaultStrokeWidth    = "2"
	defaultStrokeLinecap  = "round"
	defaultStrokeLinejoin = "round"
)

var defaultAttributes = templ.Attributes{
	"xmlns":           defaultXmlns,
	"width":           defaultWidth,
	"height":          defaultHeight,
	"viewBox":         defaultViewBox,
	"fill":            defaultFill,
	"stroke":          defaultStroke,
	"stroke-width":    defaultStrokeWidth,
	"stroke-linecap":  defaultStrokeLinecap,
	"stroke-linejoin": defaultStrokeLinejoin,
}
