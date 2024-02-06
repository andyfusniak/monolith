package main

import (
	"context"
	"fmt"
	"os"

	"github.com/andyfusniak/monolith/internal/cli"

	"github.com/spf13/cobra"
)

var (
	version   string
	gitCommit string
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %+v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// cli application
	cliApp := cli.NewApp(
		cli.WithVersion(version),
		cli.WithGitCommit(gitCommit),
		cli.WithStdOut(os.Stdout),
		cli.WithStdErr(os.Stderr),
	)

	root := cobra.Command{
		Use:     "monolith",
		Short:   "monolith command line tool for the monolith web service",
		Version: cliApp.Version(),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			v := ctx.Value(cli.AppKey("app"))
			_ = v.(*cli.App)
		},
	}
	root.AddCommand(cli.NewCmdInfo())
	root.AddCommand(cli.NewCmdMigrate())
	root.AddCommand(cli.NewCmdServer(version, gitCommit))

	ctx := context.WithValue(context.Background(), cli.AppKey("app"), cliApp)
	if err := root.ExecuteContext(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	return nil

}
