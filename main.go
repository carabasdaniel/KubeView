package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// Vars
var (
	AppName string
	BuiltAt string
	debug   = flag.Bool("d", false, "enables the debug mode")
	path    = flag.String("p", "assets/", "Set yaml decoding path")
	w       *astilectron.Window
)

var count int

func main() {
	// Init
	flag.Parse()
	astilog.FlagInit()

	count = 0
	// Run bootstrap
	astilog.Debugf("Running app built at %s", BuiltAt)
	if err := bootstrap.Run(bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
		},
		Debug:         *debug,
		OnWait:        sendMessage,
		RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage:       "index.html",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astilectron.PtrStr("#333"),
				Center:          astilectron.PtrBool(true),
				Height:          astilectron.PtrInt(700),
				Width:           astilectron.PtrInt(1000),
			},
		}},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}

}

func sendMessage(_ *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
	w = ws[0]
	go func() {

		files, err := ioutil.ReadDir(*path)
		if err != nil {
			panic(err)
		}
		if err := bootstrap.SendMessage(w, "wait", *path); err != nil {
			astilog.Error(errors.Wrap(err, "sending from wait failed"))
		}
		for _, f := range files {
			filepath := *path + "/" + f.Name()
			content, err := ioutil.ReadFile(filepath)
			if err != nil {
				panic(err)
			}
			decodeObjects := Decode(content)

			for _, obj := range decodeObjects {

				meta, err := json.Marshal(obj)
				if err != nil {
					fmt.Println(err)
				}
				if err := bootstrap.SendMessage(w, "object", string(meta)); err != nil {
					astilog.Error(errors.Wrap(err, "sending from object failed"))
				}
			}
		}
	}()
	return nil
}

func handleMessages(w *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "event.name":
		// Unmarshal payload
		var s string
		if err = json.Unmarshal(m.Payload, &s); err != nil {
			payload = err.Error()
			return
		}
		fmt.Println(payload)
	case "clicked":
		var s string
		if err = json.Unmarshal(m.Payload, &s); err != nil {
			payload = err.Error()
			return
		}
		fmt.Println("Loading from path: ", s)

		files, err := ioutil.ReadDir(s)
		if err != nil {
			panic(err)
		}
		if err := bootstrap.SendMessage(w, "wait", s); err != nil {
			astilog.Error(errors.Wrap(err, "sending from wait failed"))
		}
		for _, f := range files {
			filepath := s + "/" + f.Name()
			content, err := ioutil.ReadFile(filepath)
			if err != nil {
				panic(err)
			}
			decodeObjects := Decode(content)

			for _, obj := range decodeObjects {

				meta, err := json.Marshal(obj)
				if err != nil {
					return nil, err
				}
				if err := bootstrap.SendMessage(w, "object", string(meta)); err != nil {
					astilog.Error(errors.Wrap(err, "sending from object failed"))
				}
			}
		}
		payload = s
		return payload, nil
	}
	return
}
