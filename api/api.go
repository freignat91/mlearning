package mlapi

import (
	"log"
	"strings"
)

const (
	//LOGERROR ERROR value
	LOGERROR = 0
	//LOGWARN WARN value
	LOGWARN = 1
	//LOGINFO INFO value
	LOGINFO = 2
	//LOGDEBUG debug value
	LOGDEBUG = 3
)

//MlAPI ml api
type MlAPI struct {
	serverAddress string
	logLevel      int
}

// New create an mlearning api instance
func New(servers string) *MlAPI {
	api := &MlAPI{
		serverAddress: "localhost",
		logLevel:      LOGWARN,
	}
	return api
}

func (api *MlAPI) getClient() (*mlClient, error) {
	client := mlClient{}
	err := client.init(api)
	if err != nil {
		log.Printf("Error connecting server: %v\n", err)
		return nil, err
	}
	return &client, nil
}

//SetLogLevel set the log level
func (api *MlAPI) SetLogLevel(level string) {
	if strings.ToLower(level) == "error" {
		api.logLevel = LOGERROR
	} else if strings.ToLower(level) == "warn" {
		api.logLevel = LOGWARN
	} else if strings.ToLower(level) == "info" {
		api.logLevel = LOGINFO
	} else if strings.ToLower(level) == "debug" {
		api.logLevel = LOGDEBUG
	}
}

//LogLevelString return a string log level
func (api *MlAPI) LogLevelString() string {
	switch api.logLevel {
	case LOGERROR:
		return "error"
	case LOGWARN:
		return "warn"
	case LOGINFO:
		return "info"
	case LOGDEBUG:
		return "debug"
	default:
		return "?"
	}
}

func (api *MlAPI) error(format string, args ...interface{}) {
	if api.logLevel >= LOGERROR {
		log.Printf(format, args...)
	}
}

func (api *MlAPI) warn(format string, args ...interface{}) {
	if api.logLevel >= LOGWARN {
		log.Printf(format, args...)
	}
}

func (api *MlAPI) info(format string, args ...interface{}) {
	if api.logLevel >= LOGINFO {
		log.Printf(format, args...)
	}
}

func (api *MlAPI) debug(format string, args ...interface{}) {
	if api.logLevel >= LOGDEBUG {
		log.Printf(format, args...)
	}
}

func (api *MlAPI) isDebug() bool {
	if api.logLevel >= LOGDEBUG {
		return true
	}
	return false
}

func (api *MlAPI) printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
