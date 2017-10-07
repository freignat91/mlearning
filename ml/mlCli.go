package main

import (
	"github.com/fatih/color"
	"os"
)

type mlCLI struct {
	printColor [6]*color.Color
	server     string
	verbose    bool
	silence    bool
	debug      bool
}

var currentColorTheme = "default"
var (
	colRegular = 0
	colInfo    = 1
	colWarn    = 2
	colError   = 3
	colSuccess = 4
	colDebug   = 5
)

func (m *mlCLI) init() error {
	m.setColors()
	//
	return nil
}

func (m *mlCLI) printf(col int, format string, args ...interface{}) {
	if m.silence {
		return
	}
	colorp := m.printColor[0]
	if col > 0 && col < len(m.printColor) {
		colorp = m.printColor[col]
	}
	if !m.verbose && col == colInfo {
		return
	}
	if !m.debug && col == colDebug {
		return
	}
	colorp.Printf(format, args...)
}

func (m *mlCLI) Fatal(format string, args ...interface{}) {
	m.printf(colError, format, args...)
	os.Exit(1)
}

func (m *mlCLI) pError(format string, args ...interface{}) {
	m.printf(colError, format, args...)
}

func (m *mlCLI) pWarn(format string, args ...interface{}) {
	m.printf(colWarn, format, args...)
}

func (m *mlCLI) pInfo(format string, args ...interface{}) {
	m.printf(colInfo, format, args...)
}

func (m *mlCLI) pSuccess(format string, args ...interface{}) {
	m.printf(colSuccess, format, args...)
}

func (m *mlCLI) pRegular(format string, args ...interface{}) {
	m.printf(colRegular, format, args...)
}

func (m *mlCLI) pDebug(format string, args ...interface{}) {
	m.printf(colDebug, format, args...)
}

func (m *mlCLI) setColors() {
	theme := config.colorTheme
	if theme == "dark" {
		m.printColor[0] = color.New(color.FgHiWhite)
		m.printColor[1] = color.New(color.FgHiBlack)
		m.printColor[2] = color.New(color.FgYellow)
		m.printColor[3] = color.New(color.FgRed)
		m.printColor[4] = color.New(color.FgGreen)
		m.printColor[5] = color.New(color.FgHiBlack)
	} else {
		m.printColor[0] = color.New(color.FgMagenta)
		m.printColor[1] = color.New(color.FgHiBlack)
		m.printColor[2] = color.New(color.FgYellow)
		m.printColor[3] = color.New(color.FgRed)
		m.printColor[4] = color.New(color.FgGreen)
		m.printColor[5] = color.New(color.FgHiBlack)
	}
	//add theme as you want.
}
