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

func getViewportXCoord(x float32) float32 {
	xPos := float32(x) / float32(screenWidth)
	return -1 + xPos*2
}

func getViewportYCoord(y float32) float32 {
	xPos := float32(y) / float32(screenHeight)
	return 1 - xPos*2
}

func rectCoords(width float32, height float32, posX float32, posY float32) ([]float32, uint32, uint32) {

	x1 := getViewportXCoord(posX)
	x2 := getViewportXCoord(posX + width)
	y1 := getViewportYCoord(posY)
	y2 := getViewportYCoord(posY + height)

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

	//Create buffer
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	//Tell, how OGL must copy vbo into vertShader (5nums * sizeof(float32))
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	return points, vao, vbo
}

func newLinearGradientTexture(ctrlGlobalMargin margui.XYZW, gradient comp.LinearGradientBrush) (texture uint32, program uint32) {
	vertex_shaderStr := `
    #version 110
    attribute vec3 vert;
    void main() {
        gl_Position = vec4(vert, 1);
    }
` + "\x00"

	fragment_shaderStr := fmt.Sprintf(`
    #version 110
	uniform vec2  gradientStartPos;
	uniform vec2  gradientEndPos;
	uniform int   numStops;
	uniform vec4  colors[16];
	uniform float stops[16];
    void main() {
        float  alpha = atan( -gradientEndPos.y + gradientStartPos.y, gradientEndPos.x - gradientStartPos.x );
    	float  gradientStartPosRotatedX = gradientStartPos.x*cos(alpha) - gradientStartPos.y*sin(alpha);
    	float  gradientEndPosRotatedX   = gradientEndPos.x*cos(alpha) - gradientEndPos.y*sin(alpha);
    	float  d = gradientEndPosRotatedX - gradientStartPosRotatedX;
    	
	
    	float y = gl_FragCoord.y;
    	float x = gl_FragCoord.x;
    	float xLocRotated = x*cos( alpha ) - y*sin( alpha );
 	
    	gl_FragColor = mix(colors[0], colors[1], smoothstep( gradientStartPosRotatedX + stops[0]*d, gradientStartPosRotatedX + stops[1]*d, xLocRotated ) );
    	for ( int i=1; i<numStops-1; ++i ) {
    	    gl_FragColor = mix(gl_FragColor, colors[i+1], smoothstep( gradientStartPosRotatedX + stops[i]*d, gradientStartPosRotatedX + stops[i+1]*d, xLocRotated ) );
    	}
    }
` + "\x00")

	vs, err := compileShader(vertex_shaderStr, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fs, err := compileShader(fragment_shaderStr, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	//invert Y axis for fg shader
	startPosX := (gradient.Start.X*ctrlGlobalMargin.Z + ctrlGlobalMargin.X)
	startPosY := screenHeight - (gradient.Start.Y*ctrlGlobalMargin.W + ctrlGlobalMargin.Y)
	endPosX := (gradient.End.X*ctrlGlobalMargin.Z + ctrlGlobalMargin.X)
	endPosY := screenHeight - (gradient.End.Y*ctrlGlobalMargin.W + ctrlGlobalMargin.Y)
	var stops []float32
	var colors [][4]float32
	for _, item := range gradient.Colors {
		stops = append(stops, item.Offset)
		colors = append(colors, [4]float32{item.Color.R, item.Color.G, item.Color.B, item.Color.A})
	}

	program = gl.CreateProgram()
	gl.AttachShader(program, fs)
	gl.AttachShader(program, vs)
	gl.LinkProgram(program)

	gl.DeleteShader(fs)
	gl.DeleteShader(vs)

	gl.ProgramUniform1i(program, gl.GetUniformLocation(program, gl.Str("numStops\x00")), int32(len(gradient.Colors)))
	gl.ProgramUniform2f(program, gl.GetUniformLocation(program, gl.Str("gradientStartPos\x00")), startPosX, startPosY)
	gl.ProgramUniform2f(program, gl.GetUniformLocation(program, gl.Str("gradientEndPos\x00")), endPosX, endPosY)
	gl.ProgramUniform4fv(program, gl.GetUniformLocation(program, gl.Str("colors\x00")), int32(len(colors)), &colors[0][0])
	gl.ProgramUniform1fv(program, gl.GetUniformLocation(program, gl.Str("stops\x00")), int32(len(stops)), &stops[0])

	return texture, program
}

func newRadialGradientTexture(ctrlGlobalMargin margui.XYZW, gradient comp.RadialGradientBrush) (texture uint32, program uint32) {
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
    uniform int   numStops;
	uniform vec4  colors[16];
	uniform float stops[16];
	uniform vec2 center;

    void main() {
        vec2 pos_ndc = 2.0 * gl_FragCoord.xy / vec2(80, 80) - 1.0;
    	float dist = length(pos_ndc);

    	gl_FragColor = mix(colors[0], colors[1], smoothstep( stops[0], stops[1], dist ) );
    	for ( int i=1; i<numStops-1; ++i ) {
        	gl_FragColor = mix(gl_FragColor, colors[i+1], smoothstep( stops[i], stops[i+1], dist ) );
    	}
    }
` + "\x00")

	vs, err := compileShader(vertex_shaderStr, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fs, err := compileShader(fragment_shaderStr, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	//invert Y axis for fg shader
	centerPosX := (gradient.Pivot.X*ctrlGlobalMargin.Z + ctrlGlobalMargin.X)
	centerPosY := screenHeight - (gradient.Pivot.Y*ctrlGlobalMargin.W + ctrlGlobalMargin.Y)
	var stops []float32
	var colors [][4]float32
	for _, item := range gradient.Colors {
		stops = append(stops, item.Offset)
		colors = append(colors, [4]float32{item.Color.R, item.Color.G, item.Color.B, item.Color.A})
	}

	program = gl.CreateProgram()
	gl.AttachShader(program, fs)
	gl.AttachShader(program, vs)
	gl.LinkProgram(program)

	gl.DeleteShader(fs)
	gl.DeleteShader(vs)

	gl.ProgramUniform1i(program, gl.GetUniformLocation(program, gl.Str("numStops\x00")), int32(len(gradient.Colors)))
	gl.ProgramUniform2f(program, gl.GetUniformLocation(program, gl.Str("center\x00")), centerPosX, centerPosY)
	gl.ProgramUniform4fv(program, gl.GetUniformLocation(program, gl.Str("colors\x00")), int32(len(colors)), &colors[0][0])
	gl.ProgramUniform1fv(program, gl.GetUniformLocation(program, gl.Str("stops\x00")), int32(len(stops)), &stops[0])

	return texture, program
}

func newSolidColorTexture(color margui.Color) (texture uint32, program uint32) {
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

	gl.DeleteShader(fs)
	gl.DeleteShader(vs)

	return texture, program
}

func draw(win *comp.Rectangle) {
	drawControl(nil, &win.Geometry)
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

func unmarshalFromFile() (*comp.Rectangle, error) {

	xmlFile, err := os.Open("main.xml")
	if err != nil {
		return nil, err
	}

	defer margui.Close(xmlFile)

	fileData, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return nil, err
	}

	var out comp.Rectangle
	err = xml.Unmarshal(fileData, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func calcGlobalMargin(parent *comp.Geometry, ctrl *comp.Geometry) margui.XYZW {
	var calculated margui.XYZW
	var parentGlobalMargin = margui.XYZW{
		X: 0,
		Y: 0,
		Z: screenWidth,
		W: screenHeight,
	}

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

	if ctrl.Dock == margui.Fill {

	} else if ctrl.Dock == margui.FillLeft || ctrl.Dock == margui.FillRight || ctrl.Dock == margui.FillVertical {
		calculated.X += (calculated.Z / 2.0) * (-ctrl.Pivot.X - 1.0)
	} else if ctrl.Dock == margui.FillTop || ctrl.Dock == margui.FillBottom || ctrl.Dock == margui.FillHorizontal {
		calculated.Y += (calculated.W / 2.0) * (ctrl.Pivot.Y - 1.0)
	} else {
		calculated.X += (calculated.Z / 2.0) * (-ctrl.Pivot.X - 1.0)
		calculated.Y += (calculated.W / 2.0) * (ctrl.Pivot.Y - 1.0)
	}
	return calculated
}

func drawControl(parent *comp.Geometry, ctrl *comp.Geometry) {

	ctrl.GlobalMargin = calcGlobalMargin(parent, ctrl)

	points, vao, vbo := rectCoords(ctrl.GlobalMargin.Z, ctrl.GlobalMargin.W, ctrl.GlobalMargin.X, ctrl.GlobalMargin.Y)

	var texture, program uint32

	if ctrl.Color != nil {
		texture, program = newSolidColorTexture(*ctrl.Color)
	} else if ctrl.BackgroundSolidColor != nil {
		texture, program = newSolidColorTexture(ctrl.BackgroundSolidColor.Color)
	} else if ctrl.BackgroundLinearGradient != nil {
		texture, program = newLinearGradientTexture(ctrl.GlobalMargin, *ctrl.BackgroundLinearGradient)
	} else if ctrl.BackgroundRadialGradient != nil {
		texture, program = newRadialGradientTexture(ctrl.GlobalMargin, *ctrl.BackgroundRadialGradient)
	} else {
		texture, program = newSolidColorTexture(margui.Color{
			R: 1,
			G: 1,
			B: 1,
			A: 0,
		})
	}

	gl.UseProgram(program)
	//gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)

	_ = texture
	//gl.Enable(gl.BLEND) // enable translucency
	//gl.ActiveTexture(gl.TEXTURE0)
	//gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(len(points)/5))
	freeCoords(vao, vbo)

	//TODO как-то автоматом сделать
	for _, pCtrl := range ctrl.Rectangles {
		drawControl(ctrl, &pCtrl.Geometry)
	}
	for _, pCtrl := range ctrl.Ellipses {
		drawControl(ctrl, &pCtrl.Geometry)
	}
	for _, pCtrl := range ctrl.Paths {
		drawControl(ctrl, &pCtrl.Geometry)
	}
	for _, pCtrl := range ctrl.Polygons {
		drawControl(ctrl, &pCtrl.Geometry)
	}
}

func GetAllControls(parent *comp.Control) []*comp.Control {
	var result []*comp.Control
	result = append(result, parent)
	for _, pCtrl := range parent.Rectangles {
		result = append(result, GetAllControls(&pCtrl.Control)...)
	}
	for _, pCtrl := range parent.Ellipses {
		result = append(result, GetAllControls(&pCtrl.Control)...)
	}
	for _, pCtrl := range parent.Paths {
		result = append(result, GetAllControls(&pCtrl.Control)...)
	}
	for _, pCtrl := range parent.Polygons {
		result = append(result, GetAllControls(&pCtrl.Control)...)
	}
	return result
}

func GetAllInteractables(parent interface{}, parentCtrl *comp.Control) []comp.Interactable {
	var result []comp.Interactable
	interactable, ok := parent.(comp.Interactable)
	if ok && parentCtrl.IsInteractable {
		result = append(result, interactable)
	}
	for _, pCtrl := range parentCtrl.Rectangles {
		result = append(result, GetAllInteractables(pCtrl, &pCtrl.Control)...)
	}
	for _, pCtrl := range parentCtrl.Ellipses {
		result = append(result, GetAllInteractables(pCtrl, &pCtrl.Control)...)
	}
	for _, pCtrl := range parentCtrl.Paths {
		result = append(result, GetAllInteractables(pCtrl, &pCtrl.Control)...)
	}
	for _, pCtrl := range parentCtrl.Polygons {
		result = append(result, GetAllInteractables(pCtrl, &pCtrl.Control)...)
	}
	return result
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

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile) //error if using deprecated OpenGL
	//glfw.WindowHint(glfw.Samples, 4)
	//glfw.WindowHint(glfw.AlphaBits, 1)
	//glfw.WindowHint(glfw.Decorated, 0)

	window, err := glfw.CreateWindow(screenWidth, screenHeight, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}

	window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	window.MakeContextCurrent()

	window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		for _, pCtrl := range GetAllInteractables(wnd, &wnd.Control) {
			pCtrl.CursorButtonEventHandler(button, action, mod)
		}
	})

	window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		for _, pCtrl := range GetAllInteractables(wnd, &wnd.Control) {
			pCtrl.CursorEventHandler(xpos, ypos)
		}
	})

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press && key == glfw.KeyEscape {
			w.SetShouldClose(true)
		}
		for _, pCtrl := range GetAllInteractables(wnd, &wnd.Control) {
			pCtrl.KeyEventHandler(key, scancode, action, mods)
		}
	})

	if err = gl.Init(); err != nil {
		panic(err)
	}
	gl.Viewport(0, 0, screenWidth, screenHeight)
	gl.Enable(gl.DEPTH_TEST) //Always
	gl.DepthFunc(gl.LEQUAL)
	gl.Disable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	//TODO glOrtho

	gl.ClearColor(1, 1, 1, 0.0)
	for !window.ShouldClose() {
		glfw.PollEvents()
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		draw(wnd)
		window.SwapBuffers()
	}
}
