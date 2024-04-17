/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"encoding/json"
	graphql_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/loader"
	primary_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/spf13/cobra"
	"time"
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

		ipAddress, _ := Ctx.PlatformClient.GetIPAddress()
		input := graphql_api.Session{
			ID:        sessionID,
			RawScore:  float64(score),
			MaxScore:  float64(maxScore),
			Completed: true,
			IPAddress: ipAddress,
		}

		spinner := loader.NewLoader(cmd.Context(), "Ending session...", Ctx.ConfigManager)

		session, err := Ctx.PlatformClient.UpdateSession(cmd.Context(), input)
		if err != nil {
			Ctx.ConfigManager.Println(":exclamation: Unable to end session: ", err)
			return nil
		}

		sessionDuration, err := time.ParseDuration(session.Duration)
		if err != nil {
			Ctx.ConfigManager.Println(":exclamation: Unable to parse session duration: ", err)
			return nil
		}

		eventInput := struct {
			SessionID *int                   `json:"sessionID"`
			IP        string                 `json:"ipAddress,omitempty"`
			JSONData  *primary_api.JSONEvent `json:"jsonData,omitempty"`
			DeviceID  string                 `json:"deviceId,omitempty"`
			UUID      string                 `json:"uuid,omitempty" `
			EventType string                 `json:"eventType,omitempty"`
			UserID    int                    `json:"user_id,omitempty"`
			OrgID     int                    `json:"org_id,omitempty"`
			ModuleID  int                    `json:"moduleId,omitempty"`
		}{
			SessionID: &sessionID,
			EventType: "PIXOVR_SESSION_COMPLETE",
			IP:        session.IPAddress,
			DeviceID:  session.DeviceID,
			UUID:      session.UUID,
			UserID:    session.UserID,
			OrgID:     session.OrgID,
			ModuleID:  session.ModuleID,
			JSONData: &primary_api.JSONEvent{
				SessionDuration: sessionDuration.Seconds(),
				Score:           &session.RawScore,
				ScoreMax:        &session.MaxScore,
				ScoreScaled:     &session.ScaledScore,
			},
		}

		eventBytes, err := json.Marshal(eventInput)
		if err != nil {
			Ctx.ConfigManager.Println(":exclamation: Unable to end session: ", err)
			return nil
		}

		_, err = Ctx.PlatformClient.Post("event", eventBytes)
		spinner.Stop()
		if err != nil {
			Ctx.ConfigManager.Println(":exclamation: Unable to end session: ", err)
			return nil
		}

		percentScore := int(session.ScaledScore * 100)

		Ctx.ConfigManager.Println("\n:white_check_mark: Session completed")
		Ctx.ConfigManager.Printf(":input_numbers: Score: %d/%d\n", score, maxScore)
		Ctx.ConfigManager.Printf(":hundred_points: Percent: %d%s\n", percentScore, "%")
		Ctx.ConfigManager.Printf(":hourglass_done: Duration: %s\n", session.Duration)
		return nil
	},
}

func init() {
	sessionsCmd.AddCommand(sessionsEndCmd)

	sessionsEndCmd.Flags().Int("session-id", 0, "Session ID")
	sessionsEndCmd.Flags().Int("score", 0, "Score for the session")
	sessionsEndCmd.Flags().Int("max-score", 0, "Max possible score for the session")
}
