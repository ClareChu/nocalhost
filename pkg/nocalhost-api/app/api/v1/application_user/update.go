package application_user

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"nocalhost/internal/nocalhost-api/service"
	"nocalhost/pkg/nocalhost-api/app/api"
	"nocalhost/pkg/nocalhost-api/pkg/errno"
	"nocalhost/pkg/nocalhost-api/pkg/log"
)

func BatchInsert(c *gin.Context) {
	// userId, _ := c.Get("userId")
	applicationId := cast.ToUint64(c.Param("id"))

	var req ApplicationUsersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf("BatchInsert bind params err: %v", err)
		api.SendResponse(c, errno.ErrBind, nil)
		return
	}

	var users []uint64
	for _, user := range req.Users {
		users = append(users, uint64(user))
	}

	err := service.Svc.ApplicationUser().BatchInsert(c, applicationId, users)

	if err != nil {
		api.SendResponse(c, errno.ErrInsertApplicationUser, nil)
		return
	}
	api.SendResponse(c, nil, nil)
}

func BatchDelete(c *gin.Context) {
	// userId, _ := c.Get("userId")
	applicationId := cast.ToUint64(c.Param("id"))

	var req ApplicationUsersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf("BatchDelete bind params err: %v", err)
		api.SendResponse(c, errno.ErrBind, nil)
		return
	}

	var users []uint64
	for _, user := range req.Users {
		users = append(users, uint64(user))
	}

	err := service.Svc.ApplicationUser().BatchDelete(c, applicationId, users)

	if err != nil {
		api.SendResponse(c, errno.ErrDeleteApplicationUser, nil)
		return
	}
	api.SendResponse(c, nil, nil)
}
