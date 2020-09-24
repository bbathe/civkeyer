package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/lxn/walk"
	declarative "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"gopkg.in/yaml.v2"
)

var (
	appName    = "CIV Keyer"
	errorTitle = fmt.Sprintf("%s Error", appName)
)

// msgError displays dialog to user with error details
func msgError(p *walk.MainWindow, err error) {
	walk.MsgBox(p, errorTitle, err.Error(), walk.MsgBoxIconError)
}

// xFunc is a wrapper to call executeFunction and notify user of any errors
func xFunc(p *walk.MainWindow, f int) {
	err := executeFunction(f)
	if err != nil {
		log.Printf("%+v", err)
		msgError(p, err)
	}
}

// civkeyerWindow creates the main window and begins processing of user input
func civkeyerWindow() error {
	var tempWin *walk.MainWindow
	var mainWin *walk.MainWindow

	// load app icon
	ico, err := walk.Resources.Icon("3")
	if err != nil {
		log.Printf("%+v", err)
		return err
	}

	// this window is used to be the parent of any error messages during initialization
	tw := declarative.MainWindow{
		AssignTo: &tempWin,
		Title:    appName,
		Icon:     ico,
		Size:     declarative.Size{Width: 1, Height: 1},
		Layout:   declarative.Grid{},
		Visible:  false,
	}

	// create temporary window
	err = tw.Create()
	if err != nil {
		log.Printf("%+v", err)
		return err
	}

	//
	// complete app initialization here so we can message the user if there's an issue
	//

	// process command line
	flg := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flg.StringVar(&configFile, "config", "", "Configuration file")
	err = flg.Parse(os.Args[1:])
	if err != nil {
		e := fmt.Errorf("%s\n\nUsage of %s\n  -config string\n    Configuration file", err.Error(), os.Args[0])
		msgError(tempWin, e)
		log.Printf("%+v", err)
		return err
	}

	// log file is in the same directory as the executable with the same base name
	fn, err := os.Executable()
	if err != nil {
		msgError(tempWin, err)
		log.Printf("%+v", err)
		return err
	}
	basefn := strings.TrimSuffix(fn, path.Ext(fn))

	// log to file
	f, err := os.OpenFile(basefn+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		msgError(tempWin, err)
		log.Printf("%+v", err)
		return err
	}
	defer f.Close()
	log.SetOutput(f)

	// read config
	var cfn string
	if len(configFile) > 0 {
		// if user passed a filename, use that
		cfn = configFile
	} else {
		// default config file is in the same directory as the executable with the same base name
		cfn = basefn + ".yaml"
	}

	// #nosec G304
	bytes, err := ioutil.ReadFile(cfn)
	if err != nil {
		msgError(tempWin, err)
		log.Printf("%+v", err)
		return err
	}
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		msgError(tempWin, err)
		log.Printf("%+v", err)
		return err
	}

	// now create main window
	mw := declarative.MainWindow{
		AssignTo: &mainWin,
		Title:    appName,
		Icon:     ico,
		Size:     declarative.Size{Width: 1, Height: 1},
		Layout:   declarative.Grid{},
	}

	// create a button for each function defined
	mw.Children = make([]declarative.Widget, len(config.Functions))
	for i, f := range config.Functions {
		fnum := i
		mw.Children[i] = declarative.PushButton{Text: f.Label, Row: (i / 4), Column: (i % 4), OnClicked: func() { xFunc(mainWin, fnum) }}
	}

	// create a menu item for each function defined so the hotkeys work
	m := declarative.Menu{Visible: false}
	m.Items = make([]declarative.MenuItem, len(config.Functions))
	for i, f := range config.Functions {
		fnum := i
		m.Items[i] = declarative.Action{Text: f.Label, OnTriggered: func() { xFunc(mainWin, fnum) }, Shortcut: declarative.Shortcut{Modifiers: 0, Key: walk.Key((int(walk.KeyF1) + i))}}
	}
	mw.MenuItems = []declarative.MenuItem{m}

	// create window
	err = mw.Create()
	if err != nil {
		msgError(tempWin, err)
		log.Printf("%+v", err)
		return err
	}

	// disable maximize, minimize, and resizing
	hwnd := mainWin.Handle()
	win.SetWindowLong(hwnd, win.GWL_STYLE, win.GetWindowLong(hwnd, win.GWL_STYLE) & ^(win.WS_MAXIMIZEBOX|win.WS_MINIMIZEBOX|win.WS_SIZEBOX))

	// close temporary window
	win.DestroyWindow(tempWin.Handle())

	// start message loop
	mainWin.Run()

	return nil
}
