package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//// 捷径函数 ////

func flag(cmd *cobra.Command, name string, defaultVal interface{}, usage string) {
	f := cmd.PersistentFlags()
	switch y := defaultVal.(type) {
	case bool:
		f.Bool(name, y, usage)
	case int:
		f.Int(name, y, usage)
	case string:
		f.String(name, y, usage)
	}
	viper.SetDefault(name, defaultVal)
}

func flagP(cmd *cobra.Command, name, shorthand string, defaultVal interface{}, usage string) {
	f := cmd.PersistentFlags()
	switch y := defaultVal.(type) {
	case bool:
		f.BoolP(name, shorthand, y, usage)
	case int:
		f.IntP(name, shorthand, y, usage)
	case string:
		f.StringP(name, shorthand, y, usage)
	}
	viper.SetDefault(name, defaultVal)
}

func flagV(cmd *cobra.Command, name string, defaultVal interface{}, usage string) {
	flag(cmd, name, defaultVal, usage)
	viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
}

func flagPV(cmd *cobra.Command, name, shorthand string, defaultVal interface{}, usage string) {
	flagP(cmd, name, shorthand, defaultVal, usage)
	viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
}
