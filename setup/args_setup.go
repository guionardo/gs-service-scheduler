package setup

import (
	"fmt"
	"os"
	"strings"
)

type Args struct {
	SetupFolder string
	DoInstall   bool
	DoUninstall bool
	DoHelp      bool
}

func getArgs(runArgs []string) (args *Args, err error) {
	args = &Args{}
	for _, arg := range runArgs {
		if strings.HasPrefix(arg, "--setup-folder=") {
			words := strings.SplitN(arg, "=", 2)
			if len(words) == 2 {
				args.SetupFolder = words[1]
			}
		} else if arg == "--install" {
			args.DoInstall = true
		} else if arg == "--uninstall" {
			args.DoUninstall = true
		} else if arg == "--help" {
			args.DoHelp = true
		}
	}
	// Validate
	if err = isValidPath(args.SetupFolder, true); err != nil {
		err = fmt.Errorf("setup folder not found: %s", args.SetupFolder)
	}
	if args.DoInstall && args.DoUninstall {
		err = fmt.Errorf("--install and --uninstall arguments are mutually exclusive")
	}
	return
}

func GetArgs() (args *Args, err error) {
	return getArgs(os.Args)
}
