/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"context"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(modulesCmd)
	modulesCmd.PersistentFlags().StringP("module", "m", "", "Module Abbreviation")
}

func moduleQuestion() forms.Question {
	return forms.Question{
		Type: forms.SelectID,
		Key:  "module",
		LabelFunc: func(i interface{}) string {
			return i.(platform.Module).Abbreviation
		},
		GetItemsFunc: func(ctx context.Context) (interface{}, error) {
			return Ctx.PlatformClient.GetModules(ctx)
		},
	}
}

func modulePlatformQuestion() forms.Question {
	return forms.Question{
		Type: forms.MultiSelectIDs,
		Key:  "platforms",
		GetItemsFunc: func(ctx context.Context) (interface{}, error) {
			return Ctx.PlatformClient.GetPlatforms(ctx)
		},
	}
}

func moduleControlQuestion() forms.Question {
	return forms.Question{
		Type: forms.MultiSelectIDs,
		Key:  "controls",
		GetItemsFunc: func(ctx context.Context) (interface{}, error) {
			return Ctx.PlatformClient.GetControlTypes(ctx)
		},
	}
}

// modulesCmd represents the modules rootCmd
var modulesCmd = &cobra.Command{
	Use:   "modules",
	Short: "Manage modules",
	Long: `Manage modules and their versions.
`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}
