package input

import (
	"bufio"
	"fmt"
	"github.com/kyokomi/emoji"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/david_mbuvi/go_asterisks"
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

	val := GetConfigValue(flagName, envVarName)
	if val != "" {
		return val
	}

	prompt := fmt.Sprintf("Enter %s", strings.ReplaceAll(flagName, "-", " "))
	if len(defaultVal) > 0 {
		prompt = fmt.Sprintf("%s (press enter for default - %s): ", prompt, defaultVal[0])
	} else {
		prompt = fmt.Sprintf("%s: ", prompt)
	}

	val = ReadSensitiveFromUser(cmd, prompt)
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

	val := GetConfigValue(flagName, envVarName)
	if val != "" {
		return val
	}

	prompt := fmt.Sprintf("Enter %s", strings.ReplaceAll(flagName, "-", " "))
	if len(defaultVal) > 0 {
		prompt = fmt.Sprintf("%s (press enter for default - %s): ", prompt, defaultVal[0])
	} else {
		prompt = fmt.Sprintf("%s: ", prompt)
	}

	val = ReadFromUser(cmd, prompt)
	if val != "" {
		return val
	}

	if len(defaultVal) > 0 {
		return defaultVal[0]
	}

	return cmd.Flag(flagName).DefValue
}

func GetIntValue(cmd *cobra.Command, flagName, envVarName string) int {
	var val string
	flag := cmd.Flag(flagName)
	if flag != nil && flag.Value != nil && flag.Value.String() != "" {
		return ToInt(cmd.Flag(flagName).Value.String())
	}

	val = GetConfigValue(flagName, envVarName)
	if val != "" {
		return ToInt(val)
	}

	if flag != nil && flag.DefValue != "" {
		return ToInt(cmd.Flag(flagName).DefValue)
	}

	return 0
}

func GetStringValue(cmd *cobra.Command, flagName, envVarName string) string {
	if cmd.Flag(flagName).Value.String() != "" {
		return cmd.Flag(flagName).Value.String()
	}

	val := GetConfigValue(flagName, envVarName)
	if val != "" {
		return val
	}

	return cmd.Flag(flagName).DefValue
}

func GetConfigValue(flagName, envVarName string) string {
	val, ok := os.LookupEnv(envVarName)
	if ok {
		return strings.TrimSpace(val)
	}

	if val = viper.GetString(flagName); val != "" {
		return strings.TrimSpace(val)
	}

	return ""
}

func ReadSensitiveFromUser(cmd *cobra.Command, prompt string) string {
	val, err := go_asterisks.GetUsersPassword(prompt, true, os.Stdin, cmd.OutOrStdout())
	if err != nil {
		log.Error().Err(err).Msg("unable to read password")
		return ""
	}

	return strings.Trim(string(val), "\r\n")
}

func ReadFromUser(cmd *cobra.Command, prompt string) string {
	cmd.Print(emoji.Sprintf(":fountain_pen:%s", prompt))

	reader := bufio.NewReader(cmd.InOrStdin())
	message, _ := reader.ReadString('\n')

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
