package domain

import "errors"

var (
	// 사용자 관련 에러
	ErrUserNotFound       = errors.New("사용자를 찾을 수 없습니다")
	ErrEmailAlreadyExists = errors.New("이미 존재하는 이메일입니다")
	ErrInvalidCredentials = errors.New("잘못된 인증 정보입니다")

	// 토큰 관련 에러
	ErrInvalidToken = errors.New("유효하지 않은 토큰입니다")
	ErrExpiredToken = errors.New("만료된 토큰입니다")
	ErrRevokedToken = errors.New("폐기된 토큰입니다")
)
