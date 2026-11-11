// Command golanganalyzer packs a source tree into one text dump for AI review.
//
// Subcommands:
//
//	flatten — full-source dump; body line N matches original line N
//	minify  — comment/blank stripped dump with @MAP back to original lines
//
// Both commands load extensions and ignore rules from --config (default config.yaml),
// walk --source, and write --out (defaults under result/). Pair the dump with a
// matching prompt under prompt/flatten/ or prompt/minify/.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Rigsto/goproject-all-in-one/config"
	"github.com/Rigsto/goproject-all-in-one/flatten"
	"github.com/Rigsto/goproject-all-in-one/minify"
)

func main() {
	root := &cobra.Command{
		Use:   "golanganalyzer",
		Short: "Flatten or minify a project tree for AI code analysis",
		Long: `GolangAnalyzer walks a project directory, keeps configured file extensions,
skips ignore rules, and writes one text dump you can paste into an AI chat.

Use "flatten" for exact source and simple line citations, or "minify" for a
smaller dump with @MAP line maps. Then open a prompt from prompt/flatten/ or
prompt/minify/ that matches the command you ran.`,
	}

	var configPath string
	root.PersistentFlags().StringVar(&configPath, "config", config.DefaultPath, "path to config.yaml")

	root.AddCommand(newDumpCmd(
		&configPath,
		"flatten",
		"Write a full-source flatten dump for AI analysis",
		`Walk --source and write a full-source dump.

Each file is introduced with:
  ===== FILE <relpath> | LINES 1-<n>

The body is unchanged (aside from normalizing newlines). Dump body line N is
original source line N, so AI citations need no map. Pair with prompts under
prompt/flatten/.`,
		filepath.Join("result", "flatten.txt"),
		flatten.FlattenWithConfig,
	))
	root.AddCommand(newDumpCmd(
		&configPath,
		"minify",
		"Write a minified dump with line maps for AI analysis",
		`Walk --source and write a comment-stripped dump with line maps.

Each file is introduced with:
  @FILE <relpath>
  @ORIG 1-<n>
  @MAP 1>3,2>4,...

@MAP out>orig maps each kept dump line to its original source line. Pair with
prompts under prompt/minify/ so citations use original lines.`,
		filepath.Join("result", "minify.txt"),
		minify.MinifyWithConfig,
	))

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newDumpCmd(
	configPath *string,
	name, short, long, defaultOut string,
	run func(source, out string, cfg *config.Config) error,
) *cobra.Command {
	var source string
	var out string

	cmd := &cobra.Command{
		Use:   name,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			if out == "" {
				out = defaultOut
			}
			cfg, err := config.Load(*configPath)
			if err != nil {
				return err
			}
			if err := run(source, out, cfg); err != nil {
				return err
			}
			fmt.Printf("wrote %s\n", out)
			return nil
		},
	}

	cmd.Flags().StringVar(&source, "source", "", "source project directory")
	cmd.Flags().StringVar(&out, "out", "", fmt.Sprintf("output file (default: %s)", defaultOut))
	_ = cmd.MarkFlagRequired("source")
	return cmd
}
