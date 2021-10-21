package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Run() {
	a := app.New()
	w := a.NewWindow("clipOne")

	hello := widget.NewLabel("Let me stay here!")
	center := container.NewCenter(
		hello,
	)
	w.Resize(fyne.NewSize(240, 80))
	w.SetContent(center)

	w.ShowAndRun()
}
