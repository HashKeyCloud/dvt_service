package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/urfave/cli"
)

var apiUrl = "http://localhost:3000"

func setupApp() *cli.App {
	app := cli.NewApp()
	app.Usage = "DVT Tools"
	app.Version = "1.0.0"
	app.Description = fmt.Sprintf("Use API: %s", apiUrl)
	app.Commands = []cli.Command{
		{
			Name:      "run",
			ShortName: "r",
			Action:    action,
		},
	}

	app.Before = func(context *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}

	return app
}

func main() {
	if err := setupApp().Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func action(_ *cli.Context) error {
	pterm.DefaultBigText.WithLetters(putils.LettersFromString("DVT-Service")).Render()

	options := []string{
		uploadKeystore,
		setFeeRecipient,
		"Quit",
	}

	selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions(options).Show()
	pterm.Info.Printfln("Selected option: %s", pterm.Green(selectedOption))

	switch selectedOption {
	case setFeeRecipient:
		setFeeRecipientCommandAction()
	case uploadKeystore:
		uploadCommandAction()
	}

	return nil
}
