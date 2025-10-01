package users

import (
	"net/http"
	"user-service/domain/dto"
	"user-service/helpers"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/v2"
)

func (u *UserController) SignUp(ctx *gin.Context) {
	span, spanctx := apm.StartSpan(ctx.Request.Context(), "UserController.SignUp", "controller")
	defer span.End()

	request := &dto.RegisterRequest{}

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		logrus.Errorf("error binding JSON: %s", err)
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
		logrus.Errorf("validation error: %s", err)
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

	user, err := u.service.Register(spanctx, request)
	if err != nil {
		logrus.Errorf("error registering user: %s", err)
		helpers.HttpResponse(helpers.ParamHTTPResp{
			Code: http.StatusInternalServerError,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	helpers.HttpResponse(helpers.ParamHTTPResp{
		Code: http.StatusCreated,
		Data: user.User,
		Gin:  ctx,
	})
}

func (u *UserController) SignIn(ctx *gin.Context) {
	span, spanctx := apm.StartSpan(ctx.Request.Context(), "UserController.SignIn", "controller")
	defer span.End()
	request := &dto.LoginRequest{}

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		logrus.Errorf("error binding JSON: %s", err)
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
		logrus.Errorf("validation error: %s", err)
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

	user, err := u.service.Login(spanctx, request)
	if err != nil {
		logrus.Errorf("error logging in: %s", err)
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
	span, spanctx := apm.StartSpan(ctx.Request.Context(), "UserController.GetUserLogin", "controller")
	defer span.End()

	user, err := u.service.GetUserLogin(spanctx)
	if err != nil {
		logrus.Errorf("error getting user login: %s", err)
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
	span, spanctx := apm.StartSpan(ctx.Request.Context(), "UserController.GetUserByUUID", "controller")
	defer span.End()

	user, err := u.service.GetUserByUUID(spanctx, ctx.Param("uuid"))
	if err != nil {
		logrus.Errorf("error getting user by UUID: %s", err)
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
	span, spanctx := apm.StartSpan(ctx.Request.Context(), "UserController.GetAllCustomer", "controller")
	defer span.End()

	user, err := u.service.GetAllCustomer(spanctx)
	if err != nil {
		logrus.Errorf("error getting all customers: %s", err)
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
	span, spanctx := apm.StartSpan(ctx.Request.Context(), "UserController.GetAllAdmin", "controller")
	defer span.End()

	user, err := u.service.GetAllAdmin(spanctx)
	if err != nil {
		logrus.Errorf("error getting all admins: %s", err)
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
	span, spanctx := apm.StartSpan(ctx.Request.Context(), "UserController.GetAllUser", "controller")
	defer span.End()

	user, err := u.service.GetAllUser(spanctx)
	if err != nil {
		logrus.Errorf("error getting all users: %s", err)
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
	span, spanctx := apm.StartSpan(ctx.Request.Context(), "UserController.Update", "controller")
	defer span.End()

	request := &dto.UpdateRequest{}
	uuid := ctx.Param("uuid")

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		logrus.Errorf("error binding JSON: %s", err)
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
		logrus.Errorf("validation error: %s", err)
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

	user, err := u.service.Update(spanctx, request, uuid)
	if err != nil {
		logrus.Errorf("error updating user: %s", err)
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
