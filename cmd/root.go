package cmd

import (
	"os"
	"strings"

	"github.com/crimsonn/manifest-inspector/internal/hls"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "manifest-inspector",
	Short: "Inspect HLS streaming manifests (testing tool, work in progress)",
	Long: `manifest-inspector is a CLI testing tool for streaming manifests.

It fetches an HLS master playlist and reports validation issues
(variants, codecs, audio/subtitle groups, and related attributes).

DASH/MPD is not supported yet. This project is early and incomplete;
checks and output will change.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		url, _ := cmd.Flags().GetString("url")
		if url != "" {
			lowerURL := strings.ToLower(url)
			if strings.Contains(lowerURL, "m3u8") {
				parser := hls.NewParser()
				err := parser.Parse(url)
				if err != nil {
					return err
				}
			} else if strings.Contains(lowerURL, "mpd") {
				cmd.PrintErrln("warning: DASH/MPD manifests cannot be parsed yet")
			}
		}
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("url", "u", "", "Use the HLS or MPD manifest url")
}
