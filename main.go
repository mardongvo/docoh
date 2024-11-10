package main

import (
	"docoh/db"
	"errors"
	"os"

	cli "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "docoh",
		Usage: "Organize document coherency between source code and external docs",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "db",
				Value: ".docohdb",
			},
		},
		ErrWriter: os.Stderr,
		Commands: []*cli.Command{
			{
				Name:      "add",
				Aliases:   []string{"a"},
				Args:      true,
				Usage:     "add documenting rule",
				ArgsUsage: "add <target> <source path|glob pattern>",
				Action:    addRule,
			},
			{
				Name:      "refresh",
				Aliases:   []string{"r"},
				Args:      true,
				Usage:     "refresh file hashes for rule",
				ArgsUsage: "refresh [-n N | -t target]",
				Action:    refreshRule,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name: "n",
					},
					&cli.StringFlag{
						Name: "t",
					},
				},
			},
			{
				Name:      "report",
				Args:      true,
				Usage:     "report changes for rule(s)",
				ArgsUsage: "report [-n N | -t target]",
				Action:    reportRules,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name: "n",
					},
					&cli.StringFlag{
						Name: "t",
					},
				},
			},
		},
	}

	app.Run(os.Args)
}

func tryLoad(path string) (*db.DB, error) {
	f, err := os.Open(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return db.NewDB(), nil
	}
	defer f.Close()

	res := db.NewDB()
	err = res.Load(f)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func trySave(store *db.DB, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	err = store.Save(f)
	if err != nil {
		return err
	}
	return nil
}
