package comp

import (
	"github.com/GLeBaTi/margui"
)

type Control struct {
	Text string `xml:",chardata"`
	Id   string `xml:"Id,attr"`

	//Margin X:Left Y:Top Z:Right(Width) W:Bottom(Height)
	Margin margui.XYZW `xml:"Margin,attr"`
	//Margin X:Left Y:Top Z:Width W:Height
	GlobalMargin margui.XYZW      `xml:"-"`
	Pivot        margui.XY        `xml:"Pivot,attr"`
	Rotation     margui.XYZ       `xml:"Rotation,attr"`
	Dock         margui.DockStyle `xml:"Dock,attr"`
	//Scale          margui.XY `xml:"Scale,attr"`
	//GlobalScale    margui.XY `xml:"-"`
	Color margui.Color `xml:"Color,attr"`

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
