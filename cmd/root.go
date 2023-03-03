package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/9d4/semaphore/server"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var rootCmd = cobra.Command{
	Use:   "semaphore",
	Short: "Start semaphore server.",
	Long:  "Semaphore is blablabla..........",
	Run: func(cmd *cobra.Command, args []string) {
		srvErr, oauthSrvErr := server.Start(server.ParseViper(v))
		log.Fatal(<-oauthSrvErr)
		for {
			select {
			case err := <-oauthSrvErr:
				jww.FATAL.Fatal(err)
			case err := <-srvErr:
				jww.FATAL.Fatal(err)
			default:

			}
		}
	},
}

var (
	v           = viper.NewWithOptions(viper.EnvKeyReplacer(strings.NewReplacer("-", "_")))
	globalFlags = flag.NewFlagSet(rootCmd.Name(), flag.ContinueOnError)
	serverFlags = flag.NewFlagSet(rootCmd.Name(), flag.ContinueOnError)
)

func init() {
	loadFlags()
	cobra.OnInitialize(func() { loadEnv(); loadConfig(); loadLogger() })

	rootCmd.PersistentFlags().AddFlagSet(globalFlags)
	rootCmd.Flags().AddFlagSet(serverFlags)

	err := v.BindPFlags(globalFlags)
	if err != nil {
		return
	}

	err = v.BindPFlags(serverFlags)
	if err != nil {
		return
	}
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		jww.FATAL.Fatal("Error loading .env file")
	}
}

func loadConfig() {
	v.AutomaticEnv()
	err := v.ReadInConfig()
	if err != nil {
		return
	}
}

func loadFlags() {
	serverFlags.StringP("address", "a", "0.0.0.0:3500", "Address to listen on")
	serverFlags.String("address-oauth2", "0.0.0.0:3501", "Address for OAuth2 server to listen on")

	globalFlags.String("db-host", "127.0.0.1", "Database host")
	globalFlags.String("db-port", "5432", "Database port")
	globalFlags.String("db-name", "semaphore", "Database name")
	globalFlags.String("db-username", "semaphore", "Database user")
	globalFlags.String("db-password", "smphr", "Database password")
	globalFlags.BoolP("log-request", "l", false, "Print incoming request log")
}

func loadLogger() {
	logWriter, err := os.OpenFile("semaphore.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Println("Unable to create log file:", err)
	}

	if err == nil {
		jww.SetLogOutput(logWriter)
	}

	jww.SetFlags(log.Lshortfile)
	jww.SetLogThreshold(jww.LevelTrace)
	jww.SetStdoutThreshold(jww.LevelInfo)
}
