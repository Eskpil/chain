package cmd

import (
	"chain/context"
	"chain/logger"
	"chain/structures"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			logger.Error.Printf("Insufficent arguments, expected filename argument\n")
			os.Exit(1)
		}

		data, err := os.ReadFile(args[0])

		if err != nil {
			logger.Error.Printf("Unexpected error when reading data\n")
			logger.PrintError(fmt.Sprintf("%s\n", err))
			os.Exit(1)
		}

		isProcedure, _ := cmd.Flags().GetBool("procedure")

		scope := context.Scope{
			Parent:   nil,
			BuildDir: "bld",
			Prefix:   ".",
		}

		if isProcedure {
			procedure := structures.ProcedureStructure{}

			err = yaml.Unmarshal(data, &procedure)

			if err != nil {
				logger.Error.Printf("Unexpected error when unmarshaling data\n")
				logger.PrintError(fmt.Sprintf("%s\n", err))
				os.Exit(1)
			}

			scope.RunProcedure(procedure)
		} else {
			project := structures.ProjectStructure{}

			err = yaml.Unmarshal(data, &project)

			if err != nil {
				logger.Error.Printf("Unexpected error when unmarshaling data\n")
				logger.PrintError(fmt.Sprintf("%s\n", err))
				os.Exit(1)
			}

			project.Validate(args[0])

			scope.RunProject(&project)

		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().BoolP("procedure", "p", false, "Tell Chain we only wan't to run this procedure")
}
