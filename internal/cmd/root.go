package cmd

import (
	"fmt"
	"os"

	"github.com/pyhub-apps/sejong-cli/internal/config"
	"github.com/pyhub-apps/sejong-cli/internal/i18n"
	"github.com/pyhub-apps/sejong-cli/internal/logger"
	"github.com/spf13/cobra"
)

var (
	// Version information - will be set during build
	// Following LINE HeadVer versioning: {head}.{yearweek}.{build}
	// NOTE: real value should be injected via -ldflags; keep a neutral default.
	Version   = "dev"
	BuildDate = "unknown"
	GitCommit = "unknown"

	// Language flag
	langFlag string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd *cobra.Command

// initRootCmd initializes the root command with i18n support
func initRootCmd() {
	rootCmd = &cobra.Command{
		Use:   "sejong",
		Short: i18n.T("cli.short"),
		Long:  i18n.T("cli.long"),
		Example: `  # 법령 검색
  sejong law "개인정보 보호법"
  
  # JSON 형식으로 출력
  sejong law "도로교통법" --format json
  
  # API 키 설정
  sejong config set law.key YOUR_API_KEY`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Set language if specified
			if langFlag != "" {
				i18n.SetLanguage(langFlag)
				// Re-initialize commands with new language
				updateCommandDescriptions()
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// If no subcommand is provided, show help
			return cmd.Help()
		},
	}
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	// Initialize i18n first
	if err := i18n.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize i18n: %v\n", err)
	}

	// Initialize root command with i18n support
	initRootCmd()
	setupFlags()

	// Initialize and add subcommands
	initVersionCmd()
	initConfigCmd()
	initConfigSetCmd()
	initConfigGetCmd()
	initConfigPathCmd()
	initLawCmd()
	initOrdinanceCmd()
	initPrecedentCmd()
	initAdmruleCmd()
	initInterpretationCmd()

	// Add version command to root
	rootCmd.AddCommand(versionCmd)

	// Add config command and its subcommands
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
	rootCmd.AddCommand(configCmd)

	// Add law command to root
	rootCmd.AddCommand(lawCmd)

	// Add ordinance command to root
	rootCmd.AddCommand(ordinanceCmd)

	// Add precedent command to root
	rootCmd.AddCommand(precedentCmd)

	// Add administrative rule command to root
	rootCmd.AddCommand(admruleCmd)

	// Add legal interpretation command to root
	rootCmd.AddCommand(interpretationCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func setupFlags() {
	// Initialize configuration
	cobra.OnInitialize(initConfig)

	// Language flag
	rootCmd.PersistentFlags().StringVar(&langFlag, "lang", "", "Language (ko, en)")

	// Global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, i18n.T("cli.verbose"))

	// Version flag
	rootCmd.Version = fmt.Sprintf("%s (built %s, commit %s)", Version, BuildDate, GitCommit)
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "version %s" .Version}}
`)
}

// updateCommandDescriptions updates command descriptions after language change
func updateCommandDescriptions() {
	rootCmd.Short = i18n.T("cli.short")
	rootCmd.Long = i18n.T("cli.long")

	// Update verbose flag description
	if flag := rootCmd.PersistentFlags().Lookup("verbose"); flag != nil {
		flag.Usage = i18n.T("cli.verbose")
	}

	// Update subcommands (these will be updated in their respective files)
	updateVersionCommand()
	updateConfigCommand()
	updateLawCommand()
	updateOrdinanceCommand()
	updatePrecedentCommand()
	updateAdmruleCommand()
	updateInterpretationCommand()
}

func init() {
	// Commands will be added in their respective init functions
}

// initConfig initializes the configuration
func initConfig() {
	// Set up logging based on verbose flag
	if verbose, _ := rootCmd.PersistentFlags().GetBool("verbose"); verbose {
		logger.SetVerbose(true)
	}

	if err := config.Initialize(); err != nil {
		logger.Warn("Failed to initialize config: %v", err)
	}
}

// SetVersionInfo sets the version information for the CLI
func SetVersionInfo(version, commit, date string) {
	Version = version
	GitCommit = commit
	BuildDate = date
	// rootCmd will be initialized later in Execute(), so we don't set Version here
}
