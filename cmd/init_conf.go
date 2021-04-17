package cmd

import (
	"git.goasum.de/jasper/overtime-cli/internal/conf"
	"git.goasum.de/jasper/overtime-cli/internal/utils"
)

func InitConf() error {
	url := utils.ReadTextFromStdin("overtime api url:\n")
	token := utils.ReadTextFromStdin("access token:\n")
	defaultActivityDesc := utils.ReadTextFromStdin("default activity description (empty if not needed):\n")
	c := conf.Config{
		Host:                url,
		Token:               token,
		DefaultActivityDesc: defaultActivityDesc,
	}
	return conf.WriteConfig(c)
}
