package cmd

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"

	"hz.tools/pulseaudio"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use: "check",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := getConfig(cmd)
		if err != nil {
			return err
		}

		timeDuration, err := cmd.Flags().GetString("duration")
		if err != nil {
			return err
		}

		duration, err := time.ParseDuration(timeDuration)
		if err != nil {
			return err
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, duration)
		defer cancel()

		reader, err := pulseaudio.NewReader(config)
		if err != nil {
			return err
		}
		defer reader.Close()

		writer, err := pulseaudio.NewWriter(config)
		if err != nil {
			return err
		}
		defer writer.Close()
		for {
			if err := reader.Flush(); err != nil {
				return err
			}
			log.Println("Recording")
			out := make([]float32, int(config.Rate*uint(duration.Seconds())))
			if err := reader.Read(out); err != nil {
				return err
			}
			log.Println("Playing")
			if err := writer.Write(out); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().String("duration", "10s", "How long to run for")
}
