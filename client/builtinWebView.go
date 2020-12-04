package client

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"strings"
)

func RunBuiltinWebView(url string) {

	//debug := true
	//w := webview.New(debug)
	//defer w.Destroy()
	//w.SetTitle("Minimal webview example")
	//w.SetSize(800, 600, webview.HintNone)
	//w.Navigate(url)
	//w.Run()

	var le *walk.LineEdit
	var wv *walk.WebView

	MainWindow{
		Title:   "内置浏览器",
		MinSize: Size{800, 600},
		Layout:  VBox{MarginsZero: true},
		Children: []Widget{
			LineEdit{
				AssignTo: &le,
				Text:     Bind("wv.URL"),
				OnKeyDown: func(key walk.Key) {
					if key == walk.KeyReturn {
						//idx := strings.Index(le.Text(), "http://")
						//if idx == -1 {
						//	wv.SetURL("http://" + le.Text())
						//	return
						//}

						wv.SetURL(le.Text())
					}
				},
			},
			WebView{
				AssignTo:                 &wv,
				Name:                     "wv",
				URL:                      url,
				ShortcutsEnabled:         true,
				NativeContextMenuEnabled: true,
			},
		},
		Functions: map[string]func(args ...interface{}) (interface{}, error){
			"icon": func(args ...interface{}) (interface{}, error) {
				if strings.HasPrefix(args[0].(string), "https") {
					return "check", nil
				}

				return "stop", nil
			},
		},
	}.Run()
}
