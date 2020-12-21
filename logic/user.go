package logic

// 存放业务逻辑的代码

func SignUp(p *models.ParamSignUp) (err error) {
	// 1.判断用户存不存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}
	// 2.生成用户ID

	// 3.保存进数据库
}