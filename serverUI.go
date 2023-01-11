package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func makeServer(win fyne.Window, s *Server) {
	btnStart := widget.NewButton("Server Start", serverStart(s))
	container := container.NewVBox(btnStart)
	win.SetContent(container)
	win.Resize(fyne.NewSize(500, 300))
}
