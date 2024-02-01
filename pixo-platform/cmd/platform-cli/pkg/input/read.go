package input

import (
	"bufio"
	"fmt"
	"github.com/kyokomi/emoji"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gitlab.com/david_mbuvi/go_asterisks"
	"io"
	"os"
	"strconv"
	"strings"
)

func GetIntValueOrAskUser(cmd *cobra.Command, flagName, envVarName string) int {
	return ToInt(GetStringValueOrAskUser(cmd, flagName, envVarName))
}

func GetSensitiveStringValueOrAskUser(cmd *cobra.Command, flagName, envVarName string, defaultVal ...string) string {
	if cmd.Flag(flagName).Value.String() != "" {
		return cmd.Flag(flagName).Value.String()
	}

	//val := GetConfigValue(flagName, envVarName)
	//if val != "" {
	//	return val
	//}

	prompt := fmt.Sprintf("Enter %s", strings.ReplaceAll(flagName, "-", " "))
	if len(defaultVal) > 0 {
		prompt = fmt.Sprintf("%s (press enter for default - %s): ", prompt, defaultVal[0])
	} else {
		prompt = fmt.Sprintf("%s: ", prompt)
	}

	val := ReadSensitiveFromUser(cmd.OutOrStdout(), prompt)
	if val != "" {
		return val
	}

	if len(defaultVal) > 0 {
		return defaultVal[0]
	}

	return cmd.Flag(flagName).DefValue
}

func GetStringValueOrAskUser(cmd *cobra.Command, flagName, envVarName string, defaultVal ...string) string {
	if cmd.Flag(flagName).Value.String() != "" {
		return cmd.Flag(flagName).Value.String()
	}

	//val := GetConfigValue(flagName, envVarName)
	//if val != "" {
	//	return val
	//}

	prompt := fmt.Sprintf("Enter %s", strings.ReplaceAll(flagName, "-", " "))
	if len(defaultVal) > 0 {
		prompt = fmt.Sprintf("%s (press enter for default - %s): ", prompt, defaultVal[0])
	} else {
		prompt = fmt.Sprintf("%s: ", prompt)
	}

	val := ReadFromUser(cmd.InOrStdin(), cmd.OutOrStdout(), prompt)
	if val != "" {
		return val
	}

	if len(defaultVal) > 0 {
		return defaultVal[0]
	}

	return cmd.Flag(flagName).DefValue
}

func GetIntValue(cmd *cobra.Command, flagName, envVarName string) int {
	flag := cmd.Flag(flagName)
	if flag != nil && flag.Value != nil && flag.Value.String() != "" {
		return ToInt(cmd.Flag(flagName).Value.String())
	}

	//if val := GetConfigValue(flagName, envVarName); val != "" {
	//	return ToInt(val)
	//}

	if flag != nil && flag.DefValue != "" {
		return ToInt(cmd.Flag(flagName).DefValue)
	}

	return 0
}

func GetStringValue(cmd *cobra.Command, flagName, envVarName string) string {
	if cmd.Flag(flagName).Value.String() != "" {
		return cmd.Flag(flagName).Value.String()
	}

	//if val := GetConfigValue(flagName, envVarName); val != "" {
	//	return val
	//}

	return cmd.Flag(flagName).DefValue
}

func ReadSensitiveFromUser(writer io.Writer, prompt string) string {
	if writer == nil {
		return ""
	}

	prompt = emoji.Sprintf(":lock: %s", prompt)
	val, err := go_asterisks.GetUsersPassword(prompt, true, os.Stdin, writer)
	if err != nil {
		log.Error().Err(err).Msg("unable to read password")
		return ""
	}

	return strings.Trim(string(val), "\r\n")
}

func ReadFromUser(reader io.Reader, writer io.Writer, prompt string) string {
	if writer == nil || reader == nil {
		return ""
	}

	if _, err := writer.Write([]byte(emoji.Sprintf(":fountain_pen: Enter %s: ", prompt))); err != nil {
		log.Error().Err(err).Msg("unable to write to writer")
		return ""
	}

	bytesReader := bufio.NewReader(reader)
	message, _ := bytesReader.ReadString('\n')

	return strings.Trim(message, "\r\n")
}

func ToInt(val string) int {
	intVal, err := strconv.Atoi(val)
	if err != nil {
		log.Error().Err(err).Msgf("unable to convert value %s to int", val)
		return 0
	}

	return intVal
}
