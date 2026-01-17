package controller

import (
	"context"
	"encoding/json"
	"garnet-scaff/internal/infrastructures/custresp"
	"garnet-scaff/internal/usecases/request"
	"garnet-scaff/pkg/constants"

	"github.com/google/uuid"
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/response"
	"github.com/nocturna-ta/golib/response/rest"
	"github.com/nocturna-ta/golib/router"
	"github.com/nocturna-ta/golib/tracing"
)

// GetAllUsers godoc
// @Summary GetAllUsers
// @Description GetAllUsers to get all user
// @Tags user
// @Accept json
// @Produce json
// @Success 200	{object}	jsonResponse{data=response.UserResponse}
// @Router /v1/users [get]
func (api *API) GetAllUsers(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.GetAllUsers")
	defer span.End()

	res, err := api.userUc.GetAllUsers(ctx)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}

// GetUser godoc
// @Summary GetUser
// @Description GetUser to get a user
// @Tags user
// @Accept json
// @Produce json
// @Param 		userId 					path 		string 	false	"user id"
// @Success 200	{object}	jsonResponse{data=response.UserResponse}
// @Router /v1/users/{userId} [get]
func (api *API) GetUser(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.GetUser")
	defer span.End()

	id := req.Params("userId")
	if id == constants.EmptyString {
		return custresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "empty user id",
			Code:    400,
			Type:    response.ErrBadRequest,
		})
	}

	userId, err := uuid.Parse(id)
	if err != nil {
		return custresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "invalid user id",
			Code:    400,
			Type:    response.ErrBadRequest,
			Cause:   err,
		})
	}
	res, err := api.userUc.GetUserByID(ctx, userId)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}

// CreateUser godoc
// @Summary CreateUser
// @Description CreateUser to create a user
// @Tags user
// @Accept json
// @Produce json
// @Param		user					body		request.UserRegistrationRequest	true	"Create User payload"
// @Success 200	{object}	jsonResponse{data=response.UserResponse}
// @Router /v1/users [post]
func (api *API) CreateUser(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.CreateUser")
	defer span.End()

	var userReq request.UserRegistrationRequest
	err := json.Unmarshal(req.RawBody(), &userReq)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	_, err = api.userUc.RegisterUser(ctx, &userReq)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetMessage("User Created"), nil
}

// UpdateUser godoc
// @Summary UpdateUser
// @Description UpdateUser to update a user
// @Tags user
// @Accept json
// @Produce json
// @Param		user					body		request.UserUpdateRequest	true	"Update User payload"
// @Success 200	{object}	jsonResponse{data=response.UserResponse}
// @Router /v1/users [patch]
func (api *API) UpdateUser(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.UpdateUser")
	defer span.End()

	var userReq request.UserUpdateRequest
	err := json.Unmarshal(req.RawBody(), &userReq)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	res, err := api.userUc.UpdateUser(ctx, &userReq)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res).SetMessage("User Update Successfully"), nil
}

// DeleteUser godoc
// @Summary DeleteUser
// @Description DeleteUser to delete a user
// @Tags user
// @Accept json
// @Produce json
// @Param 		userId 					path 		string 	false	"user id"
// @Success 200	{object}	jsonResponse{data=response.UserResponse}
// @Router /v1/users/{userId} [delete]
func (api *API) DeleteUser(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.DeleteUser")
	defer span.End()

	id := req.Params("userId")
	if id == constants.EmptyString {
		return custresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "empty user id",
			Code:    400,
			Type:    response.ErrBadRequest,
		})

	}

	userId, err := uuid.Parse(id)
	if err != nil {
		return custresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "invalid user id",
			Code:    400,
			Type:    response.ErrBadRequest,
			Cause:   err,
		})
	}

	err = api.userUc.SoftDeleteUser(ctx, userId)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetMessage("Delete User Succefully"), nil
}
