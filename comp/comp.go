package comp

import (
	"github.com/GLeBaTi/margui"
)

type Control struct {
	Text string `xml:",chardata"`
	Id   string `xml:"Id,attr"`

	Position margui.Position `xml:"Position,attr"`
	Size     margui.Size     `xml:"Size,attr"`
	Color    margui.Color    `xml:"Color,attr"`

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
