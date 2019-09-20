package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/GLeBaTi/margui"
	"github.com/GLeBaTi/margui/comp"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	screenWidth  = 1000
	screenHeight = 1000
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func freeCoords(vao, vbo uint32) {
	gl.BindVertexArray(0)
	gl.DeleteVertexArrays(1, &vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.DeleteBuffers(1, &vbo)
}

func drawTexture(texture uint32, points []float32) {
	//gl.ActiveTexture(gl.TEXTURE0)
	//gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(len(points)/5))
}

func rectCoords(width float32, height float32, posX float32, posY float32) ([]float32, uint32, uint32) {

	xPos := float32(posX) / float32(screenWidth)
	x1 := -1 + xPos*2
	x2Pos := float32(posX+width) / float32(screenWidth)
	x2 := -1 + x2Pos*2

	yPos := float32(posY) / float32(screenHeight)
	y1 := 1 - yPos*2
	y2Pos := float32(posY+height) / float32(screenHeight)
	y2 := 1 - y2Pos*2

	points := []float32{
		// coord x, y, x texture x, y
		x1, y2, 0, 0.0, 1.0, // top left
		x1, y1, 0, 0.0, 0.0, // bottom left
		x2, y2, 0, 1.0, 1.0, // top right
		x2, y1, 0, 1.0, 0.0, // bottom right
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	return points, vao, vbo
}

func newTexture(color margui.Color) (texture uint32, program uint32) {
	//gl.GenTextures(1, &texture)
	//gl.ActiveTexture(gl.TEXTURE0)
	//gl.BindTexture(gl.TEXTURE_2D, texture)
	//gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	//gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	//gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	//gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	vertex_shaderStr := `
    #version 110
    attribute vec3 vert;
    void main() {
        gl_Position = vec4(vert, 1);
    }
` + "\x00"

	fragment_shaderStr := fmt.Sprintf(`
    #version 110
    out vec4 frag_colour;
    void main() {
        frag_colour = vec4(%f, %f, %f, %f);
    }
`+"\x00", color.R, color.G, color.B, color.A)

	vs, err := compileShader(vertex_shaderStr, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fs, err := compileShader(fragment_shaderStr, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program = gl.CreateProgram()
	gl.AttachShader(program, fs)
	gl.AttachShader(program, vs)
	gl.LinkProgram(program)

	return texture, program
}

func draw(win *comp.Window) {
	drawControl(nil, &win.Control)
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		info := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(info))

		return 0, fmt.Errorf("failed to compile %v: %v", source, info)
	}

	return shader, nil
}

func unmarshalFromFile() (*comp.Window, error) {

	xmlFile, err := os.Open("main.xml")
	if err != nil {
		return nil, err
	}

	defer margui.Close(xmlFile)

	fileData, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return nil, err
	}

	var out comp.Window
	err = xml.Unmarshal(fileData, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func calcGlobalMargin(parent *comp.Control, ctrl *comp.Control) margui.XYZW {
	var calculated margui.XYZW
	var parentGlobalMargin = margui.XYZW{
		X: 0,
		Y: 0,
		Z: screenWidth,
		W: screenHeight,
	}

	//TODO parent pivot topLeft by default (0,0 or -1,1???)

	if parent != nil {
		parentGlobalMargin = parent.GlobalMargin
	}

	if ctrl.Dock == margui.Fill {
		calculated.Z = parentGlobalMargin.Z - ctrl.Margin.X - ctrl.Margin.Z
		calculated.W = parentGlobalMargin.W - ctrl.Margin.Y - ctrl.Margin.W
		calculated.X = ctrl.Margin.X
		calculated.Y = ctrl.Margin.Y
	} else if ctrl.Dock == margui.None || ctrl.Dock == margui.Center {
		calculated.Z = ctrl.Margin.Z
		calculated.W = ctrl.Margin.W
		calculated.X = ctrl.Margin.X + parentGlobalMargin.Z/2.0
		calculated.Y = ctrl.Margin.Y + parentGlobalMargin.W/2.0
	} else if ctrl.Dock == margui.Left {
		calculated.Z = ctrl.Margin.Z
		calculated.W = ctrl.Margin.W
		calculated.X = ctrl.Margin.X
		calculated.Y = ctrl.Margin.Y + parentGlobalMargin.W/2.0
	} else if ctrl.Dock == margui.Top {
		calculated.Z = ctrl.Margin.Z
		calculated.W = ctrl.Margin.W
		calculated.X = ctrl.Margin.X + parentGlobalMargin.Z/2.0
		calculated.Y = ctrl.Margin.Y
	} else if ctrl.Dock == margui.Right {
		calculated.Z = ctrl.Margin.Z
		calculated.W = ctrl.Margin.W
		calculated.X = ctrl.Margin.X + parentGlobalMargin.Z
		calculated.Y = ctrl.Margin.Y + parentGlobalMargin.W/2.0
	} else if ctrl.Dock == margui.Bottom {
		calculated.Z = ctrl.Margin.Z
		calculated.W = ctrl.Margin.W
		calculated.X = ctrl.Margin.X + parentGlobalMargin.Z/2.0
		calculated.Y = ctrl.Margin.Y + parentGlobalMargin.W
	} else if ctrl.Dock == margui.LeftTop {
		calculated.Z = ctrl.Margin.Z
		calculated.W = ctrl.Margin.W
		calculated.X = ctrl.Margin.X
		calculated.Y = ctrl.Margin.Y
	} else if ctrl.Dock == margui.RightTop {
		calculated.Z = ctrl.Margin.Z
		calculated.W = ctrl.Margin.W
		calculated.X = ctrl.Margin.X + parentGlobalMargin.Z
		calculated.Y = ctrl.Margin.Y
	} else if ctrl.Dock == margui.RightBottom {
		calculated.Z = ctrl.Margin.Z
		calculated.W = ctrl.Margin.W
		calculated.X = ctrl.Margin.X + parentGlobalMargin.Z
		calculated.Y = ctrl.Margin.Y + parentGlobalMargin.W
	} else if ctrl.Dock == margui.LeftBottom {
		calculated.Z = ctrl.Margin.Z
		calculated.W = ctrl.Margin.W
		calculated.X = ctrl.Margin.X
		calculated.Y = ctrl.Margin.Y + parentGlobalMargin.W
	} else if ctrl.Dock == margui.FillHorizontal {
		calculated.Z = parentGlobalMargin.Z - ctrl.Margin.X - ctrl.Margin.Z
		calculated.W = ctrl.Margin.W
		calculated.X = ctrl.Margin.X
		calculated.Y = ctrl.Margin.Y + parentGlobalMargin.W/2.0
	} else if ctrl.Dock == margui.FillVertical {
		calculated.Z = ctrl.Margin.Z
		calculated.W = parentGlobalMargin.W - ctrl.Margin.Y - ctrl.Margin.W
		calculated.X = ctrl.Margin.X + parentGlobalMargin.Z/2.0
		calculated.Y = ctrl.Margin.Y
	} else if ctrl.Dock == margui.FillLeft {
		calculated.Z = ctrl.Margin.Z
		calculated.W = parentGlobalMargin.Z - ctrl.Margin.Y - ctrl.Margin.W
		calculated.X = ctrl.Margin.X
		calculated.Y = ctrl.Margin.Y
	} else if ctrl.Dock == margui.FillTop {
		calculated.Z = parentGlobalMargin.Z - ctrl.Margin.X - ctrl.Margin.Z
		calculated.W = ctrl.Margin.W
		calculated.X = ctrl.Margin.X
		calculated.Y = ctrl.Margin.Y
	} else if ctrl.Dock == margui.FillRight {
		calculated.Z = ctrl.Margin.Z
		calculated.W = parentGlobalMargin.Z - ctrl.Margin.Y - ctrl.Margin.W
		calculated.X = ctrl.Margin.X + parentGlobalMargin.Z
		calculated.Y = ctrl.Margin.Y
	} else if ctrl.Dock == margui.FillBottom {
		calculated.Z = parentGlobalMargin.Z - ctrl.Margin.X - ctrl.Margin.Z
		calculated.W = ctrl.Margin.W
		calculated.X = ctrl.Margin.X
		calculated.Y = ctrl.Margin.Y + parentGlobalMargin.W
	}

	//Add parent position + self margin left top
	calculated.X += parentGlobalMargin.X
	calculated.Y += parentGlobalMargin.Y

	return calculated
}

func drawControl(parent *comp.Control, ctrl *comp.Control) {

	//TODO Ко всему умножить пивот
	ctrl.GlobalMargin = calcGlobalMargin(parent, ctrl)

	//screenWidth
	//screenHeight

	points, vao, vbo := rectCoords(ctrl.GlobalMargin.Z, ctrl.GlobalMargin.W, ctrl.GlobalMargin.X, ctrl.GlobalMargin.Y)
	texture, program := newTexture(ctrl.Color)
	gl.UseProgram(program)
	//gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)

	//gl.Enable(gl.BLEND) // enable translucency
	drawTexture(texture, points)
	freeCoords(vao, vbo)

	for _, bCtrl := range ctrl.Buttons {
		drawControl(ctrl, &bCtrl.Control)
	}

	for _, pCtrl := range ctrl.Panels {
		drawControl(ctrl, &pCtrl.Control)
	}
}

func main() {
	wnd, err := unmarshalFromFile()
	if err != nil {
		panic(err)
	}

	err = glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	// make the window hidden, we will set it up and then show it later
	//glfw.WindowHint(glfw.Visible, 0)

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)
	glfw.WindowHint(glfw.Samples, 4)
	//glfw.WindowHint(glfw.AlphaBits, 1)
	//glfw.WindowHint(glfw.Decorated, 0)

	window, err := glfw.CreateWindow(screenWidth, screenHeight, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press && key == glfw.KeyEscape {
			w.SetShouldClose(true)
		}
	})

	if err = gl.Init(); err != nil {
		panic(err)
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
	gl.Disable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)

	gl.ClearColor(1, 1, 1, 0.0)
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		draw(wnd)
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
