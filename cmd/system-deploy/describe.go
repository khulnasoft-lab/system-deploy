package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/fatih/color"
	"github.com/khulnasoft-lab/system-conf/conf"
	"github.com/khulnasoft-lab/system-deploy/pkg/actions"
	"github.com/khulnasoft-lab/system-deploy/pkg/deploy"
	"github.com/mitchellh/go-wordwrap"
	"github.com/spf13/cobra"
)

var markdown bool

func mainHeader(val string) string {
	if markdown {
		return "# " + val
	}

	return color.New(color.Bold, color.Underline).Sprint(strings.ToUpper("[ " + val + " ]"))
}

func header(val string) string {
	if markdown {
		return "## " + val
	}

	return color.New(color.Bold, color.Underline).Sprint(strings.ToUpper(val))
}

func bold(val string) string {
	if markdown {
		return "**" + val + "**"
	}

	return color.New(color.Bold).Sprint(val)
}

func codeBlock(code string) string {
	if markdown {
		return "```ini\n" + code + "\n```"
	}

	return code
}

func underlineOrItalic(val string) string {
	if markdown {
		return "*" + val + "*"
	}

	return color.New(color.Underline).Sprint(val)
}

func init() {
	describe.Flags().BoolVar(&markdown, "markdown", false, "Print output in markdown")
}

func getTaskDummyPlugin() actions.Plugin {
	return actions.Plugin{
		Name:        "Task",
		Description: taskDescription,
		Website:     "https://khulnasoft-lab.github.io/system-deploy",
		Options:     deploy.TaskOptions(),
	}
}

var describe = &cobra.Command{
	Use:   "describe",
	Short: "Display documentation for an action",
	Run: func(_ *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Printf(" - %s\n", strings.Join(actions.ListActions(), "\n - "))
			return
		}

		if len(args) > 1 {
			log.Fatal("only one parameter expected")
		}

		var plg actions.Plugin
		if args[0] == "task" {
			plg = getTaskDummyPlugin()
		} else {
			var ok bool
			plg, ok = actions.GetPlugin(args[0])
			if !ok {
				log.Fatalf("Action %s does not exist", args[0])
			}
		}

		fmt.Println(mainHeader(plg.Name))

		if plg.Description != "" {
			fmt.Printf("\n%s\n", wrap(plg.Description, ""))
		}

		for _, section := range plg.Help {
			if section.Title != "" {
				fmt.Printf("\n%s\n", header(section.Title))
			}

			if section.Description != "" {
				fmt.Printf("\n%s\n", wrap(section.Description, ""))
			}
		}

		if !conf.IsAllowAny(plg.Options) {
			fmt.Printf("\n%s\n\n", header("Options"))

			for _, opt := range plg.Options {
				// skip internal options.
				if opt.Internal {
					continue
				}

				required := ""
				defaultValue := ""

				if opt.Required {
					required = " (required)"
				}

				if opt.Default != "" {
					defaultValue = fmt.Sprintf(" (Default: %q)", opt.Default)
				}

				fmt.Printf("   %s= (%s)", bold(opt.Name), opt.Type.String())
				for _, alias := range opt.Aliases {
					fmt.Printf("  \n   %s=", bold(alias))
				}

				fmt.Printf("  \n      %s\n\n",
					wrap(opt.Description+required+defaultValue, "      "),
				)
			}
		} else {
			fmt.Println("Any options allowed")
		}

		if plg.Example != "" {
			fmt.Printf("\n%s\n", header("Example"))
			fmt.Printf("\n%s\n", codeBlock(plg.Example))
		}

		if plg.Author != "" || plg.Website != "" {
			fmt.Printf("\n%s\n", header("Contact"))
			fmt.Printf("\n%s  ", underlineOrItalic(plg.Author))
			fmt.Printf("\n%s  ", plg.Website)

			fmt.Println()
		}
	},
}

// wrap ensures that text is no longer thatn 80 characters per line.
// It automatically breaks text into multiple lines that fit into a
// 80 character (including indention) limit.
func wrap(text string, indent string) string {
	lines := strings.Split(wordwrap.WrapString(text, uint(80-len(indent))), "\n")
	return strings.Join(lines, "\n"+indent)
}

const taskDescription = "The `[Task]` section must be available in each deploy unit file and contains metadata for the task like a human-readable " +
	"description or the tasks state (Disabled or Masked). Users may also defined condition and assertions in the tasks meta section that may disable or " +
	"fail the task based on environmental conditions. All properties starting with `Condition` will disable the task if not matched, all properties starting " +
	"with `Assert` will cause *system-deploy* to throw an error and exit."
