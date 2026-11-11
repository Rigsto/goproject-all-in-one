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

	"golanganalyzer/config"
	"golanganalyzer/flatten"
	"golanganalyzer/minify"
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
	root.PersistentFlags().StringVar(&configPath, "config", "config.yaml", "path to config.yaml")

	root.AddCommand(newFlattenCmd(&configPath))
	root.AddCommand(newMinifyCmd(&configPath))

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newFlattenCmd(configPath *string) *cobra.Command {
	var source string
	var out string

	cmd := &cobra.Command{
		Use:   "flatten",
		Short: "Write a full-source flatten dump for AI analysis",
		Long: `Walk --source and write a full-source dump.

Each file is introduced with:
  ===== FILE <relpath> | LINES 1-<n>

The body is unchanged (aside from normalizing newlines). Dump body line N is
original source line N, so AI citations need no map. Pair with prompts under
prompt/flatten/.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if source == "" {
				return fmt.Errorf("--source is required")
			}
			if out == "" {
				out = filepath.Join("result", "flatten.txt")
			}
			cfg, err := config.Load(*configPath)
			if err != nil {
				return err
			}
			if err := flatten.FlattenWithConfig(source, out, cfg); err != nil {
				return err
			}
			fmt.Printf("wrote %s\n", out)
			return nil
		},
	}

	cmd.Flags().StringVar(&source, "source", "", "source project directory")
	cmd.Flags().StringVar(&out, "out", "", "output file (default: result/flatten.txt)")
	_ = cmd.MarkFlagRequired("source")
	return cmd
}

func newMinifyCmd(configPath *string) *cobra.Command {
	var source string
	var out string

	cmd := &cobra.Command{
		Use:   "minify",
		Short: "Write a minified dump with line maps for AI analysis",
		Long: `Walk --source and write a comment-stripped dump with line maps.

Each file is introduced with:
  @FILE <relpath>
  @ORIG 1-<n>
  @MAP 1>3,2>4,...

@MAP out>orig maps each kept dump line to its original source line. Pair with
prompts under prompt/minify/ so citations use original lines.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if source == "" {
				return fmt.Errorf("--source is required")
			}
			if out == "" {
				out = filepath.Join("result", "minify.txt")
			}
			cfg, err := config.Load(*configPath)
			if err != nil {
				return err
			}
			if err := minify.MinifyWithConfig(source, out, cfg); err != nil {
				return err
			}
			fmt.Printf("wrote %s\n", out)
			return nil
		},
	}

	cmd.Flags().StringVar(&source, "source", "", "source project directory")
	cmd.Flags().StringVar(&out, "out", "", "output file (default: result/minify.txt)")
	_ = cmd.MarkFlagRequired("source")
	return cmd
}
