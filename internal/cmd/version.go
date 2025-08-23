package cmd

import (
	"fmt"
	"runtime"

	"github.com/pyhub-kr/pyhub-sejong-cli/internal/i18n"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd *cobra.Command

// initVersionCmd initializes the version command
func initVersionCmd() {
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: i18n.T("version.short"),
		Long:  i18n.T("version.long"),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(i18n.T("version.output.title"))
			fmt.Println(i18n.Tf("version.output.version", Version))
			fmt.Println(i18n.Tf("version.output.commit", GitCommit))
			fmt.Println(i18n.Tf("version.output.built", BuildDate))
			fmt.Println(i18n.Tf("version.output.goversion", runtime.Version()))
			fmt.Println(i18n.Tf("version.output.platform", runtime.GOOS, runtime.GOARCH))
		},
	}
}

// updateVersionCommand updates version command descriptions
func updateVersionCommand() {
	if versionCmd != nil {
		versionCmd.Short = i18n.T("version.short")
		versionCmd.Long = i18n.T("version.long")
	}
}

func init() {
	// Version command will be initialized and added in Execute()
}
