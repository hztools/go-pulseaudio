package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"hz.tools/pulseaudio"
)

const (
	appName    = "kc3nwj"
	streamName = "pa-raw"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pa-raw",
	Short: "Preform raw read and write actions with pulseaudio",
}

// Execute will run the rootCmd and exit.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getConfig(cmd *cobra.Command) (pulseaudio.Config, error) {
	pflags := cmd.Flags()
	sps, err := pflags.GetUint("samples-per-second")
	if err != nil {
		return pulseaudio.Config{}, err
	}

	formatString, err := pflags.GetString("format")
	if err != nil {
		return pulseaudio.Config{}, err
	}
	var format pulseaudio.SampleFormat
	switch formatString {
	case "f32ne":
		format = pulseaudio.SampleFormatFloat32NE
	default:
		return pulseaudio.Config{}, fmt.Errorf("pa-raw: unknown format type")

	}

	return pulseaudio.Config{
		Format:     format,
		Rate:       sps,
		AppName:    appName,
		StreamName: streamName,
		Channels:   1,
	}, nil
}

func init() {
	pflags := rootCmd.PersistentFlags()

	pflags.Uint("samples-per-second", 44100, "samples per second")
	pflags.String("format", "f32ne", "[f32ne]")
}
