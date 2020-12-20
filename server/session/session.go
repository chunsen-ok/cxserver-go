package session

// Session 表示一个会话。
// 一个 Session 对象中应该存放哪些数据呢？
// 服务器通过客户端上传的 session id 获得对应的 Session 对象后，
// 即可表示该用户已经登录了。但还不知道究竟是哪个用户，所以 Session 对象中需要
// 存放有关用户本身的信息。如用户ID，用户名等。起码需要一个能去数据库中获取到该用户
// 详细数据的关键信息，例如用户ID。
//
// 在 Session 对象中只存放一个用户ID的好处是不需要占用太多内存空间。缺点是在需要用户的
// 其他信息的时候就需要查询数据库。
//
// 基于以上原因，在 Session 对象中可以存放一个用户常用的信息，一方面减少内存占用，另一方面
// 尽可能减少数据库访问。
type Session struct {
	ExpireTime int
	User       UserData
}

type UserData struct {
	id int
}
