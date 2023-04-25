package model

type Member struct {
	MemberID      int `json:"member_id"`
	MemberUserID  int `json:"user_id"`
	MemberBalance int `json:"balance"`
	MemberPoint   int `json:"point"`
}
