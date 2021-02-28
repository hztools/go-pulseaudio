package cmd

import (
	"context"
	"encoding/binary"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"

	"hz.tools/pulseaudio"
)

// recordCmd represents the record command
var recordCmd = &cobra.Command{
	Use: "record",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := getConfig(cmd)
		if err != nil {
			return err
		}

		outputPath, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}

		fd, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer fd.Close()

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

		log.Println("Recording")
		out := make([]float32, config.Rate/10)
		for {
			if err := reader.Read(out); err != nil {
				return err
			}
			binary.Write(fd, binary.LittleEndian, out)
			if err := ctx.Err(); err != nil {
				if err == context.DeadlineExceeded {
					break
				}
				return err
			}
		}
		log.Println("Done!")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(recordCmd)

	recordCmd.Flags().String("duration", "10s", "How long to run for")
	recordCmd.Flags().String("output", "audio.raw", "file to write to")
}
