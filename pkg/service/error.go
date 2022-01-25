package service

import "errors"

// Errors that can occur in services.
var (
	ErrUnauthorized              = errors.New("unauthorized")
	ErrForbidden                 = errors.New("forbidden")
	ErrNotAccessRequestRequester = errors.New("service is not access request requester")

	ErrNotGrantable    = errors.New("access request is not grantable")
	ErrNotDeniable     = errors.New("access request is not deniable")
	ErrNotClarifiable  = errors.New("access request is not clarifiable")
	ErrNotWithdrawable = errors.New("access request is not withdrawable")

	ErrServiceNotFound     = errors.New("service not found")
	ErrSameService         = errors.New("requesting service is same as recipient service")
	ErrScopeAlreadyPending = errors.New("at least one scope is already requested")

	ErrInvitedUserExists = errors.New("user with the provided email already exists")
	ErrInvitationExists  = errors.New("invitation for user with the provided email already exists")

	ErrNonSamlUserForbidden   = errors.New("only SAML/AD users can change their avatars")
	ErrPictureSizeExceeded    = errors.New("picture size should be less than 2MB")
	ErrUnsupportedPictureType = errors.New("picture type should be JPG or PNG")
)
