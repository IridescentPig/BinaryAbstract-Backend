package define

import "github.com/dgrijalva/jwt-go"

/*
.*Req struct are strictly defined according to the api
*/
type UserRegisterReq struct {
	UserName string `form:"userName" binding:"required" json:"userName" map:"userName,omitempty"`
	Password string `form:"password" binding:"required" json:"password" map:"password,omitempty"`
}

type UserLoginReq struct {
	UserName string `form:"userName" binding:"required" json:"userName" map:"userName,omitempty"`
	Password string `form:"password" binding:"required" json:"password" map:"password,omitempty"`
}

type UriInfo struct {
	UserName string `uri:"userName"`
}

type ResetReq struct {
	Method   int    `json:"method"`
	Identity int    `json:"identity"`
	Password string `json:"password"`
}

/*
Basic info of user, can be included in other info struct
*/
type UserBasicInfo struct {
	UserID          uint   `json:"user_id"`
	UserName        string `json:"username"`
	EntitySuper     bool   `json:"entity_super"`
	DepartmentSuper bool   `json:"department_super"`
	SystemSuper     bool   `json:"system_super"`
	EntityID        uint   `json:"entity_id"`
	DepartmentID    uint   `json:"department_id"`
}

/*
Used for jwt claims
*/
type UserClaims struct {
	UserBasicInfo
	jwt.StandardClaims
}

type UserInfoResponse struct {
	UserBasicInfo
}
