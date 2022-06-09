package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func NewOpenCommand() *OpenCommand {
	gc := &OpenCommand{
		fs: flag.NewFlagSet("open", flag.ContinueOnError),
	}

	gc.fs.StringVar(&gc.url, "url", "http://example.com", "Target URL")

	return gc
}

type OpenCommand struct {
	fs *flag.FlagSet

	url string
}

func (g *OpenCommand) Name() string {
	return g.fs.Name()
}

func (g *OpenCommand) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *OpenCommand) Run() error {
	fmt.Println("Target:", g.url)
	resp, err := http.Get(g.url)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
	return nil
}

type Runner interface {
	Init([]string) error
	Run() error
	Name() string
}

func root(args []string) error {
	if len(args) < 1 {
		return errors.New("You must pass a sub-command")
	}

	cmds := []Runner{
		NewOpenCommand(),
	}

	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			cmd.Init(os.Args[2:])
			return cmd.Run()
		}
	}

	return fmt.Errorf("Unknown subcommand: %s", subcommand)
}

func main() {
	if err := root(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
