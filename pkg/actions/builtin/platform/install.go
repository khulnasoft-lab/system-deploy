package platform

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/khulnasoft-lab/system-conf/conf"
	"github.com/khulnasoft-lab/system-deploy/pkg/actions"
	"github.com/khulnasoft-lab/system-deploy/pkg/deploy"
)

var (
	pacmanNothingToDoRegex = regexp.MustCompile("\n[ \t]{1}there is nothing to do\n")
	aptChangedRegex        = regexp.MustCompile("[1-9]+[0-9]* (upgraded|newly|to remove)")
)

func init() {
	actions.MustRegister(actions.Plugin{
		Name:        "InstallPackages",
		Author:      "Md Sulaiman <infosulaimanbd@gmail.com>",
		Website:     "https://github.com/khulnasoft-lab/system-deploy",
		Description: "Install software packages using various package managers. For more control on the installation behavior use the Exec section instead.",
		Options: []conf.OptionSpec{
			{
				Name:        "AptPkgs",
				Description: "Packages to install if APT is available",
				Type:        conf.StringSliceType,
			},
			{
				Name:        "PacmanPkgs",
				Description: "Packages to install if Pacman is available",
				Type:        conf.StringSliceType,
			},
			// TODO(khulnasoft-lab): add support for DNF
			// TODO(khulnasoft-lab): add support for snap
			// TODO(khulnasoft-lab): add support for arch-linux AUR (maybe using yay?)
			/*
				{
					Name:        "DnfPkgs",
					Description: "Packages to install if DNF is available",
					Type:        deploy.StringSliceType,
				},
				{
					Name:        "SnapPkgs",
					Description: "Packages to install if Snap is available",
					Type:        deploy.StringSliceType,
				},
			*/
		},
		Setup: setupInstallAction,
	})
}

func setupInstallAction(task deploy.Task, sec conf.Section) (actions.Action, error) {
	aptPkgs := getPackages("AptPkgs", sec)
	pacmanPkgs := getPackages("PacmanPkgs", sec)
	dnfPkgs := getPackages("DnfPkgs", sec)
	snapPkgs := getPackages("SnapPkgs", sec)

	if len(aptPkgs) == 0 && len(pacmanPkgs) == 0 && len(dnfPkgs) == 0 && len(snapPkgs) == 0 {
		return nil, fmt.Errorf("no packages to install")
	}

	return &installAction{
		aptPkgs:    aptPkgs,
		pacmanPkgs: pacmanPkgs,
		dnfPkgs:    dnfPkgs,
		snapPkgs:   snapPkgs,
	}, nil
}

func getPackages(configKey string, sec conf.Section) []string {
	var pkgs []string
	pkgOpts := sec.GetStringSlice(configKey)

	for _, p := range pkgOpts {
		pkgs = append(pkgs, strings.Fields(p)...)
	}

	return pkgs
}

type installAction struct {
	actions.Base

	aptPkgs    []string
	pacmanPkgs []string
	dnfPkgs    []string
	snapPkgs   []string
}

func (ia *installAction) Name() string {
	return "Installing packages"
}

func (ia *installAction) Prepare(graph actions.ExecGraph) error {
	return nil
}

func (ia *installAction) Execute(ctx context.Context) (bool, error) {
	managers := getPackageManagers()

	var changed bool
	for _, m := range managers {
		var err error

		switch m {
		case Pacman:
			if len(ia.pacmanPkgs) == 0 {
				continue
			}
			changed, err = installPacman(ctx, ia.pacmanPkgs...)

		case APT:
			if len(ia.aptPkgs) == 0 {
				continue
			}
			changed, err = installApt(ctx, ia.aptPkgs...)

		default:
			continue
		}

		if err != nil {
			return false, err
		}
	}

	return changed, nil
}

func installPacman(ctx context.Context, pkgs ...string) (bool, error) {
	args := []string{
		"-S",
		"--needed",
		"--quiet",
		"--noconfirm",
	}

	args = append(args, pkgs...)
	cmd := exec.CommandContext(ctx, "pacman", args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "LC_ALL=C")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to install packages: %w\n%s", err, string(output))
	}

	nothingToDo := pacmanNothingToDoRegex.Match(output)
	return !nothingToDo, nil
}

func installApt(ctx context.Context, pkg ...string) (bool, error) {
	args := []string{
		"install",
		"-y",
	}
	args = append(args, pkg...)

	cmd := exec.CommandContext(ctx, "apt", args...)
	cmd.Env = os.Environ()

	cmd.Env = append(cmd.Env, "DEBCONF_FRONTEND='noninteractive'")
	cmd.Env = append(cmd.Env, "LC_ALL=C")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to install packages: %w\n%s", err, string(output))
	}

	hasChanged := aptChangedRegex.Match(output)

	return hasChanged, nil
}
