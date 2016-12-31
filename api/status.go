// Copyright (c) 2016 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package api

import (
	"net/http"

	l4g "github.com/alecthomas/log4go"

	"github.com/mattermost/platform/app"
	"github.com/mattermost/platform/model"
	"github.com/mattermost/platform/store"
	"github.com/mattermost/platform/utils"
)

func InitStatus() {
	l4g.Debug(utils.T("api.status.init.debug"))

	BaseRoutes.Users.Handle("/status", ApiUserRequired(getStatusesHttp)).Methods("GET")
	BaseRoutes.Users.Handle("/status/ids", ApiUserRequired(getStatusesByIdsHttp)).Methods("POST")
	BaseRoutes.WebSocket.Handle("get_statuses", ApiWebSocketHandler(getStatusesWebSocket))
	BaseRoutes.WebSocket.Handle("get_statuses_by_ids", ApiWebSocketHandler(getStatusesByIdsWebSocket))
}

func getStatusesHttp(c *Context, w http.ResponseWriter, r *http.Request) {
	statusMap, err := app.GetAllStatuses()
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(model.StringInterfaceToJson(statusMap)))
}

func getStatusesWebSocket(req *model.WebSocketRequest) (map[string]interface{}, *model.AppError) {
	statusMap, err := app.GetAllStatuses()
	if err != nil {
		return nil, err
	}

	return statusMap, nil
}

func getStatusesByIdsHttp(c *Context, w http.ResponseWriter, r *http.Request) {
	userIds := model.ArrayFromJson(r.Body)

	if len(userIds) == 0 {
		c.SetInvalidParam("getStatusesByIdsHttp", "user_ids")
		return
	}

	statusMap, err := app.GetStatusesByIds(userIds)
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(model.StringInterfaceToJson(statusMap)))
}

func getStatusesByIdsWebSocket(req *model.WebSocketRequest) (map[string]interface{}, *model.AppError) {
	var userIds []string
	if userIds = model.ArrayFromInterface(req.Data["user_ids"]); len(userIds) == 0 {
		l4g.Error(model.StringInterfaceToJson(req.Data))
		return nil, NewInvalidWebSocketParamError(req.Action, "user_ids")
	}

	statusMap, err := app.GetStatusesByIds(userIds)
	if err != nil {
		return nil, err
	}

	return statusMap, nil
}
