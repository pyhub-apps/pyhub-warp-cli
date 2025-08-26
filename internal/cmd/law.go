package cmd

import (
	"context"
	"errors"
	"strings"

	"github.com/pyhub-apps/pyhub-warp-cli/internal/api"
	cliErrors "github.com/pyhub-apps/pyhub-warp-cli/internal/errors"
	"github.com/pyhub-apps/pyhub-warp-cli/internal/i18n"
	"github.com/pyhub-apps/pyhub-warp-cli/internal/logger"
	"github.com/pyhub-apps/pyhub-warp-cli/internal/onboarding"
	"github.com/spf13/cobra"
)

var (
	outputFormat string
	pageNo       int
	pageSize     int
	sourceFlag   string // "all", "nlic", "elis"

	// testAPIClient allows injecting a mock client for testing
	testAPIClient APIClient
)

// lawCmd represents the law command
var lawCmd *cobra.Command

// initLawCmd initializes the law command
func initLawCmd() {
	lawCmd = &cobra.Command{
		Use:   "law",
		Short: i18n.T("law.short"),
		Long:  i18n.T("law.long"),
		Example: `  # 법령 검색
  warp law search "개인정보 보호법"
  warp law "개인정보 보호법"  # search는 생략 가능
  
  # 법령 상세 조회
  warp law detail 001234
  
  # 법령 이력 조회
  warp law history 001234`,
		// Run default search when args provided without subcommand
		RunE: func(cmd *cobra.Command, args []string) error {
			// If args are provided without subcommand, run search
			if len(args) > 0 {
				return runLawCommand(cmd, args)
			}
			// Otherwise show help
			return cmd.Help()
		},
	}

	// Initialize subcommands
	initLawSearchCmd()
	initLawDetailCmd()
	initLawHistoryCmd()

	// Add subcommands
	lawCmd.AddCommand(lawSearchCmd)
	lawCmd.AddCommand(lawDetailCmd)
	lawCmd.AddCommand(lawHistoryCmd)

	// Flags for backward compatibility (when using law without subcommand)
	lawCmd.Flags().StringVarP(&outputFormat, "format", "f", "table", i18n.T("law.flag.format"))
	lawCmd.Flags().IntVarP(&pageNo, "page", "p", 1, i18n.T("law.flag.page"))
	lawCmd.Flags().IntVarP(&pageSize, "size", "s", 50, i18n.T("law.flag.size"))
	lawCmd.Flags().StringVar(&sourceFlag, "source", "nlic", i18n.T("law.flag.source"))
}

// updateLawCommand updates law command descriptions
func updateLawCommand() {
	if lawCmd != nil {
		lawCmd.Short = i18n.T("law.short")
		lawCmd.Long = i18n.T("law.long")

		// Update flag descriptions
		if flag := lawCmd.Flags().Lookup("format"); flag != nil {
			flag.Usage = i18n.T("law.flag.format")
		}
		if flag := lawCmd.Flags().Lookup("page"); flag != nil {
			flag.Usage = i18n.T("law.flag.page")
		}
		if flag := lawCmd.Flags().Lookup("size"); flag != nil {
			flag.Usage = i18n.T("law.flag.size")
		}

		// Update subcommands
		updateLawSearchCommand()
		updateLawDetailCommand()
		updateLawHistoryCommand()
	}
}

func init() {
	// Law command will be initialized and added in Execute()
}

// APIClient interface for dependency injection and testing
type APIClient interface {
	Search(ctx context.Context, req *api.UnifiedSearchRequest) (*api.SearchResponse, error)
}

func runLawCommand(cmd *cobra.Command, args []string) error {
	// Get search query
	query := strings.TrimSpace(args[0])
	if query == "" {
		logger.Debug("Empty query provided")
		return cliErrors.ErrEmptyQuery
	}

	logger.Debug("Starting law search for query: %s", query)

	// Use test client if available (for testing)
	var client APIClient
	if testAPIClient != nil {
		client = testAPIClient
	} else {
		// Determine which API to use based on source flag
		var apiType api.APIType
		switch sourceFlag {
		case "all":
			apiType = api.APITypeAll
		case "elis":
			apiType = api.APITypeELIS
		case "nlic":
			fallthrough
		default:
			apiType = api.APITypeNLIC
		}

		// Create API client using the factory
		apiClient, err := api.CreateClient(apiType)
		if err != nil {
			// Check if it's an API key error (either CLIError or regular error with API key message)
			var cliErr *cliErrors.CLIError
			if errors.As(err, &cliErr) && cliErr.Code == cliErrors.ErrCodeNoAPIKey {
				guide := onboarding.NewGuideWithWriter(cmd.OutOrStdout(), false)
				guide.ShowAPIKeySetup()
				return nil // Return nil to avoid printing the error twice
			}

			// Also check for direct API key error message from factory
			if strings.Contains(err.Error(), "API 키가 설정되지 않았습니다") {
				guide := onboarding.NewGuideWithWriter(cmd.OutOrStdout(), false)
				guide.ShowAPIKeySetup()
				return nil // Return nil to avoid printing the error twice
			}

			verbose, _ := cmd.Flags().GetBool("verbose")
			logger.LogError(err, verbose)
			return err
		}
		client = apiClient
	}

	// Get verbose flag
	verbose, _ := cmd.Flags().GetBool("verbose")

	// Use searchLaws for the actual search logic
	return searchLaws(client, query, outputFormat, pageNo, pageSize, cmd.OutOrStdout(), verbose)
}
