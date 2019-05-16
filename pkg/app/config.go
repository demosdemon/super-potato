package app

import "github.com/spf13/cobra"

type Config interface {
	Use() string
	Args(cmd *cobra.Command, args []string) error
}

type PreRunner interface {
	Config
	PreRun(cmd *cobra.Command, args []string) error
}

type Runner interface {
	Config
	Run(cmd *cobra.Command, args []string) error
}

type PostRunner interface {
	Config
	PostRun(cmd *cobra.Command, args []string) error
}

type RootRunner interface {
	MasterRunner
	PersistentPreRun(cmd *cobra.Command, args []string) error
	PersistentPostRun(cmd *cobra.Command, args []string) error
}

type MasterRunner interface {
	Config
	SubCommands() []Config
}
