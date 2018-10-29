package config

import "de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/helpers"

type Configuration struct {
	LoginDataFolderPath string
	AccessTokenCookie helpers.Cookie
}
