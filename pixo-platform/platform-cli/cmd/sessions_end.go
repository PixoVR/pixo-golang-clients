/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"encoding/json"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/legacy"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/loader"
	"github.com/spf13/cobra"
	"time"
)

// sessionsEndCmd represents the sessions end command
var sessionsEndCmd = &cobra.Command{
	Use:   "end",
	Short: "End a session",
	Long:  `End a session to mimic headset interactions`,
	RunE: func(cmd *cobra.Command, args []string) error {
		questions := []config.Value{
			{Question: forms.Question{Type: forms.Input, Key: "session-id"}},
			{Question: forms.Question{Type: forms.Input, Key: "score"}},
			{Question: forms.Question{Type: forms.Input, Key: "max-score"}},
		}

		answers, err := Ctx.ConfigManager.GetValuesOrSubmitForm(questions, cmd)
		if err != nil {
			return err
		}

		input := platform.Session{
			ID:        forms.Int(answers["session-id"]),
			RawScore:  float64(forms.Int(answers["score"])),
			MaxScore:  float64(forms.Int(answers["max-score"])),
			Completed: true,
		}

		spinner := loader.NewLoader(cmd.Context(), "Ending session...", Ctx.Printer)

		session, err := Ctx.PlatformClient.UpdateSession(cmd.Context(), input)
		if err != nil {
			Ctx.Printer.Println(":exclamation: Unable to end session: ", err)
			return nil
		}

		sessionDuration, err := time.ParseDuration(session.Duration)
		if err != nil {
			Ctx.Printer.Println(":exclamation: Unable to parse session duration: ", err)
			return nil
		}

		eventInput := struct {
			SessionID *int                 `json:"sessionID"`
			UUID      string               `json:"uuid,omitempty" `
			DeviceID  string               `json:"deviceId,omitempty"`
			Type      string               `json:"eventType,omitempty"`
			Payload   *legacy.EventPayload `json:"jsonData,omitempty"`
		}{
			SessionID: &session.ID,
			Type:      "PIXOVR_SESSION_COMPLETE",
			DeviceID:  session.DeviceID,
			Payload: &legacy.EventPayload{
				LessonStatus:    &[]string{"passed"}[0],
				SessionDuration: sessionDuration.Seconds(),
				Score:           &session.RawScore,
				ScoreMax:        &session.MaxScore,
				ScoreScaled:     &session.ScaledScore,
				Result: &legacy.Result{
					Completion: true,
					Success:    true,
					Duration:   session.Duration,
				},
			},
		}

		eventBytes, err := json.Marshal(eventInput)
		if err != nil {
			Ctx.Printer.Println(":exclamation: Unable to end session: ", err)
			return nil
		}

		_, err = Ctx.PlatformClient.Post("event", eventBytes)
		spinner.Stop()
		if err != nil {
			Ctx.Printer.Println(":exclamation: Unable to end session: ", err)
			return nil
		}

		percentScore := int(session.ScaledScore * 100)

		Ctx.Printer.Println("\n:white_check_mark: Session completed")
		Ctx.Printer.Printf(":input_numbers: Score: %.2f/%.2f\n", session.RawScore, session.MaxScore)
		Ctx.Printer.Printf(":hundred_points: Percent: %d%s\n", percentScore, "%")
		Ctx.Printer.Printf(":hourglass_done: Duration: %s\n", session.Duration)
		return nil
	},
}

func init() {
	sessionsCmd.AddCommand(sessionsEndCmd)

	sessionsEndCmd.Flags().String("session-id", "", "Session ID")
	sessionsEndCmd.Flags().String("score", "", "Score for the session")
	sessionsEndCmd.Flags().String("max-score", "", "Max possible score for the session")
}
