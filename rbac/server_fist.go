package rbac

import (
	"github.com/emicklei/go-restful"
	"github.com/fanux/fist/tools"
)

//FistRegister is fist auth controller
func FistRegister(auth *restful.WebService) {
	auth.Path("/").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON) // you can specify this per route as well
	//login http server
	auth.Route(auth.POST("/login").To(handleLogin))
	//logout http server
	auth.Route(auth.POST("/logout").Filter(cookieFilter).To(handleLogout))
	//user manager
	//GET_USER ALL
	auth.Route(auth.GET("/user").Filter(cookieFilter).To(handleListUserInfo))
	//GET_USER SINGLE
	auth.Route(auth.GET("/user/{user_name}").Filter(cookieFilter).To(handleGetUserInfo))
	//ADD_USER
	auth.Route(auth.POST("/user").Filter(cookieFilter).To(handleAddUserInfo))
	//UPDATE_USER
	auth.Route(auth.PUT("/user").Filter(cookieFilter).To(handleUpdateUserInfo))
	//DELETE_USER
	auth.Route(auth.DELETE("/user/{user_name}").Filter(cookieFilter).To(handleDelUserInfo))
}

// cookie Filter
func cookieFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	if filterCookieValidate(req) {
		chain.ProcessFilter(req, resp)
	} else {
		tools.ResponseAuthError(resp)
	}
}
func handleLogout(request *restful.Request, response *restful.Response) {
	logoutCookieSetter(response)
	tools.ResponseSuccess(response, "")
}
func handleLogin(request *restful.Request, response *restful.Response) {
	t := &UserInfo{}
	err := request.ReadEntity(t)
	if err != nil {
		tools.ResponseSystemError(response, err)
		return
	}
	uerInfo := DoAuthentication(t.Username, t.Password)
	if uerInfo == nil {
		tools.ResponseError(response, tools.ErrUserAuth)
		return
	}
	loginCookieSetter(response, uerInfo)
	tools.ResponseSuccess(response, uerInfo)
}

func handleGetUserInfo(request *restful.Request, response *restful.Response) {
	userName := request.PathParameter("user_name")
	// is exists
	if !validateUserNameExist(userName) {
		tools.ResponseError(response, tools.ErrUserNotExists)
		return
	}
	userInfo := GetUserInfo(userName, false)
	if userInfo == nil {
		tools.ResponseError(response, tools.ErrUserGet)
		return
	}
	tools.ResponseSuccess(response, userInfo)
}

func handleListUserInfo(request *restful.Request, response *restful.Response) {
	arr := ListAllUserInfo(false)
	tools.ResponseSuccess(response, arr)
}

func handleAddUserInfo(request *restful.Request, response *restful.Response) {
	t := &UserInfo{}
	err := request.ReadEntity(t)
	if err != nil {
		tools.ResponseSystemError(response, err)
		return
	}
	//1 user name is error ?
	if validateUserName(t.Username) {
		tools.ResponseSystemError(response, tools.ErrUserName)
		return
	}
	//3 user is  not exists ?
	if validateUserNameExist(t.Username) {
		tools.ResponseSystemError(response, tools.ErrUserExists)
		return
	}
	err = AddUserInfo(t)
	if err != nil {
		tools.ResponseError(response, tools.ErrUserAdd)
		return
	}
	tools.ResponseSuccess(response, nil)
}

func handleUpdateUserInfo(request *restful.Request, response *restful.Response) {
	t := &UserInfo{}
	err := request.ReadEntity(t)
	if err != nil {
		tools.ResponseSystemError(response, err)
		return
	}
	//1 user name is error ?
	if validateUserName(t.Username) {
		tools.ResponseSystemError(response, tools.ErrUserName)
		return
	}
	//3 user is   exists ?
	if !validateUserNameExist(t.Username) {
		tools.ResponseSystemError(response, tools.ErrUserNotExists)
		return
	}
	err = UpdateUserInfo(t)
	if err != nil {
		tools.ResponseError(response, tools.ErrUserUpdate)
		return
	}
	tools.ResponseSuccess(response, nil)
}

func handleDelUserInfo(request *restful.Request, response *restful.Response) {
	userName := request.PathParameter("user_name")
	// is exists
	if !validateUserNameExist(userName) {
		tools.ResponseError(response, tools.ErrUserNotExists)
		return
	}
	err := DelUserInfo(userName)
	if err != nil {
		tools.ResponseError(response, tools.ErrUserDel)
		return
	}
	tools.ResponseSuccess(response, nil)
}
