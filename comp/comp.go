package comp

import (
	"github.com/GLeBaTi/margui"
)

type Control struct {
	Text string `xml:",chardata"`
	Id   string `xml:"Id,attr"`

	BackgroundColor          *ColorBrush          `xml:"BackgroundColor"`
	BackgroundLinearGradient *LinearGradientBrush `xml:"BackgroundLinearGradient"`
	BackgroundRadialGradient *RadialGradientBrush `xml:"BackgroundRadialGradient"`
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
	Color *margui.Color `xml:"Color,attr"`

	Buttons []*Button `xml:"Button"`
	Panels  []*Panel  `xml:"Panel"`
}

type Window struct {
	Control
}

type Panel struct {
	Control
}

type Button struct {
	Control
}

type Brush interface{}

type ColorBrush struct {
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
