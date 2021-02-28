package cmd

import (
	"encoding/binary"
	"io"
	"os"

	"github.com/spf13/cobra"

	"hz.tools/pulseaudio"
)

// playCmd represents the record command
var playCmd = &cobra.Command{
	Use: "play",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := getConfig(cmd)
		if err != nil {
			return err
		}

		outputPath, err := cmd.Flags().GetString("input")
		if err != nil {
			return err
		}

		fd, err := os.Open(outputPath)
		if err != nil {
			return err
		}
		defer fd.Close()

		writer, err := pulseaudio.NewWriter(config)
		if err != nil {
			return err
		}
		defer writer.Close()

		in := make([]float32, config.Rate/10)
		for {
			if err := binary.Read(fd, binary.LittleEndian, in); err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}
			if err := writer.Write(in); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(playCmd)
	playCmd.Flags().String("input", "audio.raw", "file to read from")
}
