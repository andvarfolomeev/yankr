```
██    ██  █████  ███    ██ ██   ██ ██████
 ██  ██  ██   ██ ████   ██ ██  ██  ██   ██
  ████   ███████ ██ ██  ██ █████   ██████
   ██    ██   ██ ██  ██ ██ ██  ██  ██   ██
   ██    ██   ██ ██   ████ ██   ██ ██   ██
             — snip. fill. yank. done.
```

# Installing

```bash
go install github.com/andvarfolomeev/yankr
```

# Using

```
NAME:
   yankr - Snippet manager with clipboard integration

USAGE:
   yankr [global options] command [command options] [arguments...]

DESCRIPTION:
   A CLI tool for managing and using code snippets with parameterized templates.

      ENVIRONMENT VARIABLES:
        YANKR_SNIPPETS_DIR - Override the default snippets directory location

      SNIPPET PARAMETERS:
        Parameters in snippets are defined using double curly braces: {{parameter_name}}
        You can provide values via --param flag or will be prompted interactively

COMMANDS:
   list     List all available snippets
   yank     Process a snippet and copy it to clipboard
   create   Create a new snippet
   path     Show the current snippets directory path
   params   Show parameters in a snippet
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

# License

MIT
