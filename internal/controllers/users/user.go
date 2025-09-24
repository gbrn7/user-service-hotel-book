package users

import (
	"net/http"
	"user-service/domain/dto"
	"user-service/helpers"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (u *UserController) SignUp(ctx *gin.Context) {
	request := &dto.RegisterRequest{}

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		helpers.HttpResponse(helpers.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := helpers.ErrValidationResponse(err)
		helpers.HttpResponse(helpers.ParamHTTPResp{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Err:     err,
			Gin:     ctx,
		})

		return
	}

	user, err := u.service.Register(ctx, request)
	if err != nil {
		helpers.HttpResponse(helpers.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	helpers.HttpResponse(helpers.ParamHTTPResp{
		Code: http.StatusOK,
		Data: user.User,
		Gin:  ctx,
	})
}

func (u *UserController) SignIn(ctx *gin.Context) {
	request := &dto.LoginRequest{}

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		helpers.HttpResponse(helpers.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := helpers.ErrValidationResponse(err)
		helpers.HttpResponse(helpers.ParamHTTPResp{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Err:     err,
			Gin:     ctx,
		})

		return
	}

	user, err := u.service.Login(ctx, request)
	if err != nil {
		helpers.HttpResponse(helpers.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	helpers.HttpResponse(helpers.ParamHTTPResp{
		Code:  http.StatusOK,
		Data:  user.User,
		Token: &user.Token,
		Gin:   ctx,
	})
}

func (u *UserController) GetUserLogin(ctx *gin.Context) {
	user, err := u.service.GetUserLogin(ctx.Request.Context())
	if err != nil {
		helpers.HttpResponse(helpers.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})

		return
	}

	helpers.HttpResponse(helpers.ParamHTTPResp{
		Code: http.StatusOK,
		Data: user,
		Gin:  ctx,
	})
}

func (u *UserController) GetUserByUUID(ctx *gin.Context) {
	user, err := u.service.GetUserByUUID(ctx, ctx.Param("uuid"))
	if err != nil {
		helpers.HttpResponse(helpers.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})

		return
	}

	helpers.HttpResponse(helpers.ParamHTTPResp{
		Code: http.StatusOK,
		Data: user,
		Gin:  ctx,
	})
}

func (u *UserController) GetAllCustomer(ctx *gin.Context) {
	user, err := u.service.GetAllCustomer(ctx)
	if err != nil {
		helpers.HttpResponse(helpers.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})

		return
	}

	helpers.HttpResponse(helpers.ParamHTTPResp{
		Code: http.StatusOK,
		Data: user,
		Gin:  ctx,
	})

}

func (u *UserController) GetAllAdmin(ctx *gin.Context) {
	user, err := u.service.GetAllAdmin(ctx)
	if err != nil {
		helpers.HttpResponse(helpers.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})

		return
	}

	helpers.HttpResponse(helpers.ParamHTTPResp{
		Code: http.StatusOK,
		Data: user,
		Gin:  ctx,
	})

}

func (u *UserController) GetAllUser(ctx *gin.Context) {
	user, err := u.service.GetAllUser(ctx)
	if err != nil {
		helpers.HttpResponse(helpers.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})

		return
	}

	helpers.HttpResponse(helpers.ParamHTTPResp{
		Code: http.StatusOK,
		Data: user,
		Gin:  ctx,
	})

}

func (u *UserController) Update(ctx *gin.Context) {
	request := &dto.UpdateRequest{}
	uuid := ctx.Param("uuid")

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		helpers.HttpResponse(helpers.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})

		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := helpers.ErrValidationResponse(err)
		helpers.HttpResponse(helpers.ParamHTTPResp{
			Code:    http.StatusBadRequest,
			Err:     err,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     ctx,
		})

		return
	}

	user, err := u.service.Update(ctx, request, uuid)
	if err != nil {
		helpers.HttpResponse(helpers.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})

		return
	}

	helpers.HttpResponse(helpers.ParamHTTPResp{
		Code: http.StatusCreated,
		Data: user,
		Gin:  ctx,
	})
}
