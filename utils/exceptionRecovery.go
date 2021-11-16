package utils

import (
	"github.com/gin-gonic/gin"
)

type ExceptionProcess struct {
	Ctx *gin.Context
}

func (ep *ExceptionProcess) ExceptionRecovery(fnName string, args... string)  {
	if r := recover(); r !=nil {
		//logger.Error(r.(string))

		e := RunFromCallMethodName(ep, fnName, args...)

		if e != nil {
			//logger.Error(e.Error())
		}
	}
}

/*
	Exception 처리 Methods
 */

// LiveApi Profile Exception
func (ep *ExceptionProcess) GetProfileException()  {
	//logger.Error("GetProfileException")
}
