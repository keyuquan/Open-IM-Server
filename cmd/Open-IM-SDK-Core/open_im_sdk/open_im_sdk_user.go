package open_im_sdk

import (
	"encoding/json"
)

func (u *UserRelated) GetUsersInfo(uIDList string, cb Base) {
	go func() {
		var uidList []string
		err := json.Unmarshal([]byte(uIDList), &uidList)
		if err != nil {
			cb.OnError(ErrCodeUserInfo, err.Error())
			return
		}
		resp, err := post2Api(getUserInfoRouter, paramsGetUserInfo{UidList: uidList, OperationID: operationIDGenerator()}, u.token)
		if err != nil {
			cb.OnError(ErrCodeUserInfo, err.Error())
			return
		}
		var vgetUserInfoResp getUserInfoResp
		_ = json.Unmarshal(resp, &vgetUserInfoResp)
		if vgetUserInfoResp.ErrCode != 0 {
			cb.OnError(vgetUserInfoResp.ErrCode, vgetUserInfoResp.ErrMsg)
			return
		}
		jsonResp2Ui, _ := json.Marshal(vgetUserInfoResp.Data)
		cb.OnSuccess(string(jsonResp2Ui))
	}()
}

func (u *UserRelated) SetSelfInfo(info string, cb Base) {
	go func() {
		var uiUpdateUserInfo ui2UpdateUserInfo
		err := json.Unmarshal([]byte(info), &uiUpdateUserInfo)
		if err != nil {
			cb.OnError(ErrCodeUserInfo, err.Error())
			return
		}
		resp, err := post2Api(updateUserInfoRouter, paramsUpdateUserInfo{
			Name:        uiUpdateUserInfo.Name,
			Icon:        uiUpdateUserInfo.Icon,
			Gender:      uiUpdateUserInfo.Gender,
			Mobile:      uiUpdateUserInfo.Mobile,
			Birth:       uiUpdateUserInfo.Birth,
			Email:       uiUpdateUserInfo.Email,
			Ex:          uiUpdateUserInfo.Ex,
			OperationID: operationIDGenerator(),
		}, u.token)
		if err != nil {
			cb.OnError(ErrCodeUserInfo, err.Error())
			return
		}
		var cmResp commonResp
		_ = json.Unmarshal(resp, &cmResp)
		if cmResp.ErrCode != 0 {
			cb.OnError(cmResp.ErrCode, cmResp.ErrMsg)
			return
		}

		user, err := u.getServerUserInfo()
		if err != nil {
			cb.OnError(ErrCodeUserInfo, err.Error())
			return
		}

		err = u.replaceIntoUser(user)
		if err != nil {
			cb.OnError(ErrCodeUserInfo, err.Error())
			return
		}

		cb.OnSuccess("")
		u.cb.OnSelfInfoUpdated(structToJsonString(user))
	}()
}
