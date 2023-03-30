package handle

import (
	"net"
	"sync"
)

type ConnectUserList struct {
	sync.Map
}

type ConnectUser struct {
	UserID  string   //用户ID
	Connect net.Conn //用户链接
	Level   int      //用户等级
}

// 增改
func (ul *ConnectUserList) AddUser(user ConnectUser) (ConnectUser, bool) {
	// 同一时间只有一个登录
	actual, loaded := ul.LoadOrStore(user.UserID, user)
	oldConnectUser := actual.(ConnectUser)
	return oldConnectUser, loaded

}

// 删除
func (ul *ConnectUserList) RemoveUser(userID string) {
	ul.Delete(userID)
}

// 查
func (ul *ConnectUserList) GetUser(userID string) (bool, ConnectUser) {
	user, ok := ul.Load(userID)
	if !ok {
		return false, ConnectUser{}
	}
	return true, user.(ConnectUser)
}

// 找最大
func (ul *ConnectUserList) GetUserMaxLevel() (ConnectUser, bool) {
	var MaxLeveL int
	var MaxLeveLUser ConnectUser
	ul.Range(func(key, value any) bool {
		user := value.(ConnectUser)
		if user.Level > MaxLeveL || MaxLeveLUser.Connect == nil {
			MaxLeveL = user.Level
			MaxLeveLUser = user
		}
		return true
	})
	if MaxLeveLUser.Connect != nil {
		return MaxLeveLUser, true
	}
	return ConnectUser{}, false
}
