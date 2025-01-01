package domain

import "errors"

var (
	// 토큰 관련 에러
	ErrInvalidToken = errors.New("유효하지 않은 토큰입니다")
	ErrExpiredToken = errors.New("만료된 토큰입니다")
	ErrRevokedToken = errors.New("폐기된 토큰입니다")

	// 인증 관련 에러
	ErrAuthenticationFailed = errors.New("인증에 실패했습니다")
	ErrUnauthorized         = errors.New("권한이 없습니다")
	ErrInvalidCredentials   = errors.New("잘못된 인증 정보입니다")
)
