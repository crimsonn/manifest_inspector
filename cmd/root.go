package cmd

import (
	"os"
	"strings"

	"github.com/crimsonn/manifest-inspector/internal/hls"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "manifest-inspector",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		url, _ := cmd.Flags().GetString("url")
		if url != "" {
			isHLS := strings.Contains(url, "m3u8")
			if isHLS {
				parser := hls.NewParser()
				err := parser.Parse(url)
				if err != nil {
					return err
				}
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
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringP("url", "u", "", "Use the HLS or MPD manifest url")
}
