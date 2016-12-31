// Copyright (c) 2016 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	l4g "github.com/alecthomas/log4go"
	"github.com/mattermost/platform/model"
	"github.com/mattermost/platform/utils"
)

func MakeDirectChannelVisible(channelId string) {
	var members []model.ChannelMember
	if result := <-Srv.Store.Channel().GetMembers(channelId); result.Err != nil {
		l4g.Error(utils.T("api.post.make_direct_channel_visible.get_members.error"), channelId, result.Err.Message)
		return
	} else {
		members = result.Data.([]model.ChannelMember)
	}

	if len(members) != 2 {
		l4g.Error(utils.T("api.post.make_direct_channel_visible.get_2_members.error"), channelId)
		return
	}

	// make sure the channel is visible to both members
	for i, member := range members {
		otherUserId := members[1-i].UserId

		if result := <-Srv.Store.Preference().Get(member.UserId, model.PREFERENCE_CATEGORY_DIRECT_CHANNEL_SHOW, otherUserId); result.Err != nil {
			// create a new preference since one doesn't exist yet
			preference := &model.Preference{
				UserId:   member.UserId,
				Category: model.PREFERENCE_CATEGORY_DIRECT_CHANNEL_SHOW,
				Name:     otherUserId,
				Value:    "true",
			}

			if saveResult := <-Srv.Store.Preference().Save(&model.Preferences{*preference}); saveResult.Err != nil {
				l4g.Error(utils.T("api.post.make_direct_channel_visible.save_pref.error"), member.UserId, otherUserId, saveResult.Err.Message)
			} else {
				message := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_PREFERENCE_CHANGED, "", "", member.UserId, nil)
				message.Add("preference", preference.ToJson())

				go Publish(message)
			}
		} else {
			preference := result.Data.(model.Preference)

			if preference.Value != "true" {
				// update the existing preference to make the channel visible
				preference.Value = "true"

				if updateResult := <-Srv.Store.Preference().Save(&model.Preferences{preference}); updateResult.Err != nil {
					l4g.Error(utils.T("api.post.make_direct_channel_visible.update_pref.error"), member.UserId, otherUserId, updateResult.Err.Message)
				} else {
					message := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_PREFERENCE_CHANGED, "", "", member.UserId, nil)
					message.Add("preference", preference.ToJson())

					go Publish(message)
				}
			}
		}
	}
}
