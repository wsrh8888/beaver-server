// Code generated by goctl. DO NOT EDIT.
package types

type AuthenticationReq struct {
	Token     string `header:"Token,optional"`
	ValidPath string `header:"ValidPath,optional"`
}

type AuthenticationRes struct {
	UserID string `json:"userId"`
}

type LoginReq struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type LoginRes struct {
	Token  string `json:"token"`
	UserID string `json:"userId"`
}

type RegisterReq struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type RegisterRes struct {
}
