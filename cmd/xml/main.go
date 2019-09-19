package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"

	"github.com/GLeBaTi/margui/comp"

	"github.com/GLeBaTi/margui"
)

func main() {
	//unmarshalFromFile()
	marshlAndUnmarshal()
}

func unmarshalFromFile() {

	xmlFile, err := os.Open("main.xml")
	if err != nil {
		panic(err)
	}

	defer margui.Close(xmlFile)

	fileData, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		panic(err)
	}

	var out comp.Window
	err = xml.Unmarshal(fileData, &out)
	if err != nil {
		panic(err)
	}
	margui.LogIf("%+v\n", out)
}

func marshlAndUnmarshal() {
	bn1 := &comp.Button{
		Control: comp.Control{
			Position: margui.NewPos(10, 10),
		},
	}
	bn1.Text = "Ok"
	bn2 := &comp.Button{}
	bn3 := &comp.Button{}
	bn3.Text = "Ok2"

	var ii []*comp.Button
	bn2.Buttons = append(ii, bn3)
	var i []*comp.Button
	panel1 := &comp.Panel{}
	panel1.Buttons = append(i, bn1, bn2)

	wind := &comp.Window{}
	wind.Panels = append(wind.Panels, panel1)

	bytes, err := xml.MarshalIndent(wind, "  ", "    ")
	if err != nil {
		panic(err)
	}
	margui.LogI(string(bytes))
	margui.LogI("==========================\n")
	var out comp.Window

	err = xml.Unmarshal(bytes, &out)
	if err != nil {
		panic(err)
	}
	margui.LogIf("%+v\n", out)
}
