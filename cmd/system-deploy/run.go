package main

import (
	"context"
	"log"
	"strings"

	"github.com/khulnasoft-lab/system-conf/conf"
	"github.com/khulnasoft-lab/system-deploy/pkg/actions"
	"github.com/khulnasoft-lab/system-deploy/pkg/deploy"
	"github.com/khulnasoft-lab/system-deploy/pkg/runner"
	"github.com/spf13/cobra"
)

// Flags for the runActionCommand
var (
	flagRunOptions []string
)

var runActionCommand = &cobra.Command{
	Use:   "run",
	Short: "Execute a single action",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		_, ok := actions.GetPlugin(name)
		if !ok {
			log.Fatalf("unknown plugin: %s", name)
		}

		var opts conf.Options
		for _, o := range flagRunOptions {
			parts := strings.Split(o, "=")
			key := parts[0]
			value := strings.Join(parts[1:], "=")

			opts = append(opts, conf.Option{
				Name:  key,
				Value: value,
			})
		}

		s := conf.Section{
			Options: opts,
			Name:    name,
		}
		task := deploy.Task{
			Sections: []conf.Section{s},
		}

		r, err := runner.NewRunner(actions.NewLogger(), []deploy.Task{task})
		if err != nil {
			log.Fatalf("failed to prepare runner: %s", err)
		}

		if err := r.Deploy(context.Background()); err != nil {
			log.Fatalf("failed to deploy: %s", err)
		}
	},
}

func init() {
	flags := runActionCommand.Flags()
	{
		flags.StringSliceVarP(&flagRunOptions, "option", "o", nil, "Options for the action.")
	}
}
