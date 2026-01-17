package errorcode

type ErrorDefinition struct {
	Code    int
	Message string
}

var (
	NotFound = ErrorDefinition{
		Code:    40005,
		Message: "Not Found",
	}
	UserAlreadyExists = ErrorDefinition{
		Code:    40011,
		Message: "User already exists",
	}
	InvalidUUID = ErrorDefinition{
		Code:    40012,
		Message: "Invalid UUID",
	}
)
