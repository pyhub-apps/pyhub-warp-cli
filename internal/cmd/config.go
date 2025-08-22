package cmd

import (
	"fmt"
	"strings"

	"github.com/pyhub-kr/pyhub-sejong-cli/internal/config"
	"github.com/pyhub-kr/pyhub-sejong-cli/internal/onboarding"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "설정 관리",
	Long: `Sejong CLI의 설정을 관리합니다.

API 키와 같은 설정값을 저장하고 조회할 수 있습니다.
설정 파일 위치는 'sejong config path' 명령으로 확인할 수 있습니다.`,
	Example: `  # API 키 설정
  sejong config set law.key YOUR_API_KEY
  
  # API 키 확인
  sejong config get law.key
  
  # 설정 파일 경로 확인
  sejong config path`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no subcommand is provided, show help
		return cmd.Help()
	},
}

// configSetCmd represents the config set command
var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "설정값 저장",
	Long:  `지정한 키에 값을 저장합니다.`,
	Example: `  # API 키 설정
  sejong config set law.key YOUR_API_KEY`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := strings.TrimSpace(args[0])
		value := strings.TrimSpace(args[1])

		// Validate key format
		if !isValidConfigKey(key) {
			return fmt.Errorf("잘못된 설정 키 형식: %s (허용: law.key)", key)
		}
		
		// Validate value is not empty
		if value == "" {
			return fmt.Errorf("설정값이 비어있습니다")
		}

		// Special handling for API key
		if key == "law.key" {
			if err := config.SetAPIKey(value); err != nil {
				return fmt.Errorf("API 키 설정 실패: %w", err)
			}
			guide := onboarding.NewGuide()
			guide.ShowSuccess("API 키가 성공적으로 설정되었습니다")
			fmt.Printf("설정 파일: %s\n", config.GetConfigPath())
			return nil
		}

		// Generic config set
		config.Set(key, value)
		if err := config.Save(); err != nil {
			return fmt.Errorf("설정 저장 실패: %w", err)
		}

		guide := onboarding.NewGuide()
		guide.ShowSuccess(fmt.Sprintf("설정이 저장되었습니다: %s = %s", key, value))
		return nil
	},
}

// configGetCmd represents the config get command
var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "설정값 조회",
	Long:  `지정한 키의 값을 조회합니다.`,
	Example: `  # API 키 확인
  sejong config get law.key`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := strings.TrimSpace(args[0])

		// Validate key format
		if !isValidConfigKey(key) {
			return fmt.Errorf("잘못된 설정 키 형식: %s (허용: law.key)", key)
		}

		// Special handling for API key
		if key == "law.key" {
			if !config.IsAPIKeySet() {
				guide := onboarding.NewGuide()
				guide.ShowAPIKeySetup()
				return nil
			}
			
			apiKey := config.GetAPIKey()
			// Mask API key for security (show first 10 chars only)
			if len(apiKey) > 10 {
				fmt.Printf("%s: %s...(%d자)\n", key, apiKey[:10], len(apiKey))
			} else {
				fmt.Printf("%s: %s\n", key, apiKey)
			}
			return nil
		}

		// Generic config get
		value := config.Get(key)
		switch v := value.(type) {
		case nil:
			fmt.Printf("❌ 설정값이 없습니다: %s\n", key)
			return nil
		case string:
			if strings.TrimSpace(v) == "" {
				fmt.Printf("❌ 설정값이 없습니다: %s\n", key)
				return nil
			}
			fmt.Printf("%s: %s\n", key, v)
		default:
			fmt.Printf("%s: %v\n", key, v)
		}
		return nil
	},
}

// configPathCmd represents the config path command
var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "설정 파일 경로 확인",
	Long:  `설정 파일의 경로를 확인합니다.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("설정 파일 경로: %s\n", config.GetConfigPath())
		return nil
	},
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

func init() {
	// Add subcommands to config
	// Avoid printing usage on handled errors
	configCmd.SilenceUsage = true
	configCmd.SilenceErrors = true
	configSetCmd.SilenceUsage = true
	configSetCmd.SilenceErrors = true
	configGetCmd.SilenceUsage = true
	configGetCmd.SilenceErrors = true
	configPathCmd.SilenceUsage = true
	configPathCmd.SilenceErrors = true
	
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configPathCmd)

	// Add config command to root
	rootCmd.AddCommand(configCmd)
}