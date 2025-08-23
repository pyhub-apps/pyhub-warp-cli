package cmd

import (
	"fmt"
	"runtime"
	"strings"

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
			fmt.Fprintln(cmd.OutOrStdout(), i18n.T("version.output.title"))

			// Get the label width for proper alignment
			// We need to calculate the maximum label width for alignment
			labels := []string{
				i18n.T("version.output.version"),
				i18n.T("version.output.commit"),
				i18n.T("version.output.built"),
				i18n.T("version.output.goversion"),
				i18n.T("version.output.platform"),
			}

			// Remove the %s placeholder from labels to get actual label text
			maxWidth := 0
			for _, label := range labels {
				// Find the position of ": %s" and extract just the label part
				if idx := strings.Index(label, ": %s"); idx > 0 {
					labelText := label[:idx]
					if len(labelText) > maxWidth {
						maxWidth = len(labelText)
					}
				} else if idx := strings.Index(label, ": %s/%s"); idx > 0 {
					labelText := label[:idx]
					if len(labelText) > maxWidth {
						maxWidth = len(labelText)
					}
				}
			}

			// Add padding for better alignment
			maxWidth += 2

			// Format and print with proper alignment
			fmt.Fprintf(cmd.OutOrStdout(), "  %-*s %s\n", maxWidth, getLabel(i18n.T("version.output.version")), Version)
			fmt.Fprintf(cmd.OutOrStdout(), "  %-*s %s\n", maxWidth, getLabel(i18n.T("version.output.commit")), GitCommit)
			fmt.Fprintf(cmd.OutOrStdout(), "  %-*s %s\n", maxWidth, getLabel(i18n.T("version.output.built")), BuildDate)
			fmt.Fprintf(cmd.OutOrStdout(), "  %-*s %s\n", maxWidth, getLabel(i18n.T("version.output.goversion")), runtime.Version())
			fmt.Fprintf(cmd.OutOrStdout(), "  %-*s %s/%s\n", maxWidth, getLabel(i18n.T("version.output.platform")), runtime.GOOS, runtime.GOARCH)
		},
	}
}

// getLabel extracts the label part from a translation string
// For example, "Version: %s" returns "Version:"
func getLabel(translationStr string) string {
	if idx := strings.Index(translationStr, ": %s"); idx > 0 {
		return translationStr[:idx] + ":"
	}
	if idx := strings.Index(translationStr, ": %s/%s"); idx > 0 {
		return translationStr[:idx] + ":"
	}
	return translationStr
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
