package conf

import (
	"git.goasum.de/jasper/overtime-cli/internal/utils"
)

func InitConf() error {
	url := utils.ReadTextFromStdin("overtime api url:\n")
	token := utils.ReadTextFromStdin("access token:\n")
	defaultActivityDesc := utils.ReadTextFromStdin("default activity description (empty if not needed):\n")
	c := Config{
		Host:                url,
		Token:               token,
		DefaultActivityDesc: defaultActivityDesc,
	}
	return WriteConfig(c)
}
