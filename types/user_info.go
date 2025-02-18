package types

type UserInfo struct {
	ID   uint64 `json:"id"`
	Role Role   `json:"role"`
}

type Role uint8
