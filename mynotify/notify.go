package mynotify

import (
	"github.com/lxn/walk"
	"log"
)

func notify(info string, level string) {

	title := "CAN 数据配置工具"

	// We need either a walk.MainWindow or a walk.Dialog for their message loop.
	// We will not make it visible in this example, though.
	mw, err := walk.NewMainWindow()
	if err != nil {
		log.Fatal(err)
	}

	// We load our icon from a file.
	icon, err := walk.Resources.Icon("/favicon.ico")
	if err != nil {
		log.Fatal(err)
	}

	// Create the notify icon and make sure we clean it up on exit.
	ni, err := walk.NewNotifyIcon(mw)
	if err != nil {
		log.Fatal(err)
	}
	defer ni.Dispose()

	// Set the icon and a tool tip text.
	if err := ni.SetIcon(icon); err != nil {
		log.Fatal(err)
	}
	if err := ni.SetToolTip("点击查看信息，或右键菜单退出"); err != nil {
		log.Fatal(err)
	}

	// When the left mouse button is pressed, bring up our balloon.
	ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}

		if err := ni.ShowCustom(
			title,
			info,
			icon); err != nil {

			log.Fatal(err)
		}
	})

	// We put an exit action into the context menu.
	exitAction := walk.NewAction()
	if err := exitAction.SetText("退出"); err != nil {
		log.Fatal(err)
	}
	exitAction.Triggered().Attach(func() { walk.App().Exit(0) })
	if err := ni.ContextMenu().Actions().Add(exitAction); err != nil {
		log.Fatal(err)
	}

	// The notify icon is hidden initially, so we have to make it visible.
	if err := ni.SetVisible(true); err != nil {
		log.Fatal(err)
	}

	// Now that the icon is visible, we can bring up an info balloon.
	if level == "info" {
		if err := ni.ShowInfo(title, info); err != nil {
			log.Fatal(err)
		}
	}

	if level == "error" {
		if err := ni.ShowError(title, info); err != nil {
			log.Fatal(err)
		}
	}

	if level == "message" {
		if err := ni.ShowMessage(title, info); err != nil {
			log.Fatal(err)
		}
	}

	if level == "warning" {
		if err := ni.ShowWarning(title, info); err != nil {
			log.Fatal(err)
		}
	}

	// Run the message loop.
	mw.Run()
}

func Info(info string) {
	notify(info, "info")
}

func Message(info string) {
	notify(info, "message")
}

func Warning(info string) {
	notify(info, "warning")
}

func Error(info string) {
	notify(info, "error")
}
