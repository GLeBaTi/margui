package comp

import (
	"github.com/GLeBaTi/margui"
)

type Control struct {
	Text string `xml:",chardata"`
	Id   string `xml:"Id,attr"`

	Margin       margui.XYZW `xml:"Margin,attr"`
	Padding      margui.XYZW `xml:"Padding,attr"`
	GlobalMargin margui.XYZW `xml:"-"`
	Pivot        margui.XY   `xml:"Pivot,attr"`
	Rotation     margui.XYZ  `xml:"Rotation,attr"`
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
