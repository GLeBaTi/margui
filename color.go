package margui

import (
	"encoding/xml"
	"fmt"
)

type Color struct {
	R float32
	G float32
	B float32
	A float32
}

func (c *Color) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: fmt.Sprintf("%.2f", c.R) + " " + fmt.Sprintf("%.2f", c.G) + " " + fmt.Sprintf("%.2f", c.B) + " " + fmt.Sprintf("%.2f", c.A)}, nil
}

func (c *Color) UnmarshalXMLAttr(attr xml.Attr) error {
	_, err := fmt.Sscanf(attr.Value, "%f %f %f %f", &c.R, &c.G, &c.B, &c.A)
	if err != nil {
		c.R = 0
		c.G = 0
		c.B = 0
		c.A = 0
	}
	return nil
}