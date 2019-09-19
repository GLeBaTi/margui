package margui

import (
	"encoding/xml"
	"fmt"
)

// Size describes something with width and height.
type Size struct {
	Width  float32 // The number of units along the X axis.
	Height float32 // The number of units along the Y axis.
}

// Add returns a new Size that is the result of increasing the current size by
// s2 Width and Height.
func (s1 Size) Add(s2 Size) Size {
	return Size{s1.Width + s2.Width, s1.Height + s2.Height}
}

// Subtract returns a new Size that is the result of decreasing the current size
// by s2 Width and Height.
func (s1 Size) Subtract(s2 Size) Size {
	return Size{s1.Width - s2.Width, s1.Height - s2.Height}
}

// Union returns a new Size that is the maximum of the current Size and s2.
func (s1 Size) Union(s2 Size) Size {
	maxW := MaxF32(s1.Width, s2.Width)
	maxH := MaxF32(s1.Height, s2.Height)

	return NewSize(maxW, maxH)
}

// NewSize returns a newly allocated Size of the specified dimensions.
func NewSize(w float32, h float32) Size {
	return Size{w, h}
}

func (c *Size) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: fmt.Sprintf("%.2f", c.Width) + " " + fmt.Sprintf("%.2f", c.Height)}, nil
}

func (c *Size) UnmarshalXMLAttr(attr xml.Attr) error {
	_, err := fmt.Sscanf(attr.Value, "%f %f", &c.Width, &c.Height)
	if err != nil {
		c.Width = 0
		c.Height = 0
	}
	return nil
}

// Position describes a generic X, Y coordinate relative to a parent Canvas
// or CanvasObject.
type Position struct {
	X float32 // The position from the parent' left edge
	Y float32 // The position from the parent's top edge
}

// Add returns a new Position that is the result of offsetting the current
// position by p2 X and Y.
func (p1 Position) Add(p2 Position) Position {
	return Position{p1.X + p2.X, p1.Y + p2.Y}
}

// Subtract returns a new Position that is the result of offsetting the current
// position by p2 -X and -Y.
func (p1 Position) Subtract(p2 Position) Position {
	return Position{p1.X - p2.X, p1.Y - p2.Y}
}

// NewPos returns a newly allocated Position representing the specified coordinates.
func NewPos(x float32, y float32) Position {
	return Position{x, y}
}

func (c *Position) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: fmt.Sprintf("%.2f", c.X) + " " + fmt.Sprintf("%.2f", c.Y)}, nil
}

func (c *Position) UnmarshalXMLAttr(attr xml.Attr) error {
	_, err := fmt.Sscanf(attr.Value, "%f %f", &c.X, &c.Y)
	if err != nil {
		c.X = 0
		c.Y = 0
	}
	return nil
}
