package cmd

import (
	"log"
	"os"

	"github.com/9d4/semaphore/server"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var rootCmd = cobra.Command{
	Use:   "semaphore",
	Short: "Start semaphore server.",
	Long:  "Semaphore is blablabla..........",
	Run: boot(func(cmd *cobra.Command, args []string, passData *bootData) {
		log.Fatal(server.Start())
	}),
}

var serverFlags = flag.NewFlagSet(rootCmd.Name(), flag.ContinueOnError)

func init() {
	initFlags()
	cobra.OnInitialize(func() { initConfig(); initLogger() })

	rootCmd.PersistentFlags().AddFlagSet(serverFlags)
	viper.BindPFlags(serverFlags)
}

func initConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.ReadInConfig()
}

func initFlags() {
	serverFlags.StringP("addr", "a", "0.0.0.0:3500", "Address to listen on")
	serverFlags.String("dbhost", "127.0.0.1", "Database host")
	serverFlags.String("dbport", "5432", "Database port")
	serverFlags.String("dbname", "semaphore", "Database name")
	serverFlags.String("dbuser", "semaphore", "Database user")
	serverFlags.String("dbpasswd", "smphr", "Database password")
}

func initLogger() {
	logWriter, err := os.OpenFile("semaphore.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Println("Unable to create log file:", err)
	}

	if err == nil {
		jww.SetLogOutput(logWriter)
	}

	jww.SetLogThreshold(jww.LevelTrace)
	jww.SetStdoutThreshold(jww.LevelInfo)
}
