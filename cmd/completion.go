package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

  $ source <(ttrack completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ ttrack completion bash > /etc/bash_completion.d/ttrack
  # macOS:
  $ ttrack completion bash > /usr/local/etc/bash_completion.d/ttrack

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ ttrack completion zsh > "${fpath[1]}/_ttrack"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ ttrack completion fish | source

  # To load completions for each session, execute once:
  $ ttrack completion fish > ~/.config/fish/completions/ttrack.fish

PowerShell:

  PS> ttrack completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> ttrack completion powershell > ttrack.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			if err := cmd.Root().GenBashCompletion(os.Stdout); err != nil {
				log.Fatal(err)
			}
		case "zsh":
			if err := cmd.Root().GenZshCompletion(os.Stdout); err != nil {
				log.Fatal(err)
			}
		case "fish":
			if err := cmd.Root().GenFishCompletion(os.Stdout, true); err != nil {
				log.Fatal(err)
			}
		case "powershell":
			if err := cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout); err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
