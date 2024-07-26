/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var libusbDbgLvl int
var devAddr DevAddrFlags

type DevAddrFlags struct {
	vid string
	pid string
	bus string
	dev string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "usbexp",
	Short: "Display info about the USB devices attached to the system",
	Long: `USB Explorer discovers the attached USB devices and allows
viewing the device info as well as resetting or performing
power operations on the device.
	
The USB device address consists of VID [PID [bus [address]]]
Example: -a 0x04B4 0x00F9  // Minrray cameras on all busses`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("rootCmd.Run vid %s pid %s\n", devAddr.vid, devAddr.pid)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("rootCmd.Execute failed! %v", err)
		return err
	}
	return nil
}

func init() {
	cobra.OnInitialize(initConfig)

	// define global flags and configuration settings.

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "cfg_file", "c", "", "config file (default is $HOME/.usbexp.yaml)")
	viper.BindPFlag("cfg_file", rootCmd.PersistentFlags().Lookup("cfg_file"))

	rootCmd.PersistentFlags().IntVarP(&libusbDbgLvl, "dbg_lvl", "l", 0, "-l=level set libusb debug output level)")
	viper.BindPFlag("dbg_lvl", rootCmd.PersistentFlags().Lookup("dbg_lvl"))

	rootCmd.PersistentFlags().StringVarP(&devAddr.vid, "vid", "v", "", "-v vendor ID in hex]")
	viper.BindPFlag("vid", rootCmd.PersistentFlags().Lookup("vid"))

	rootCmd.PersistentFlags().StringVarP(&devAddr.pid, "pid", "p", "", "-p product ID in hex]")
	viper.BindPFlag("pid", rootCmd.PersistentFlags().Lookup("pid"))

	rootCmd.PersistentFlags().StringVarP(&devAddr.bus, "bus", "b", "", "-b bus device is on")
	viper.BindPFlag("bus", rootCmd.PersistentFlags().Lookup("bus"))

	rootCmd.PersistentFlags().StringVarP(&devAddr.dev, "dev", "d", "", "-d device address")
	viper.BindPFlag("dev", rootCmd.PersistentFlags().Lookup("dev"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".usbexp" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".usbexp")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
