package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "tadpoles-backup",
		Short: "Backup photos of your child from www.tadpoles.com",
	}
	authToken string
	//authEmail    string
	//authPassword string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&authToken, "authToken", "t", "", "Auth Token for Tadpoles account.")
	viper.BindPFlag("authToken", rootCmd.PersistentFlags().Lookup("authToken"))

	//rootCmd.PersistentFlags().StringVarP(&authEmail, "email", "e", "", "Loginn email for Tadpoles account.")
	//viper.BindPFlag("authEmail", rootCmd.PersistentFlags().Lookup("authEmail"))
	//
	//rootCmd.PersistentFlags().StringVarP(&authPassword, "password", "p", "", "Login Password for Tadpoles account.")
	//viper.BindPFlag("authPassword", rootCmd.PersistentFlags().Lookup("authPassword"))

	rootCmd.AddCommand(statCmd)
	rootCmd.AddCommand(backupCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
