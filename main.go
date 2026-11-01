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
