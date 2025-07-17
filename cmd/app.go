package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/andvarfolomeev/yankr/pkg/config"
	"github.com/andvarfolomeev/yankr/pkg/snippet"
	"github.com/urfave/cli/v2"
)

func BuildApp(cfg *config.Config) *cli.App {
	app := &cli.App{
		Name:  "yankr",
		Usage: "Snippet manager with clipboard integration",
		Description: "A CLI tool for managing and using code snippets with parameterized templates.\n\n" +
			"   ENVIRONMENT VARIABLES:\n" +
			"     YANKR_SNIPPETS_DIR - Override the default snippets directory location\n\n" +
			"   SNIPPET PARAMETERS:\n" +
			"     Parameters in snippets are defined using double curly braces: {{parameter_name}}\n" +
			"     You can provide values via --param flag or will be prompted interactively",
		Commands: []*cli.Command{
			buildListCommand(cfg),
			buildYankCommand(cfg),
			buildCreateCommand(cfg),
			buildPathCommand(cfg),
			buildParamsCommand(cfg),
		},
	}

	return app
}

func buildListCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:        "list",
		Usage:       "List all available snippets",
		Description: "Shows all snippets available in the configured snippets directory.",
		Action: func(c *cli.Context) error {
			snippets, err := snippet.List(cfg.SnippetsDir)
			if err != nil {
				return err
			}
			if len(snippets) == 0 {
				fmt.Println("No snippets found.")
				return nil
			}
			fmt.Println("Available snippets:")
			for _, s := range snippets {
				fmt.Println("-", s)
			}
			return nil
		},
	}
}

func buildYankCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:      "yank",
		Usage:     "Process a snippet and copy it to clipboard",
		ArgsUsage: "SNIPPET_NAME",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "param",
				Aliases: []string{"p"},
				Usage:   "Define parameter values in format 'name=value' (can be used multiple times)",
			},
		},
		Description: "This command processes a snippet, replacing all parameters with their values,\n" +
			"   and copies the result to the system clipboard.\n\n" +
			"   Example: yankr yank email-template --param recipient=John --param subject=Meeting",
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return fmt.Errorf("snippet name is required")
			}

			snippetName := c.Args().First()
			snippetPath := filepath.Join(cfg.SnippetsDir, snippetName)

			if _, err := os.Stat(snippetPath); os.IsNotExist(err) {
				return fmt.Errorf("snippet '%s' not found", snippetName)
			}

			params := make(map[string]string)
			paramSlice := c.StringSlice("param")
			for _, p := range paramSlice {
				parts := strings.SplitN(p, "=", 2)
				if len(parts) == 2 {
					params[parts[0]] = parts[1]
				} else {
					return fmt.Errorf("invalid parameter format: %s (use 'name=value')", p)
				}
			}

			return snippet.Process(snippetPath, params)
		},
	}
}

func buildCreateCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create a new snippet",
		Description: "Creates a new snippet file. You can use {{parameter_name}} syntax in your snippet\n" +
			"   to define parameters that will be replaced when using the 'yank' command.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "Name of the snippet",
			},
		},
		Action: func(c *cli.Context) error {
			name := c.String("name")
			if name == "" {
				return fmt.Errorf("snippet name is required (--name or -n)")
			}

			return snippet.Create(cfg.SnippetsDir, name)
		},
	}
}

func buildPathCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "path",
		Usage: "Show the current snippets directory path",
		Description: "Displays the current path where snippets are stored and indicates\n" +
			"   whether this path is the default or has been overridden by the YANKR_SNIPPETS_DIR environment variable.",
		Action: func(c *cli.Context) error {
			fmt.Printf("Current snippets directory: %s\n", cfg.SnippetsDir)
			return nil
		},
	}
}

func buildParamsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:      "params",
		Usage:     "Show parameters in a snippet",
		ArgsUsage: "SNIPPET_NAME",
		Description: "Displays all parameters defined in a snippet and shows example usage\n" +
			"   of how to provide values for these parameters when using the 'yank' command.",
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return fmt.Errorf("snippet name is required")
			}

			snippetName := c.Args().First()
			snippetPath := filepath.Join(cfg.SnippetsDir, snippetName)

			if _, err := os.Stat(snippetPath); os.IsNotExist(err) {
				return fmt.Errorf("snippet '%s' not found", snippetName)
			}

			params, err := snippet.GetParams(snippetPath)
			if err != nil {
				return err
			}

			if len(params) == 0 {
				fmt.Printf("Snippet '%s' has no parameters.\n", snippetName)
				return nil
			}

			fmt.Printf("Parameters in snippet '%s':\n", snippetName)
			for param := range params {
				fmt.Printf("- %s\n", param)
			}

			fmt.Printf("\nCommand-line usage example:\n")
			var paramList []string
			for param := range params {
				paramList = append(paramList, fmt.Sprintf("--param %s=value", param))
			}
			fmt.Printf("  yankr yank %s %s\n", snippetName, strings.Join(paramList, " "))

			return nil
		},
	}
}
