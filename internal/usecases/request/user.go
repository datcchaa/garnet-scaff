package request

import (
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/response"
)

type UserRegistrationRequest struct {
	Username string `json:"username"`
}

type UserUpdateRequest struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (req *UserRegistrationRequest) ValidateRegistrationRequest() error {
	if req == nil {
		return &custerr.ErrChain{
			Message: "Request cannot be nil",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	return nil
}

func (req *UserUpdateRequest) ValidateUpdateRequest() error {
	if req == nil {
		return &custerr.ErrChain{
			Message: "Request cannot be nil",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	return nil
}
