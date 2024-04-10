/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	graphql_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/loader"
	"github.com/spf13/cobra"
)

// sessionsEndCmd represents the sessions end command
var sessionsEndCmd = &cobra.Command{
	Use:   "end",
	Short: "End a session",
	Long:  `End a session to mimic headset interactions`,
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionID, ok := Ctx.ConfigManager.GetIntConfigValueOrAskUser("session-id", cmd)
		if !ok {
			Ctx.ConfigManager.Println(":exclamation: Session ID not provided")
		}

		score, _ := Ctx.ConfigManager.GetIntConfigValueOrAskUser("score", cmd)
		maxScore, _ := Ctx.ConfigManager.GetIntConfigValueOrAskUser("max-score", cmd)

		input := graphql_api.Session{
			ID:        sessionID,
			RawScore:  float64(score),
			MaxScore:  float64(maxScore),
			Completed: true,
		}

		spinner := loader.NewLoader(cmd.Context(), "Ending session...", Ctx.ConfigManager)
		session, err := Ctx.PlatformClient.UpdateSession(cmd.Context(), input)
		spinner.Stop()
		if err != nil {
			Ctx.ConfigManager.Println(":exclamation: Unable to end session: ", err)
			return nil
		}

		percentScore := int(session.ScaledScore * 100)

		Ctx.ConfigManager.Printf(":white_check_mark: Session completed with score %d/%d - %d%s", score, maxScore, percentScore, "%")
		return nil
	},
}

func init() {
	sessionsCmd.AddCommand(sessionsEndCmd)

	sessionsEndCmd.Flags().Int("session-id", 0, "Session ID")
	sessionsEndCmd.Flags().Int("score", 0, "Score for the session")
	sessionsEndCmd.Flags().Int("max-score", 0, "Max possible score for the session")
}
