package comp

import (
	"math/rand"

	"github.com/GLeBaTi/margui"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Interactable interface {
	KeyEventHandler(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey)
	CursorEventHandler(x, y float64)
	CursorButtonEventHandler(button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey)
}

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
	Positionable

	Text           string `xml:",chardata"`
	Id             string `xml:"Id,attr"`
	IsInteractable bool   `xml:"IsInteractable,attr"`

	Rectangles []*Rectangle `xml:"Rectangle"`
	Ellipses   []*Ellipse   `xml:"Ellipse"`
	Paths      []*Rectangle `xml:"Path"`
	Polygons   []*Rectangle `xml:"Polygon"`

	IsMouseInside bool
	MouseButton   glfw.MouseButton
	MouseAction   glfw.Action
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

func (e *Rectangle) KeyEventHandler(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {

}

func (e *Rectangle) CursorEventHandler(mouseX, mouseY float64) {
	e.IsMouseInside = (mouseX >= float64(e.GlobalMargin.X)) && (mouseX <= float64(e.GlobalMargin.X+e.GlobalMargin.Z)) &&
		(mouseY >= float64(e.GlobalMargin.Y)) && (mouseY <= float64(e.GlobalMargin.Y+e.GlobalMargin.W))
}

func (e *Rectangle) CursorButtonEventHandler(button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	e.MouseButton = button
	e.MouseAction = action
	if e.IsMouseInside && e.MouseButton == glfw.MouseButtonLeft && e.MouseAction == glfw.Press {
		e.Color = margui.NewColor(rand.Float32(), rand.Float32(), rand.Float32(), 1)
	}
}
