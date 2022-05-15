/*
 * Copyright (c) 2022, Linus Johansen.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"chain/compilers"
	"chain/procedures"
	"fmt"
	"os"
)

func main() {
	clang := compilers.Clang{
		Path: "/usr/bin/clang",
	}

	build := procedures.BuildProcedure{
		Compiler: clang,
		Files:    []string{"test.c"},
	}

	link := procedures.LinkProcedure{
		Linker: clang,
		Files:  []string{"test.o"},
		Into:   "test",
	}

	// All build procedures will be ran scoped when the command build is called.
	err := build.RunProcedure()

	if err != nil {
		fmt.Println("Failed to run build procedure")
		os.Exit(1)
	}

	err = link.RunProcedure()

	fmt.Println(build)
	fmt.Println(link)
}
