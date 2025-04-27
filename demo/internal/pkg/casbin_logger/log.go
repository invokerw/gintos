package casbin_logger

import (
	"fmt"
	"github/invokerw/gintos/log"
	"strings"

	log2 "github.com/casbin/casbin/v2/log"
)

// CasbinLogger is the implementation for a Logger using golang log.
type CasbinLogger struct {
	enabled bool
	logger  *log.Helper
}

var _ log2.Logger = (*CasbinLogger)(nil)

func NewCasbinLogger(helper *log.Helper) *CasbinLogger {
	return &CasbinLogger{
		enabled: true,
		logger:  helper,
	}
}

func (l *CasbinLogger) EnableLog(enable bool) {
	l.enabled = enable
}

func (l *CasbinLogger) IsEnabled() bool {
	return l.enabled
}

func (l *CasbinLogger) LogModel(model [][]string) {
	if !l.enabled {
		return
	}
	var str strings.Builder
	str.WriteString("Model: \n")
	for _, v := range model {
		str.WriteString(fmt.Sprintf("%v\n", v))
	}
	l.logger.Debug(str.String())
}

func (l *CasbinLogger) LogEnforce(matcher string, request []interface{}, result bool, explains [][]string) {
	if !l.enabled {
		return
	}

	var reqStr strings.Builder
	reqStr.WriteString("Request: \n")
	for i, rval := range request {
		if i != len(request)-1 {
			reqStr.WriteString(fmt.Sprintf("%v, ", rval))
		} else {
			reqStr.WriteString(fmt.Sprintf("%v", rval))
		}
	}
	reqStr.WriteString(fmt.Sprintf(" ---> %t\n", result))

	reqStr.WriteString("Hit Policy: \n")
	for i, pval := range explains {
		if i != len(explains)-1 {
			reqStr.WriteString(fmt.Sprintf("%v, ", pval))
		} else {
			reqStr.WriteString(fmt.Sprintf("%v \n", pval))
		}
	}

	l.logger.Debug(reqStr.String())
}

func (l *CasbinLogger) LogPolicy(policy map[string][][]string) {
	if !l.enabled {
		return
	}

	var str strings.Builder
	str.WriteString("Policy: \n")
	for k, v := range policy {
		str.WriteString(fmt.Sprintf("%s : %v\n", k, v))
	}

	l.logger.Debug(str.String())
}

func (l *CasbinLogger) LogRole(roles []string) {
	if !l.enabled {
		return
	}

	l.logger.Debug("Roles: \n", strings.Join(roles, "\n"))
}

func (l *CasbinLogger) LogError(err error, msg ...string) {
	if !l.enabled {
		return
	}
	l.logger.Error(msg, err)
}
