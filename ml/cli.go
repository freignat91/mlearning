package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	//RootCmd .
	RootCmd = &cobra.Command{
		Use:   `ml [OPTIONS] COMMAND [arg...]`,
		Short: "machine learning",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.UsageString())
		},
	}
)

func cli() {
	RootCmd.PersistentFlags().BoolVarP(&mlCli.verbose, "verbose", "v", false, `Verbose output`)
	RootCmd.PersistentFlags().BoolVarP(&mlCli.silence, "silence", "s", false, `Silence output`)
	RootCmd.PersistentFlags().BoolVar(&mlCli.debug, "debug", false, `Silence output`)
	cobra.OnInitialize(func() {
		if err := mlCli.init(); err != nil {
			fmt.Printf("Init error: %v\n", err)
			os.Exit(1)
		}
	})

	// versionCmd represents the mlearning version
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Display the version number of mLearning",
		Long:  `Display the version number of mLearning`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("mLearning version: %s\n", Version)
		},
	}
	RootCmd.AddCommand(versionCmd)

	//Execute commad
	cmd, _, err := RootCmd.Find(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := cmd.Execute(); err != nil {
		fmt.Printf("Error during: %s: %v\n", cmd.Name(), err)
		os.Exit(1)
	}

	os.Exit(0)
}

func displayList(lines []string) {
	for _, line := range lines {
		fmt.Printf("%s", line)
	}
}
