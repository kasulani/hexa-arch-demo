package app

import (
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/kasulani/hexa-arch-demo/internal/url"
)

type (
	cliPayload struct {
		args    []string
		longURL url.LongURL
		err     error
	}

	// command is cli command type.
	command struct {
		*cobra.Command
	}

	// SubCommand interface defines AddTo method.
	SubCommand interface {
		AddTo(root *rootCommand)
	}

	// subCommands is a slice of SubCommand.
	subCommands []SubCommand

	rootCommand    command
	shortenCommand command
)

func registerSubCommands(root *rootCommand, commands subCommands) {
	for _, command := range commands {
		command.AddTo(root)
	}
}

// AddTo adds shortenCommand to the root command.
func (shorten *shortenCommand) AddTo(root *rootCommand) {
	root.AddCommand(shorten.Command)
}

func newRootCommand() *rootCommand {
	return &rootCommand{
		Command: &cobra.Command{
			Use:     "shortener",
			Short:   "cli app to shorten a url",
			Long:    `cli app to shorten a url`,
			Version: "1.0.0",
		},
	}
}

func (p *cliPayload) invalid() bool {
	if len(p.args) != 1 {
		p.err = errors.New("command expects 1 argument")

		return true
	}

	p.longURL = url.LongURL(p.args[0])

	return false
}

func (p *cliPayload) error() error {
	return p.err
}

func newShortenCommand(ctx context.Context, cfg *config, useCase url.Shortener) *shortenCommand {
	return &shortenCommand{
		Command: &cobra.Command{
			Use:   "shorten",
			Short: "shortens a url",
			Long:  "This subcommand shortens a url",
			Run: func(cmd *cobra.Command, args []string) {
				input := cliPayload{args: args}

				if input.invalid() {
					fmt.Println("error:", input.error())
					os.Exit(1)
				}

				shortURL, err := useCase.ShortenURL(ctx, input.longURL)
				if err != nil {
					fmt.Println("command failed:", err)
					os.Exit(1)
				}

				fmt.Printf("%s://%s:%s/%s\n", cfg.Proto, cfg.Host, cfg.Port, shortURL.Code())
			},
		},
	}
}
