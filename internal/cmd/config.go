package cmd

import (
	"fmt"
	"strings"

	"github.com/pyhub-kr/pyhub-sejong-cli/internal/config"
	cliErrors "github.com/pyhub-kr/pyhub-sejong-cli/internal/errors"
	"github.com/pyhub-kr/pyhub-sejong-cli/internal/i18n"
	"github.com/pyhub-kr/pyhub-sejong-cli/internal/logger"
	"github.com/pyhub-kr/pyhub-sejong-cli/internal/onboarding"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd *cobra.Command

// initConfigCmd initializes the config command
func initConfigCmd() {
	configCmd = &cobra.Command{
		Use:     "config",
		Short:   i18n.T("config.short"),
		Long:    i18n.T("config.long"),
		Example: i18n.T("config.example"),
		RunE: func(cmd *cobra.Command, args []string) error {
			// If no subcommand is provided, show help
			return cmd.Help()
		},
	}
}

// configSetCmd represents the config set command
var configSetCmd *cobra.Command

// initConfigSetCmd initializes the config set command
func initConfigSetCmd() {
	configSetCmd = &cobra.Command{
		Use:     "set <key> <value>",
		Short:   i18n.T("config.set.short"),
		Long:    i18n.T("config.set.long"),
		Example: i18n.T("config.set.example"),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := strings.TrimSpace(args[0])
			value := strings.TrimSpace(args[1])

			// Validate key format
			if !isValidConfigKey(key) {
				return fmt.Errorf(i18n.T("config.set.invalidKey"), key)
			}

			// Validate value is not empty
			if value == "" {
				return fmt.Errorf(i18n.T("config.set.emptyValue"))
			}

			// Special handling for API key
			if key == "law.key" {
				if err := config.SetAPIKey(value); err != nil {
					return fmt.Errorf(i18n.T("config.set.failed"), err)
				}
				guide := onboarding.NewGuideWithWriter(cmd.OutOrStdout(), false)
				guide.ShowSuccess(i18n.T("config.set.apiKeySuccess"))
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", i18n.T("config.path.output"), config.GetConfigPath())
				return nil
			}

			// Generic config set
			config.Set(key, value)
			if err := config.Save(); err != nil {
				return fmt.Errorf(i18n.T("config.set.saveFailed"), err)
			}

			guide := onboarding.NewGuideWithWriter(cmd.OutOrStdout(), false)
			guide.ShowSuccess(fmt.Sprintf(i18n.T("config.set.success"), key, value))
			return nil
		},
	}
}

// configGetCmd represents the config get command
var configGetCmd *cobra.Command

// initConfigGetCmd initializes the config get command
func initConfigGetCmd() {
	configGetCmd = &cobra.Command{
		Use:     "get <key>",
		Short:   i18n.T("config.get.short"),
		Long:    i18n.T("config.get.long"),
		Example: i18n.T("config.get.example"),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := strings.TrimSpace(args[0])

			// Validate key format
			if !isValidConfigKey(key) {
				logger.Error("Invalid config key format: %s", key)
				return cliErrors.New(
					cliErrors.ErrCodeInvalidInput,
					fmt.Sprintf(i18n.T("config.get.invalidKey"), key),
					i18n.T("config.get.apiKeyHelp"),
				)
			}

			// Special handling for API key
			if key == "law.key" {
				if !config.IsAPIKeySet() {
					guide := onboarding.NewGuideWithWriter(cmd.OutOrStdout(), false)
					guide.ShowAPIKeySetup()
					return nil
				}

				apiKey := config.GetAPIKey()
				// Mask API key for security (show first 10 chars only)
				if len(apiKey) > 10 {
					fmt.Fprintf(cmd.OutOrStdout(), "%s: %s...(%d자)\n", key, apiKey[:10], len(apiKey))
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", key, apiKey)
				}
				return nil
			}

			// Generic config get
			value := config.Get(key)
			switch v := value.(type) {
			case nil:
				fmt.Fprintf(cmd.OutOrStdout(), "❌ %s\n", fmt.Sprintf(i18n.T("config.get.notFound"), key))
				return nil
			case string:
				if strings.TrimSpace(v) == "" {
					fmt.Fprintf(cmd.OutOrStdout(), "❌ %s\n", fmt.Sprintf(i18n.T("config.get.notFound"), key))
					return nil
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", key, v)
			default:
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %v\n", key, v)
			}
			return nil
		},
	}
}

// configPathCmd represents the config path command
var configPathCmd *cobra.Command

// initConfigPathCmd initializes the config path command
func initConfigPathCmd() {
	configPathCmd = &cobra.Command{
		Use:   "path",
		Short: i18n.T("config.path.short"),
		Long:  i18n.T("config.path.long"),
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", fmt.Sprintf(i18n.T("config.path.output"), config.GetConfigPath()))
			return nil
		},
	}
}

// isValidConfigKey validates the configuration key format
func isValidConfigKey(key string) bool {
	// Currently only support law.key
	// Can be extended for more keys in the future
	validKeys := []string{
		"law.key",
	}

	for _, validKey := range validKeys {
		if key == validKey {
			return true
		}
		// Check if it's a prefix match for nested keys
		if strings.HasPrefix(key, validKey+".") {
			return true
		}
	}
	return false
}

// updateConfigCommand updates config command descriptions
func updateConfigCommand() {
	if configCmd != nil {
		configCmd.Short = i18n.T("config.short")
		configCmd.Long = i18n.T("config.long")
		configCmd.Example = i18n.T("config.example")
	}
	if configSetCmd != nil {
		configSetCmd.Short = i18n.T("config.set.short")
		configSetCmd.Long = i18n.T("config.set.long")
		configSetCmd.Example = i18n.T("config.set.example")
	}
	if configGetCmd != nil {
		configGetCmd.Short = i18n.T("config.get.short")
		configGetCmd.Long = i18n.T("config.get.long")
		configGetCmd.Example = i18n.T("config.get.example")
	}
	if configPathCmd != nil {
		configPathCmd.Short = i18n.T("config.path.short")
		configPathCmd.Long = i18n.T("config.path.long")
	}
}

func init() {
	// Config commands will be initialized and added in Execute()
}
