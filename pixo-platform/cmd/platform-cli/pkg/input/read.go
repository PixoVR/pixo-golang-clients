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

func ToInt(val string) int {
	intVal, err := strconv.Atoi(val)
	if err != nil {
		log.Error().Err(err).Msgf("unable to convert value %s to int", val)
		return 0
	}

	return intVal
}

func GetIntValue(cmd *cobra.Command, flagName, envVarName string) int {
	var val string
	if cmd.Flag(flagName).Value.String() != "" {
		val = cmd.Flag(flagName).Value.String()
	}

	val, ok := os.LookupEnv(envVarName)
	if ok {
		return ToInt(val)
	}

	val = viper.GetString(flagName)
	if val != "" {
		return ToInt(val)
	}

	if val == "" {
		inputVal := ReadFromUser(fmt.Sprintf("Enter %s: ", flagName))
		if inputVal != "" {
			val = inputVal
		}
	}

	if val == "" {
		val = cmd.Flag(flagName).DefValue
	}

	return ToInt(val)
}

func GetStringValue(cmd *cobra.Command, flagName, envVarName string) string {
	if cmd.Flag(flagName).Value.String() != "" {
		return cmd.Flag(flagName).Value.String()
	}

	val, ok := os.LookupEnv(envVarName)
	if ok {
		return val
	}

	val = viper.GetString(flagName)
	if val != "" {
		return val
	}

	inputVal := ReadFromUser(fmt.Sprintf("Enter %s: ", flagName))
	if inputVal != "" {
		return inputVal
	}

	return cmd.Flag(flagName).DefValue
}

func ReadFromUser(prompt string) string {
	fmt.Print("\n", prompt)

	reader := bufio.NewReader(os.Stdin)
	message, _ := reader.ReadString('\n')

	return strings.Trim(message, "\r\n")
}
