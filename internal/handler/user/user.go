package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	res "github.com/lin-snow/ech0/internal/handler/response"
	authModel "github.com/lin-snow/ech0/internal/model/auth"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	model "github.com/lin-snow/ech0/internal/model/user"
	service "github.com/lin-snow/ech0/internal/service/user"
)

type UserHandler struct {
	userService service.UserServiceInterface
}

func NewUserHandler(userService service.UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Login 用户登录
// @Summary 用户登录接口
// @Description 用户通过用户名和密码登录，返回 JWT Token
// @Tags 用户认证
// @Accept application/json
// @Produce application/json
// @Param login body authModel.LoginDto true "登录请求体"
// @Success 200 {object} res.Response "登录成功，返回JWT Token"
// @Failure 200 {object} res.Response "登录失败，返回错误信息"
// @Router /login [post]
func (userHandler *UserHandler) Login() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 从请求体获取用户名和密码
		var loginDto authModel.LoginDto
		if err := ctx.ShouldBindJSON(&loginDto); err != nil {
			return res.Response{
				Msg: commonModel.INVALID_REQUEST_BODY,
				Err: err,
			}
		}

		// 调用 Service 层处理登陆
		token, err := userHandler.userService.Login(&loginDto)
		if err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		// 返回成功响应， 包含 JWT Token
		return res.Response{
			Data: token,
			Msg:  commonModel.LOGIN_SUCCESS,
		}
	})

}

// Register 用户注册
//
// @Summary 用户注册
// @Description 通过提交用户名、密码等信息完成注册
// @Tags 用户认证
// @Accept json
// @Produce json
// @Param register body authModel.RegisterDto true "注册请求体"
// @Success 200 {object} res.Response "注册成功，code=1，msg=REGISTER_SUCCESS"
// @Failure 200 {object} res.Response "请求参数错误或注册失败，code=0，msg错误描述"
// @Router /register [post]
func (userHandler *UserHandler) Register() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		var registerDto authModel.RegisterDto
		if err := ctx.ShouldBindJSON(&registerDto); err != nil {
			return res.Response{
				Msg: commonModel.INVALID_REQUEST_BODY,
				Err: err,
			}
		}

		// 调用 Service 层处理注册
		if err := userHandler.userService.Register(&registerDto); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Msg: commonModel.REGISTER_SUCCESS,
		}
	})

}

// UpdateUser 更新用户信息
//
// @Summary 更新当前用户的信息
// @Description 接口会根据请求体更新用户相关字段，需携带有效的用户身份（如 JWT）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body model.UserInfoDto true "用户更新信息"
// @Success 200 {object} res.Response "更新成功，code=1，msg=UPDATE_USER_SUCCESS"
// @Failure 200 {object} res.Response "请求参数错误或更新失败，code=0，msg错误描述"
// @Security ApiKeyAuth
// @Router /user [put]
func (userHandler *UserHandler) UpdateUser() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 解析用户请求体中的参数
		var userdto model.UserInfoDto
		if err := ctx.ShouldBindJSON(&userdto); err != nil {
			return res.Response{
				Msg: commonModel.INVALID_REQUEST_BODY,
				Err: err,
			}
		}

		// 获取当前用户 ID
		userid := ctx.MustGet("userid").(uint)
		if err := userHandler.userService.UpdateUser(userid, userdto); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Msg: commonModel.UPDATE_USER_SUCCESS,
		}
	})

}

// UpdateUserAdmin 更新用户权限
//
// @Summary 更新用户权限（管理员权限）
// @Description 通过用户ID更新其管理员权限，接口调用者需拥有相应权限
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} res.Response "更新成功，code=1，msg=UPDATE_USER_SUCCESS"
// @Failure 200 {object} res.Response "参数错误或更新失败，code=0，msg错误描述"
// @Security ApiKeyAuth
// @Router /user/admin/{id} [put]
func (userHandler *UserHandler) UpdateUserAdmin() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 获取当前用户 ID
		userid := ctx.MustGet("userid").(uint)

		idStr := ctx.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return res.Response{
				Msg: commonModel.INVALID_PARAMS,
				Err: err,
			}
		}

		if err := userHandler.userService.UpdateUserAdmin(userid, uint(id)); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Msg: commonModel.UPDATE_USER_SUCCESS,
		}
	})

}

// GetAllUsers 获取所有用户
//
// @Summary 获取所有用户
// @Description 获取系统中所有用户的详细信息，接口需要认证
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} res.Response{data=[]model.UserInfoDto} "获取成功，code=1，包含用户列表"
// @Failure 200 {object} res.Response "获取失败，code=0，msg错误描述"
// @Security ApiKeyAuth
// @Router /allusers [get]
func (userHandler *UserHandler) GetAllUsers() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		allusers, err := userHandler.userService.GetAllUsers()
		if err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Data: allusers,
			Msg:  commonModel.GET_USER_SUCCESS,
		}
	})

}

// DeleteUser 删除用户
//
// @Summary 删除用户
// @Description 根据用户ID删除用户，调用者需具备相应权限
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} res.Response "删除成功，code=1，msg=DELETE_USER_SUCCESS"
// @Failure 200 {object} res.Response "参数错误或删除失败，code=0，msg错误描述"
// @Security ApiKeyAuth
// @Router /user/{id} [delete]
func (userHandler *UserHandler) DeleteUser() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 获取当前用户 ID
		userid := ctx.MustGet("userid").(uint)

		idStr := ctx.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return res.Response{
				Msg: commonModel.INVALID_PARAMS,
				Err: err,
			}
		}

		if err := userHandler.userService.DeleteUser(userid, uint(id)); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Msg: commonModel.DELETE_USER_SUCCESS,
		}
	})

}

// GetUserInfo 获取当前用户信息
//
// @Summary 获取当前用户信息
// @Description 获取当前认证用户的详细信息，密码字段不会返回
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} res.Response{data=model.UserInfoDto} "获取成功，code=1，包含用户信息"
// @Failure 200 {object} res.Response "获取失败，code=0，msg错误描述"
// @Security ApiKeyAuth
// @Router /user [get]
func (userHandler *UserHandler) GetUserInfo() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 获取当前用户 ID
		userid := ctx.MustGet("userid").(uint)

		// 调用 Service 层获取用户信息
		user, err := userHandler.userService.GetUserByID(int(userid))
		user.Password = "" // 不返回密码
		if err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		// 返回成功响应
		return res.Response{
			Data: user,
			Msg:  commonModel.GET_USER_INFO_SUCCESS,
		}
	})

}
