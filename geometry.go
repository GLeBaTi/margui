package margui

import (
	"encoding/xml"
	"fmt"
)

// XY describes a generic X, Y coordinate relative to a parent Canvas
// or CanvasObject.
type XY struct {
	X float32 // The position from the parent' left edge
	Y float32 // The position from the parent's top edge
}

// Add returns a new XY that is the result of offsetting the current
// position by p2 X and Y.
func (p1 XY) Add(p2 XY) XY {
	return XY{p1.X + p2.X, p1.Y + p2.Y}
}

// Subtract returns a new XY that is the result of offsetting the current
// position by p2 -X and -Y.
func (p1 XY) Subtract(p2 XY) XY {
	return XY{p1.X - p2.X, p1.Y - p2.Y}
}

// NewPos returns a newly allocated XY representing the specified coordinates.
func NewPosXY(x float32, y float32) XY {
	return XY{x, y}
}

func (c *XY) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: fmt.Sprintf("%.2f", c.X) + " " + fmt.Sprintf("%.2f", c.Y)}, nil
}

func (c *XY) UnmarshalXMLAttr(attr xml.Attr) error {
	_, err := fmt.Sscanf(attr.Value, "%f %f", &c.X, &c.Y)
	if err != nil {
		c.X = 0
		c.Y = 0
	}
	return nil
}

// XYZW
type XYZW struct {
	X float32
	Y float32
	Z float32
	W float32
}

// Add returns a new XY that is the result of offsetting the current
// position by p2 X and Y.
func (p1 XYZW) Add(p2 XYZW) XYZW {
	return XYZW{p1.X + p2.X, p1.Y + p2.Y, p1.Z + p2.Z, p1.W + p2.W}
}

// Subtract returns a new XYZW that is the result of offsetting the current
// position by p2 -X and -Y.
func (p1 XYZW) Subtract(p2 XYZW) XYZW {
	return XYZW{p1.X - p2.X, p1.Y - p2.Y, p1.Z - p2.Z, p1.W - p2.W}
}

// NewPos returns a newly allocated XYZW representing the specified coordinates.
func NewPosXYZW(x float32, y float32, z float32, w float32) XYZW {
	return XYZW{x, y, z, w}
}

func (c *XYZW) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: fmt.Sprintf("%.2f", c.X) + " " + fmt.Sprintf("%.2f", c.Y) + " " + fmt.Sprintf("%.2f", c.Z) + " " + fmt.Sprintf("%.2f", c.W)}, nil
}

func (c *XYZW) UnmarshalXMLAttr(attr xml.Attr) error {
	_, err := fmt.Sscanf(attr.Value, "%f %f %f %f", &c.X, &c.Y, &c.Z, &c.W)
	if err != nil {
		LogE(err)
		c.X = 0
		c.Y = 0
		c.Z = 0
		c.W = 0
	}
	return nil
}

// XYZ
type XYZ struct {
	X float32
	Y float32
	Z float32
}

// Add returns a new XY that is the result of offsetting the current
// position by p2 X and Y.
func (p1 XYZ) Add(p2 XYZ) XYZ {
	return XYZ{p1.X + p2.X, p1.Y + p2.Y, p1.Z + p2.Z}
}

// Subtract returns a new XYZ that is the result of offsetting the current
// position by p2 -X and -Y.
func (p1 XYZ) Subtract(p2 XYZ) XYZ {
	return XYZ{p1.X - p2.X, p1.Y - p2.Y, p1.Z - p2.Z}
}

// NewPos returns a newly allocated XYZ representing the specified coordinates.
func NewPosXYZ(x float32, y float32, z float32, w float32) XYZ {
	return XYZ{x, y, z}
}

func (c *XYZ) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: fmt.Sprintf("%.2f", c.X) + " " + fmt.Sprintf("%.2f", c.Y) + " " + fmt.Sprintf("%.2f", c.Z)}, nil
}

func (c *XYZ) UnmarshalXMLAttr(attr xml.Attr) error {
	_, err := fmt.Sscanf(attr.Value, "%f %f %f %f", &c.X, &c.Y, &c.Z)
	if err != nil {
		LogE(err)
		c.X = 0
		c.Y = 0
		c.Z = 0
	}
	return nil
}

type DockStyle string

const (
	//None = Center
	None   DockStyle = ""
	Center DockStyle = "Center"
	Fill   DockStyle = "Fill"

	FillLeft       DockStyle = "FillLeft"
	FillTop        DockStyle = "FillTop"
	FillRight      DockStyle = "FillRight"
	FillBottom     DockStyle = "FillBottom"
	FillHorizontal DockStyle = "FillHorizontal"
	FillVertical   DockStyle = "FillVertical"
	Left           DockStyle = "Left"
	Right          DockStyle = "Right"
	Top            DockStyle = "Top"
	Bottom         DockStyle = "Bottom"
	LeftTop        DockStyle = "LeftTop"
	LeftBottom     DockStyle = "LeftBottom"
	RightTop       DockStyle = "RightTop"
	RightBottom    DockStyle = "RightBottom"
)

type BorderStyle string

const (
	BNone   BorderStyle = ""
	BSolid  BorderStyle = "Solid"
	BDotted BorderStyle = "Dotted"
	BDashed BorderStyle = "Dashed"
)
