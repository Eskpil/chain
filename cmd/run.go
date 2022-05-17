package cmd

import (
	"chain/context"
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
			fmt.Println("Insufficent arguments, expected filename argument.")
			return
		}

		data, err := os.ReadFile(args[0])

		if err != nil {
			fmt.Printf("Unexpected error when reading from file: %s: %s\n", args[0], err)
			return
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
				fmt.Printf("Unexpected error when unmarshaling data: %s\n", err)
				return
			}

			scope.RunProcedure(procedure)
		} else {
			project := structures.ProjectStructure{}

			err = yaml.Unmarshal(data, &project)

			if err != nil {
				fmt.Printf("Unexpected error when unmarshaling data: %s\n", err)
				return
			}

			for _, s := range project.Project.Procedures {
				filepath := fmt.Sprintf("%s/procedure.yml", s)

				data, err = os.ReadFile(filepath)

				if err != nil {
					fmt.Printf("Unexpected error when reading from file: %s: %s\n", filepath)
					return
				}

				procedure := structures.ProcedureStructure{}

				err := yaml.Unmarshal(data, &procedure)

				if err != nil {
					fmt.Printf("Unexpected error when unmarshaling data: %s\n", err)
					return
				}

				childScope := context.Scope{}

				childScope.InheritFrom(&scope, s)

				childScope.RunProcedure(procedure)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().BoolP("procedure", "p", false, "Tell Chain we only wan't to run this procedure")
}
