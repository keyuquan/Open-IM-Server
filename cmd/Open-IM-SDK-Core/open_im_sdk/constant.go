package open_im_sdk

import (
	"database/sql"
	"sync"
)

func initAddr() {
	ginAddress = SvrConf.IpApiAddr
	getUserInfoRouter = ginAddress + "/user/get_user_info"
	updateUserInfoRouter = ginAddress + "/user/update_user_info"
	addFriendRouter = ginAddress + "/friend/add_friend"
	getFriendApplicationListRouter = ginAddress + "/friend/get_friend_apply_list"
	getSelfApplicationListRouter = ginAddress + "/friend/get_self_apply_list"
	deleteFriendRouter = ginAddress + "/friend/delete_friend"
	getFriendInfoRouter = ginAddress + "/friend/get_friends_info"
	getFriendListRouter = ginAddress + "/friend/get_friend_list"
	sendMsgRouter = ginAddress + "/chat/send_msg"
	getBlackListRouter = ginAddress + "/friend/get_blacklist"
	addFriendResponse = ginAddress + "/friend/add_friend_response"
	addBlackListRouter = ginAddress + "/friend/add_blacklist"
	removeBlackListRouter = ginAddress + "/friend/remove_blacklist"
	//getFriendApplyListRouter = ginAddress + "/friend/get_friend_apply_list"
	pullUserMsgRouter = ginAddress + "/chat/pull_msg"
	pullUserMsgBySeqRouter = ginAddress + "/chat/pull_msg_by_seq"

	newestSeqRouter = ginAddress + "/chat/newest_seq"
	setFriendComment = ginAddress + "/friend/set_friend_comment"
	tencentCloudStorageCredentialRouter = ginAddress + "/third/tencent_cloud_storage_credential"

	//group
	createGroupRouter = ginAddress + "/group/create_group"
	setGroupInfoRouter = ginAddress + "/group/set_group_info"
	joinGroupRouter = ginAddress + "/group/join_group"
	quitGroupRouter = ginAddress + "/group/quit_group"
	getGroupsInfoRouter = ginAddress + "/group/get_groups_info"
	getGroupMemberListRouter = ginAddress + "/group/get_group_member_list"
	getGroupAllMemberListRouter = ginAddress + "/group/get_group_all_member_list"
	getGroupMembersInfoRouter = ginAddress + "/group/get_group_members_info"
	inviteUserToGroupRouter = ginAddress + "/group/invite_user_to_group"
	getJoinedGroupListRouter = ginAddress + "/group/get_joined_group_list"
	kickGroupMemberRouter = ginAddress + "/group/kick_group"
	transferGroupRouter = ginAddress + "/group/transfer_group"
	getGroupApplicationListRouter = ginAddress + "/group/get_group_applicationList"
	acceptGroupApplicationRouter = ginAddress + "/group/group_application_response"
	refuseGroupApplicationRouter = ginAddress + "/group/group_application_response"
}

var (
	ginAddress = ""

	getUserInfoRouter              = ""
	updateUserInfoRouter           = ""
	addFriendRouter                = ""
	getFriendInfoRouter            = ""
	getFriendApplicationListRouter = ""
	getSelfApplicationListRouter   = ""
	deleteFriendRouter             = ""
	getFriendListRouter            = ""
	sendMsgRouter                  = ""
	getBlackListRouter             = ""
	addFriendResponse              = ""
	addBlackListRouter             = ""
	removeBlackListRouter          = ""
	//	getFriendApplyListRouter            = ginAddress + "/friend/get_friend_apply_list"
	setFriendComment                    = " "
	pullUserMsgRouter                   = ""
	pullUserMsgBySeqRouter              = ""
	newestSeqRouter                     = ""
	tencentCloudStorageCredentialRouter = ""
	//group
	createGroupRouter             = ""
	setGroupInfoRouter            = ""
	joinGroupRouter               = ""
	quitGroupRouter               = ""
	getGroupsInfoRouter           = ""
	getGroupMemberListRouter      = ""
	getGroupAllMemberListRouter   = ""
	getGroupMembersInfoRouter     = ""
	inviteUserToGroupRouter       = ""
	getJoinedGroupListRouter      = ""
	kickGroupMemberRouter         = ""
	transferGroupRouter           = ""
	getGroupApplicationListRouter = ""
	acceptGroupApplicationRouter  = ""
	refuseGroupApplicationRouter  = ""
)

func (u *UserRelated) initListenerCh() {
	u.ch = make(chan cmd2Value, 1000)
	u.ConversationCh = u.ch

	u.wsNotification = make(map[string]chan GeneralWsResp, 1)
	u.seqMsg = make(map[int32]MsgData, 1000)
}

type UserRelated struct {
	ConversationCh chan cmd2Value //cmd：

	SvrConf        IMConfig
	token          string
	LoginUid       string
	wsNotification map[string]chan GeneralWsResp
	wsMutex        sync.RWMutex
	IMManager
	Friend
	ConversationListener
	groupListener

	//initDB     *sql.DB
	db         *sql.DB
	mRWMutex   *sync.RWMutex
	stateMutex sync.Mutex

	minSeqSvr        int64
	minSeqSvrRWMutex sync.RWMutex

	seqMsg      map[int32]MsgData
	seqMsgMutex sync.RWMutex
}

var UserSDKRwLock sync.RWMutex
var UserRouterMap map[string]*UserRelated
var SvrConf IMConfig
var SdkLogFlag int32

var userForSDK UserRelated

const (
	CmdFriend                     = "001"
	CmdBlackList                  = "002"
	CmdFriendApplication          = "003"
	CmdDeleteConversation         = "004"
	CmdNewMsgCome                 = "005"
	CmdGeyLoginUserInfo           = "006"
	CmdUpdateConversation         = "007"
	CmdForceSyncFriend            = "008"
	CmdFroceSyncBlackList         = "009"
	CmdForceSyncFriendApplication = "010"
	CmdForceSyncMsg               = "011"
	CmdForceSyncLoginUerInfo      = "012"
	CmdReLogin                    = "013"
	CmdUnInit                     = "014"
	CmdAcceptFriend               = "015"
	CmdRefuseFriend               = "016"
	CmdAddFriend                  = "017"
)

const (
	MessageHasNotRead = 0
	MessageHasRead    = 1
)
const (
	//ContentType
	Text           = 101
	Picture        = 102
	Sound          = 103
	Video          = 104
	File           = 105
	AtText         = 106
	Merger         = 107
	Card           = 108
	Location       = 109
	Custom         = 110
	Revoke         = 111
	HasReadReceipt = 112
	Typing         = 113
	Quote          = 114
	//////////////////////////////////////////
	SingleTipBegin             = 200
	AcceptFriendApplicationTip = 201
	AddFriendTip               = 202
	RefuseFriendApplicationTip = 203
	SetSelfInfoTip             = 204
	KickOnlineTip              = 303

	SingleTipEnd = 399
	/////////////////////////////////////////
	GroupTipBegin             = 500
	TransferGroupOwnerTip     = 501
	CreateGroupTip            = 502
	JoinGroupTip              = 504
	QuitGroupTip              = 505
	SetGroupInfoTip           = 506
	AcceptGroupApplicationTip = 507
	RefuseGroupApplicationTip = 508
	KickGroupMemberTip        = 509
	InviteUserToGroupTip      = 510

	GroupTipEnd = 599
	////////////////////////////////////////
	//MsgFrom
	UserMsgType = 100
	SysMsgType  = 200

	/////////////////////////////////////
	//SessionType
	SingleChatType = 1
	GroupChatType  = 2

	//MsgStatus
	MsgStatusSending     = 1
	MsgStatusSendSuccess = 2
	MsgStatusSendFailed  = 3
	MsgStatusHasDeleted  = 4
	MsgStatusRevoked     = 5
)

const (
	ckWsInitConnection  string = "ws-init-connection"
	ckWsLoginConnection string = "ws-login-connection"
	ckWsClose           string = "ws-close"
	ckWsKickOffLine     string = "ws-kick-off-line"
	ckTokenExpired      string = "token-expired"
	ckSelfInfoUpdate    string = "self-info-update"
)

const (
	ErrCodeInitLogin    = 1001
	ErrCodeFriend       = 2001
	ErrCodeConversation = 3001
	ErrCodeUserInfo     = 4001
	ErrCodeGroup        = 5001
)

const (
	LoginSuccess = 101
	Logining     = 102
	LoginFailed  = 103

	LogoutCmd = 201
)

const (
	DeFaultSuccessMsg = "ok"
)

const (
	ConAndUnreadChange        = 1
	AddConOrUpLatMsg          = 2
	UnreadCountSetZero        = 3
	ConChange                 = 4
	IncrUnread                = 5
	TotalUnreadMessageChanged = 6
	UpdateFaceUrlAndNickName  = 7

	HasRead = 1
	NotRead = 0
)

const (
	GroupActionCreateGroup            = 1
	GroupActionApplyJoinGroup         = 2
	GroupActionQuitGroup              = 3
	GroupActionSetGroupInfo           = 4
	GroupActionKickGroupMember        = 5
	GroupActionTransferGroupOwner     = 6
	GroupActionInviteUserToGroup      = 7
	GroupActionAcceptGroupApplication = 8
	GroupActionRefuseGroupApplication = 9
)
const ZoomScale = "200"
const MaxTotalMsgLen = 2048
const (
	FriendAcceptTip  = "You have successfully become friends, so start chatting"
	TransferGroupTip = "The owner of the group is transferred!"
	AcceptGroupTip   = "%s join the group"
)

const (
	WSGetNewestSeq     = 1001
	WSPullMsg          = 1002
	WSSendMsg          = 1003
	WSPullMsgBySeqList = 1004
	WSPushMsg          = 2001
	WSDataError        = 3001
)
