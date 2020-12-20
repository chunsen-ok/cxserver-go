// Package session 实现基于服务端的会话管理。
//
// ## 会话的生命周期
//
// - 创建会话：用户登陆、验证通过
// - 维持会话：服务器存储<session_id,session>键值对，定时更新session_id
// - 结束会话：从服务器移除<session_id,session>键值对
//
package session

type SessionMgr struct {
	Sessions map[string]Session
}

// NewSession 创建会话
// 创建成功返回会话ID。
func (sm *SessionMgr) NewSession() (string, bool) {
	return "", false
}

// Refresh 更新会话的标识(session id)
// 遍历检查过期时间，更新成功返回新的会话ID。
// func (sm *SessionMgr) Refresh() (string, bool) {
// 	return "", false
// }

// DropSession 移除会话
func (sm *SessionMgr) DropSession() bool {
	return false
}

// GetSession 获取会话
// 成功获取，返回一个不等于nil的Session对象指针。
func (sm SessionMgr) GetSession(id string) *Session {
	return nil
}
