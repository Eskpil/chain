package cmd

import (
	"chain/compilers"
	"chain/procedures"
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

		procedure := structures.ProcedureStructure{}

		err = yaml.Unmarshal(data, &procedure)

		if err != nil {
			fmt.Printf("Unexpected error when unmarshaling data: %s\n", err)
			return
		}

		compiler := compilers.Clang{
			Path: "/usr/bin/clang",
		}

		buildProcedure := procedures.BuildProcedure{
			Files:    procedure.Procedure.Build.Files,
			Compiler: compiler,
		}

		var target procedures.Target

		if procedure.Procedure.Link.Target == "library" {
			target = procedures.Library
		} else {
			target = procedures.Binary
		}

		linkProcedure := procedures.LinkProcedure{
			Files:  procedure.Procedure.Link.Files,
			Target: target,
			Into:   procedure.Procedure.Link.Into,
			Linker: compiler,
		}

		err = buildProcedure.RunProcedure()

		if err != nil {
			return
		}

		err = linkProcedure.RunProcedure()

		if err != nil {
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
