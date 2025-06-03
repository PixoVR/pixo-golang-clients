/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"encoding/json"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/headset"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/loader"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

var useLegacyAPI bool

// sessionsSimulateCmd represents the sessions start command
var sessionsSimulateCmd = &cobra.Command{
	Use:   "simulate",
	Short: "Simulate a session in headset",
	Long:  `start a session, create events, and end the session to mimic headset interactions`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		questions := []config.Value{
			{Question: moduleQuestion()},
			{Question: forms.Question{
				Type: forms.Select,
				Key:  "mode",
				Options: []forms.Option{
					{Label: "Tutorial", Value: "tutorial"},
					{Label: "Practice", Value: "practice"},
					{Label: "Challenge", Value: "challenge"},
				},
				Optional: true,
			}},
			{Question: forms.Question{
				Type:     forms.Input,
				Key:      "scenario",
				Optional: true,
			}},
			{Question: forms.Question{
				Type:     forms.Input,
				Key:      "focus",
				Optional: true,
			}},
			{Question: forms.Question{
				Type:     forms.Input,
				Key:      "specialization",
				Optional: true,
			}},
		}

		answers, err := Ctx.ConfigManager.GetValuesOrSubmitForm(questions, cmd)
		if err != nil {
			return err
		}

		moduleID := forms.Int(answers["module"])
		mode := strings.ToLower(forms.String(answers["mode"]))
		scenario := forms.String(answers["scenario"])
		focus := forms.String(answers["focus"])
		specialization := forms.String(answers["specialization"])

		var sessionID int

		spinner := loader.NewLoader(cmd.Context(), "Starting session...", Ctx.Printer)
		if useLegacyAPI {
			eventRequest := headset.EventRequest{
				ModuleID: moduleID,
				Payload: map[string]interface{}{
					"context": map[string]interface{}{
						"revision": "1.0.0",
					},
				},
			}

			if mode != "" {
				eventRequest.Payload["https://pixovr.com/xapi/extension/sessionMode"] = mode
			}
			if focus != "" {
				eventRequest.Payload["https://pixovr.com/xapi/extension/sessionFocus"] = focus
			}
			if specialization != "" {
				eventRequest.Payload["https://pixovr.com/xapi/extension/sessionSpecialization"] = specialization
			}

			res, err := Ctx.HeadsetClient.StartSession(cmd.Context(), eventRequest)
			spinner.Stop()
			if err != nil {
				return err
			}
			sessionID = *res.Event.SessionID
			Ctx.Printer.Println(":white_check_mark: Session started using legacy headset API")
		} else {
			session := &platform.Session{
				ModuleID:       moduleID,
				Scenario:       scenario,
				Mode:           mode,
				Focus:          focus,
				Specialization: specialization,
			}
			err = Ctx.PlatformClient.CreateSession(cmd.Context(), session)
			spinner.Stop()
			if err != nil {
				return err
			}
			sessionID = session.ID
			Ctx.Printer.Printf(":white_check_mark: Session started for module %s\n", session.Module.Abbreviation)
		}

		startTime := time.Now()

		Ctx.ConfigManager.SetIntConfigValue("session-id", sessionID)

		askToCreateEvent := true
		for askToCreateEvent {
			questions = []config.Value{
				{Question: forms.Question{
					Type:   forms.Confirm,
					Key:    "create-event",
					Prompt: "Create event?",
				}},
			}

			if answers, err = Ctx.ConfigManager.GetValuesOrSubmitForm(questions, cmd); err != nil {
				return err
			}

			createEvent := forms.Bool(answers["create-event"])
			askToCreateEvent = createEvent

			if createEvent {
				questions = []config.Value{
					{Question: forms.Question{
						Type:     forms.Input,
						Key:      "event-type",
						Prompt:   "EVENT TYPE: ",
						Optional: true,
					}},
					{Question: forms.Question{
						Type:     forms.Input,
						Key:      "payload",
						Prompt:   "JSON PAYLOAD: ",
						Optional: true,
					}},
				}

				if answers, err = Ctx.ConfigManager.GetValuesOrSubmitForm(questions, cmd); err != nil {
					return err
				}

				eventType := forms.String(answers["event-type"])
				payload := forms.String(answers["payload"])

				var payloadMap map[string]interface{}
				if payload != "" {
					if err = json.Unmarshal([]byte(payload), &payloadMap); err != nil {
						return err
					}
				}

				if useLegacyAPI {
					eventRequest := headset.EventRequest{
						SessionID: &sessionID,
						Type:      eventType,
						Payload:   payloadMap,
					}

					spinner = loader.NewLoader(cmd.Context(), "Creating event...", Ctx.Printer)
					_, err = Ctx.HeadsetClient.SendEvent(cmd.Context(), eventRequest)
					spinner.Stop()
					if err != nil {
						return err
					}

					Ctx.Printer.Println(":white_check_mark: Event created for session")

				} else {
					event := &platform.Event{
						SessionID: &sessionID,
						Type:      eventType,
						Payload:   payload,
					}

					spinner = loader.NewLoader(cmd.Context(), "Creating event...", Ctx.Printer)
					err = Ctx.PlatformClient.CreateEvent(cmd.Context(), event)
					spinner.Stop()
					if err != nil {
						return err
					}
				}

				Ctx.Printer.Println(":white_check_mark: Event created for session")
			}
		}

		questions = []config.Value{
			{Question: forms.Question{
				Type:     forms.Input,
				Key:      "score",
				Optional: true,
			}},
			{Question: forms.Question{
				Type:     forms.Input,
				Key:      "max-score",
				Optional: true,
			}},
			{Question: forms.Question{
				Type:   forms.Confirm,
				Key:    "session-passed",
				Prompt: "Passed?",
			}},
		}

		if answers, err = Ctx.ConfigManager.GetValuesOrSubmitForm(questions, cmd); err != nil {
			return err
		}

		score := forms.Int(answers["score"])
		maxScore := forms.Int(answers["max-score"])
		sessionPassed := forms.Bool(answers["session-passed"])

		lessonStatus := "passed"
		lessonStatusEmoji := "1st_place_medal"
		if !sessionPassed {
			lessonStatusEmoji = "x"
			lessonStatus = "failed"
		}

		var session *platform.Session
		if useLegacyAPI {
			sessionDuration := time.Since(startTime)
			eventRequest := headset.EventRequest{
				SessionID: &sessionID,
				Payload: map[string]interface{}{
					"sessionDuration": sessionDuration.Seconds(),
					"lessonStatus":    lessonStatus,
					"result": map[string]interface{}{
						"score": map[string]interface{}{
							"raw": score,
							"max": maxScore,
						},
					},
				},
			}

			spinner = loader.NewLoader(cmd.Context(), "Ending session...", Ctx.Printer)
			res, err := Ctx.HeadsetClient.EndSession(cmd.Context(), eventRequest)
			spinner.Stop()
			if err != nil {
				return err
			}
			session = res.Event.Session
			session.LessonStatus = lessonStatus
			session.Duration = sessionDuration.String()
			session.RawScore = float64(score)
			session.MaxScore = float64(maxScore)
			session.ScaledScore = float64(score) / float64(maxScore)
			session.Scenario = scenario
			session.Mode = mode
			session.Focus = focus
			session.Specialization = specialization
		} else {
			session = &platform.Session{
				ID:           sessionID,
				LessonStatus: lessonStatus,
				RawScore:     float64(forms.Int(answers["score"])),
				MaxScore:     float64(forms.Int(answers["max-score"])),
				Completed:    true,
			}

			spinner = loader.NewLoader(cmd.Context(), "Ending session...", Ctx.Printer)
			session, err = Ctx.PlatformClient.UpdateSession(cmd.Context(), *session)
			spinner.Stop()
			if err != nil {
				Ctx.Printer.Println(":exclamation: Unable to end session")
				return err
			}
		}

		percentScore := int(session.ScaledScore * 100)

		Ctx.Printer.Println("\n:white_check_mark:  Session completed")
		Ctx.Printer.Printf(":alarm_clock:  Duration: %s\n\n", session.Duration)

		if mode != "" {
			Ctx.Printer.Printf(":book: Mode: %s\n", mode)
		}
		if scenario != "" {
			Ctx.Printer.Printf(":magnifying_glass_tilted_right: Scenario: %s\n", scenario)
		}
		if focus != "" {
			Ctx.Printer.Printf(":magnifying_glass_tilted_right: Focus: %s\n", focus)
		}
		if specialization != "" {
			Ctx.Printer.Printf(":glowing_star: Specialization: %s\n", specialization)
		}

		Ctx.Printer.Printf("\n:input_numbers: Score: %.2f/%.2f\n", session.RawScore, session.MaxScore)
		Ctx.Printer.Printf(":hundred_points: Percent: %d%s\n", percentScore, "%")
		Ctx.Printer.Printf(":%s: Lesson Status: %s\n\n", lessonStatusEmoji, session.LessonStatus)
		return nil
	},
}

func init() {
	sessionsCmd.AddCommand(sessionsSimulateCmd)
	sessionsSimulateCmd.Flags().StringP("module", "m", "", "Module Abbreviation")
	sessionsSimulateCmd.Flags().String("mode", "", "Session mode: tutorial, practice, challenge")
	sessionsSimulateCmd.Flags().String("scenario", "", "Module Scenario")
	sessionsSimulateCmd.Flags().String("focus", "", "Area of focus")
	sessionsSimulateCmd.Flags().String("specialization", "", "Fine grained specialization within focus")
	sessionsSimulateCmd.Flags().BoolVar(&useLegacyAPI, "legacy", false, "Uses legacy Headset API")
}
