package input

import (
	"bufio"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
)

func GetIntValueOrAskUser(cmd *cobra.Command, flagName, envVarName string) int {
	return ToInt(GetStringValueOrAskUser(cmd, flagName, envVarName))
}

func GetStringValueOrAskUser(cmd *cobra.Command, flagName, envVarName string) string {
	if cmd.Flag(flagName).Value.String() != "" {
		return cmd.Flag(flagName).Value.String()
	}

	val := GetConfigValue(flagName, envVarName)
	if val != "" {
		return val
	}

	val = ReadFromUser(fmt.Sprintf("Enter %s: ", strings.ReplaceAll(flagName, "-", " ")))
	if val != "" {
		return val
	}

	return cmd.Flag(flagName).DefValue
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
		return val
	}

	if val = viper.GetString(flagName); val != "" {
		return val
	}

	return ""
}

func ReadFromUser(prompt string) string {
	fmt.Print("\n", prompt)

	reader := bufio.NewReader(os.Stdin)
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
