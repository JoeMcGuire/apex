// Package list outputs a list of Lambda function information.
package list

import (
	"fmt"

	"github.com/apex/log"
	"github.com/spf13/cobra"

	"github.com/apex/apex/cmd/apex/root"
	"github.com/apex/apex/colors"
)

// tfvars output format.
var tfvars bool

// example output.
const example = `  List all functions
  $ apex list

  Output list as Terraform variables (.tfvars)
  $ apex list --tfvars`

// Command config.
var Command = &cobra.Command{
	Use:     "list",
	Short:   "Output functions list",
	Example: example,
	RunE:    run,
}

// Initialize.
func init() {
	root.Register(Command)

	f := Command.Flags()
	f.BoolVar(&tfvars, "tfvars", false, "Output as Terraform variables")
}

// Run command.
func run(c *cobra.Command, args []string) error {
	if err := root.Project.LoadFunctions(); err != nil {
		return err
	}

	if tfvars {
		outputTFvars()
	} else {
		outputList()
	}

	return nil
}

// outputTFvars format.
func outputTFvars() {
	for _, fn := range root.Project.Functions {
		config, err := fn.GetConfig()
		if err != nil {
			log.Debugf("can't fetch function config: %s", err.Error())
			continue
		}

		fmt.Printf("apex_function_%s=%q\n", fn.Name, *config.Configuration.FunctionArn)
	}
}

// outputList format.
func outputList() {
	fmt.Println()
	for _, fn := range root.Project.Functions {
		fmt.Printf("  \033[%dm%s\033[0m\n", colors.Blue, fn.Name)
		if fn.Description != "" {
			fmt.Printf("    description: %v\n", fn.Description)
		}
		fmt.Printf("    runtime: %v\n", fn.Runtime)
		fmt.Printf("    memory: %vmb\n", fn.Memory)
		fmt.Printf("    timeout: %vs\n", fn.Timeout)
		fmt.Printf("    role: %v\n", fn.Role)
		fmt.Printf("    handler: %v\n", fn.Handler)

		config, err := fn.GetConfigCurrent()
		if err != nil {
			fmt.Println()
			continue // ignore
		}

		fmt.Printf("    current version: %s\n", *config.Configuration.Version)

		fmt.Println()
	}
}
