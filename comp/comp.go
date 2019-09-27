package comp

import (
	"github.com/GLeBaTi/margui"
)

type Interactable struct{}
type Positionable struct {
	//Margin X:Left Y:Top Z:Right(Width) W:Bottom(Height)
	Margin margui.XYZW `xml:"Margin,attr"`
	//Margin X:Left Y:Top Z:Width W:Height
	GlobalMargin margui.XYZW `xml:"-"`
	//Start of element coordinates [-1; 1]
	Pivot    margui.XY        `xml:"Pivot,attr"`
	Rotation margui.XYZ       `xml:"Rotation,attr"`
	Dock     margui.DockStyle `xml:"Dock,attr"`
	//Scale          margui.XY `xml:"Scale,attr"`
	//GlobalScale    margui.XY `xml:"-"`
}

type Control struct {
	Interactable
	Positionable

	Text string `xml:",chardata"`
	Id   string `xml:"Id,attr"`

	Rectangles []*Rectangle `xml:"Rectangle"`
	Ellipses   []*Ellipse   `xml:"Ellipse"`
	Paths      []*Rectangle `xml:"Path"`
	Polygons   []*Rectangle `xml:"Polygon"`
}

type SolidColorBrush struct {
	Color margui.Color `xml:"Color,attr"`
}

type LinearGradientBrush struct {
	//Start [0, 1] top left
	Start margui.XY `xml:"Start,attr"`
	//End [0, 1] top left
	End    margui.XY       `xml:"End,attr"`
	Colors []*GradientStop `xml:"GradientStop"`
}

type RadialGradientBrush struct {
	//Pivot [-1, 1] center
	Pivot  margui.XY       `xml:"Pivot,attr"`
	Colors []*GradientStop `xml:"GradientStop"`
}

type GradientStop struct {
	Color margui.Color `xml:"Color,attr"`
	//Offset [0, 1] top left
	Offset float32 `xml:"Offset,attr"`
}

type Geometry struct {
	Control

	BackgroundSolidColor     *SolidColorBrush     `xml:"BackgroundSolidColor"`
	BackgroundLinearGradient *LinearGradientBrush `xml:"BackgroundLinearGradient"`
	BackgroundRadialGradient *RadialGradientBrush `xml:"BackgroundRadialGradient"`
	BorderSolidColor         *SolidColorBrush     `xml:"BorderSolidColor"`
	BorderLinearGradient     *LinearGradientBrush `xml:"BorderLinearGradient"`
	BorderRadialGradient     *RadialGradientBrush `xml:"BorderRadialGradient"`

	BorderWidth float32            `xml:"BorderWidth,attr"`
	BorderStyle margui.BorderStyle `xml:"BorderStyle,attr"`

	Color       *margui.Color `xml:"Color,attr"`
	BorderColor *margui.Color `xml:"BorderColor,attr"`
}

type Rectangle struct {
	Geometry
}

type Ellipse struct {
	Geometry
}

type Path struct {
	Geometry
}

type Polygon struct {
	Geometry
}

type Drawable interface {
	Draw()
}
