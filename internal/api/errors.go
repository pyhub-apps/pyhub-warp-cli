package api

import "errors"

var (
	// ErrNoAPIKey indicates that no API key is configured
	ErrNoAPIKey = errors.New("API 키가 설정되지 않았습니다")
	
	// ErrNotImplemented indicates that a feature is not implemented
	ErrNotImplemented = errors.New("기능이 구현되지 않았습니다")
	
	// ErrInvalidAPIType indicates an invalid API type was specified
	ErrInvalidAPIType = errors.New("잘못된 API 타입입니다")
)