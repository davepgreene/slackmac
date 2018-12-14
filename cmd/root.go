package cmd

import (
	"crypto"
	"fmt"
	"github.com/davepgreene/slackmac/config"
	"github.com/davepgreene/slackmac/http"
	"github.com/davepgreene/slackmac/utils"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var verbose bool

// SlackMacCommand represents the base command when called without any subcommands
var SlackMacCommand = &cobra.Command{
	Use:   "slackmac",
	Short: "A tiny HTTP proxy that validates Slack payloads.",
	Long: `SlackMac is a tiny HTTP proxy that validates Slack payloads.
I wrote this because Spring Boot has an issue where it is impossible
to get the raw request body parameters in the correct order if a POST
request is sent with Content-Type application/x-www-form-urlencoded.
By offloading this work to a proxy, SlackMac can be dropped in front
of any service that needs to validate Slack payloads without the developer
ever having to worry about calculating HMACs. It's already done!`,
	Run: func(cmd *cobra.Command, args []string) {
		err := initializeConfig()
		initializeLog()
		if err != nil {
			log.Fatal(err)
		}

		SupportedAlgorithms := map[string]interface{}{
			"SHA256": crypto.SHA256,
		}

		// Prevent boot if we aren't using a supported algorithm
		confAlg := viper.GetString("slack.algorithm")
		if val, ok := SupportedAlgorithms[strings.ToUpper(confAlg)]; ok {
			// Set the actual crypto algorithm instead of a string representation
			viper.Set("slack.algorithm", val)

			// Instantiate our metrics client early
			_, err := utils.Metrics()
			if err != nil {
				log.Error(err)
			}
			log.Fatal(http.Handler())
		}

		supported := strings.Join(utils.MapKeys(SupportedAlgorithms), ", ")
		log.Fatalf("Slack currently supports the following encryption algorithms: %s. You specified %s.", supported, confAlg)
	},
}


// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the SlackMacCmd.
func Execute() {
	if err := SlackMacCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	SlackMacCommand.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
	SlackMacCommand.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose level logging")
	validConfigFilenames := []string{"json", "toml"}
	err := SlackMacCommand.PersistentFlags().SetAnnotation("config", cobra.BashCompFilenameExt, validConfigFilenames)
	if err != nil {
		log.Fatal(err)
	}
}

func initializeLog() {
	log.RegisterExitHandler(func() {
		log.Info("Shutting down")
	})

	// Set logging options based on config
	log.SetLevel(utils.GetLogLevel())

	// If using verbose mode, log at debug level
	if verbose {
		log.Info("Verbose mode specified. Setting log level to DEBUG")
		log.SetLevel(log.DebugLevel)
	}

	if viper.GetBool("log.json") {
		log.SetFormatter(&log.JSONFormatter{})
	}

	if cfgFile != "" {
		log.WithFields(log.Fields{
			"file": viper.ConfigFileUsed(),
		}).Info("Loaded config file")
	}

}

func initializeConfig(subCmdVs ...*cobra.Command) error {
	config.Defaults()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}
	viper.AutomaticEnv() // read in environment variables that match

	return nil
}
