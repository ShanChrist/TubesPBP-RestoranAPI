package model

type Courrier struct {
	CourrierID                int `json:"courrier_id"`
	CourrierUserI             int `json:"user_id"`
	CourrierBalance           int `json:"balance"`
	CourrierStatusPengantaran int `json:"status_pengantaran"`
}
