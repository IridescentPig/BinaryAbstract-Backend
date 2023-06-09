package api

import (
	"asset-management/app/dao"
	"asset-management/app/define"
	"asset-management/app/model"
	"asset-management/middleware"
	"asset-management/utils"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	headerForm = map[string]string{
		"Content-Type": "x-www-form-urlencoded",
	}
	headerJson = map[string]string{
		"Content-Type": "application/json",
	}
	headerFormToken = map[string]string{
		"Content-Type":  "x-www-form-urlencoded",
		"Authorization": "",
	}
	headerJsonToken = map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "",
	}
)

func InitForTest(r *gin.Engine) {
	InitForUser(r)
	InitForEntity(r)
	InitForAssetClass(r)
	InitForAsset(r)

	dao.InitForTest()
}

func InitForUser(r *gin.Engine) {
	group := r.Group("/user")

	group.POST("/register", utils.Handler(UserApi.UserRegister))
	group.POST("/login", utils.Handler(UserApi.UserLogin))
	group.GET("/logout", utils.Handler(middleware.JWTMiddleware()), utils.Handler(UserApi.UserLogout))
	group.POST("", utils.Handler(middleware.JWTMiddleware()), utils.Handler(UserApi.UserCreate))
	group.PATCH("/:username", utils.Handler(middleware.JWTMiddleware()), utils.Handler(UserApi.ResetContent))
	group.GET("/:username/lock", utils.Handler(middleware.JWTMiddleware()), utils.Handler(UserApi.LockUser))
	group.GET("/:username/unlock", utils.Handler(middleware.JWTMiddleware()), utils.Handler(UserApi.UnlockUser))

	group.GET("/info/:user_id", utils.Handler(middleware.JWTMiddleware()), utils.Handler(UserApi.GetUserInfoByID))
	group.GET("/list", utils.Handler(middleware.JWTMiddleware()), utils.Handler(UserApi.GetAllUsers))
	group.DELETE("/:user_id", utils.Handler(middleware.JWTMiddleware()), utils.Handler(UserApi.DeleteUser))
	group.POST("/info/:user_id/password", utils.Handler(middleware.JWTMiddleware()), utils.Handler(UserApi.ChangePassword))
	group.DELETE("/info/:user_id/entity", utils.Handler(middleware.JWTMiddleware()), utils.Handler(middleware.CheckSystemSuper()), utils.Handler(UserApi.ChangeUserEntity))
	group.DELETE("/info/:user_id/department", utils.Handler(middleware.JWTMiddleware()), utils.Handler(UserApi.ChangeUserDepartment))
}

func GetJsonBody(data interface{}) io.Reader {
	bodyData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	body := bytes.NewReader(bodyData)
	return body
}

func GetFormBody(data interface{}) io.Reader {
	bodyData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	m := make(map[string]string)
	json.Unmarshal(bodyData, &m)

	body := url.Values{}
	for k, v := range m {
		body.Set(k, v)
	}

	return strings.NewReader(body.Encode())
}

func GetRequest(method string, url string, header map[string]string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	return req
}

func TestUser(t *testing.T) {
	res := httptest.NewRecorder()
	_, r := gin.CreateTestContext(res)
	InitForTest(r)

	UserRegister := define.UserRegisterReq{
		UserName: "test",
		Password: "e10adc3949ba59abbe56e057f20f883e",
	}

	{
		req := GetRequest(http.MethodPost, "/user/register", headerForm, GetFormBody(UserRegister))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	{
		req := GetRequest(http.MethodPost, "/user/register", headerJson, GetJsonBody(UserRegister))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

	{
		req := GetRequest(http.MethodPost, "/user/register", headerJson, GetJsonBody(UserRegister))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	UserLogin := define.UserLoginReq{
		UserName: "test",
		Password: "e10adc3949ba59abbe56e057f20f883e",
	}

	{
		req := GetRequest(http.MethodPost, "/user/login", headerForm, GetFormBody(UserLogin))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	{
		req := GetRequest(http.MethodPost, "/user/login", headerJson, GetJsonBody(UserLogin))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

	b, err := io.ReadAll(res.Result().Body)
	if err != nil {
		log.Fatal(err)
	}

	{
		data := map[string]interface{}{}
		json.Unmarshal(b, &data)
		user := data["data"].(map[string]interface{})
		token := user["token"].(string)

		headerJsonToken["Authorization"] = token
		req := GetRequest(http.MethodGet, "/user/logout", headerJsonToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

}

func TestAdmin(t *testing.T) {
	res := httptest.NewRecorder()
	_, r := gin.CreateTestContext(res)
	InitForTest(r)

	admin := model.User{
		UserName:    "admin",
		Password:    utils.CreateMD5("21232f297a57a5a743894a0e4a801fc3"),
		SystemSuper: true,
		Ban:         false,
		// EntityID:     nil,
		// DepartmentID: nil,
	}
	dao.UserDao.Create(admin)

	UserLogin := define.UserLoginReq{
		UserName: "admin",
		Password: "21232f297a57a5a743894a0e4a801fc3",
	}

	{
		req := GetRequest(http.MethodPost, "/user/login", headerJson, GetJsonBody(UserLogin))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

	b, err := io.ReadAll(res.Result().Body)
	if err != nil {
		log.Fatal(err)
	}

	data := map[string]interface{}{}
	json.Unmarshal(b, &data)
	// fmt.Println(data["data"])
	user := data["data"].(map[string]interface{})
	token := user["token"].(string)
	headerJsonToken["Authorization"] = token
	headerFormToken["Authorization"] = token

	UserCreate := define.UserRegisterReq{
		UserName: "test2",
		Password: "e10adc3949ba59abbe56e057f20f883e",
	}

	{
		req := GetRequest(http.MethodPost, "/user", headerFormToken, GetFormBody(UserCreate))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	{
		req := GetRequest(http.MethodPost, "/user", headerJson, GetJsonBody(UserCreate))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}

	{
		req := GetRequest(http.MethodPost, "/user", headerJsonToken, GetJsonBody(UserCreate))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

	{
		req := GetRequest(http.MethodPost, "/user", headerJsonToken, GetJsonBody(UserCreate))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	UserLogin = define.UserLoginReq{
		UserName: "test2",
		Password: "e10adc3949ba59abbe56e057f20f883e",
	}

	{
		{
			req := GetRequest(http.MethodPost, "/user/login", headerJson, GetJsonBody(UserLogin))
			res = httptest.NewRecorder()
			r.ServeHTTP(res, req)

			assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
		}

		b, err := io.ReadAll(res.Result().Body)
		if err != nil {
			log.Fatal(err)
		}

		data := map[string]interface{}{}
		json.Unmarshal(b, &data)
		// fmt.Println(data["data"])
		user := data["data"].(map[string]interface{})
		tokenNotAdmin := user["token"].(string)
		headerJsonToken["Authorization"] = tokenNotAdmin
		headerFormToken["Authorization"] = tokenNotAdmin

		UserCreate := define.UserRegisterReq{
			UserName: "test3",
			Password: "e10adc3949ba59abbe56e057f20f883e",
		}

		UserReset := define.ResetReq{
			Method:   0,
			Identity: 1,
			Password: "",
		}

		{
			req := GetRequest(http.MethodPost, "/user", headerJsonToken, GetJsonBody(UserCreate))
			res = httptest.NewRecorder()
			r.ServeHTTP(res, req)

			assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
		}

		{
			req := GetRequest(http.MethodPatch, "/user/test2", headerJsonToken, GetJsonBody(UserReset))
			res = httptest.NewRecorder()
			r.ServeHTTP(res, req)

			assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
		}

		UserReset = define.ResetReq{
			Method:   1,
			Identity: 1,
			Password: "e10adc3949ba59abbe56e057f20f883e",
		}

		{
			req := GetRequest(http.MethodPatch, "/user/test3", headerJsonToken, GetJsonBody(UserReset))
			res = httptest.NewRecorder()
			r.ServeHTTP(res, req)

			assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
		}

		{
			req := GetRequest(http.MethodPatch, "/user/admin", headerJsonToken, GetJsonBody(UserReset))
			res = httptest.NewRecorder()
			r.ServeHTTP(res, req)

			assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
		}

		{
			req := GetRequest(http.MethodGet, "/user/admin/lock", headerJsonToken, nil)
			res = httptest.NewRecorder()
			r.ServeHTTP(res, req)

			assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
		}

		{
			req := GetRequest(http.MethodGet, "/user/admin/lock", headerJsonToken, nil)
			res = httptest.NewRecorder()
			r.ServeHTTP(res, req)

			assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
		}

		{
			req := GetRequest(http.MethodGet, "/user/admin/unlock", headerJsonToken, nil)
			res = httptest.NewRecorder()
			r.ServeHTTP(res, req)

			assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, "response failed")
		}
	}

	headerJsonToken["Authorization"] = token
	headerFormToken["Authorization"] = token
	UserReset := define.ResetReq{
		Method:   0,
		Identity: 1,
		Password: "",
	}

	{
		req := GetRequest(http.MethodPatch, "/user/test2", headerJsonToken, GetJsonBody(UserReset))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

	{
		req := GetRequest(http.MethodPatch, "/user/test3", headerJsonToken, GetJsonBody(UserReset))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	UserReset = define.ResetReq{
		Method:   1,
		Identity: 1,
		Password: "e10adc3949ba59abbe56e057f20f883e",
	}

	{
		req := GetRequest(http.MethodPatch, "/user/test2", headerJsonToken, GetJsonBody(UserReset))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

	{
		req := GetRequest(http.MethodPatch, "/user/test3", headerJsonToken, GetJsonBody(UserReset))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	{
		req := GetRequest(http.MethodPatch, "/user/test2", headerJsonToken, GetFormBody(UserReset))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	{
		UserReset.Method = 5
		req := GetRequest(http.MethodPatch, "/user/test2", headerJsonToken, GetJsonBody(UserReset))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	{
		req := GetRequest(http.MethodGet, "/user/test2/lock", headerJsonToken, GetJsonBody(UserReset))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

	UserLogin = define.UserLoginReq{
		UserName: "test2",
		Password: "e10adc3949ba59abbe56e057f20f883e",
	}

	{
		req := GetRequest(http.MethodPost, "/user/login", headerJsonToken, GetJsonBody(UserLogin))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	{
		UserLogin.Password = "834932849382"
		req := GetRequest(http.MethodPost, "/user/login", headerJsonToken, GetJsonBody(UserLogin))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	{
		req := GetRequest(http.MethodGet, "/user/test3/lock", headerJsonToken, GetJsonBody(UserReset))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	{
		req := GetRequest(http.MethodGet, "/user/test2/unlock", headerJsonToken, GetJsonBody(UserReset))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

	{
		req := GetRequest(http.MethodGet, "/user/test3/unlock", headerJsonToken, GetJsonBody(UserReset))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	/*UserLogin = define.UserLoginReq{
		UserName: "admin",
		Password: "21232f297a57a5a743894a0e4a801fc3",
	}

	{
		req := GetRequest(http.MethodPost, "/user/login", headerJson, GetJsonBody(UserLogin))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode, print_errormessage(res))
	}

	b, err = io.ReadAll(res.Result().Body)
	if err != nil {
		log.Fatal(err)
	}

	data = map[string]interface{}{}
	json.Unmarshal(b, &data)
	// fmt.Println(data["data"])
	user = data["data"].(map[string]interface{})
	token = user["token"].(string)
	headerJsonToken["Authorization"] = token
	headerFormToken["Authorization"] = token*/

	{
		req := GetRequest(http.MethodDelete, "/user/1", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, print_errormessage(res))
	}
	{
		req := GetRequest(http.MethodDelete, "/user/9", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodDelete, "/user/1", headerForm, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

	{
		req := GetRequest(http.MethodGet, "/user/info/1", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

	{
		req := GetRequest(http.MethodGet, "/user/info/9", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodGet, "/user/info/1", headerForm, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodGet, "/user/list", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	password := define.ChangePasswordReq{
		Password: "pp",
	}
	{
		req := GetRequest(http.MethodPost, "/user/info/1/password", headerFormToken, GetJsonBody(password))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/user/info/9/password", headerFormToken, GetJsonBody(password))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/user/info/1/password", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/user/info/1/password", headerForm, GetJsonBody(password))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	entity := define.ChangeUserEntityReq{
		EntityID: 1,
	}
	{
		req := GetRequest(http.MethodPost, "/user/info/1/entity", headerFormToken, GetJsonBody(entity))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/user/info/1/entity", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/user/info/9/entity", headerFormToken, GetJsonBody(entity))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/user/info/1/entity", headerForm, GetJsonBody(entity))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	department := define.ChangeUserDepartmentReq{
		DepartmentID: 1,
	}
	{
		req := GetRequest(http.MethodPost, "/user/info/1/department", headerFormToken, GetJsonBody(department))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/user/info/9/department", headerFormToken, GetJsonBody(department))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/user/info/1/department", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/user/info/1/department", headerForm, GetJsonBody(department))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
}

func TestUserInfo(t *testing.T) {
	res := httptest.NewRecorder()
	_, r := gin.CreateTestContext(res)
	InitForTest(r)

	admin := model.User{
		UserName:        "admin",
		Password:        utils.CreateMD5("21232f297a57a5a743894a0e4a801fc3"),
		SystemSuper:     true,
		EntitySuper:     true,
		DepartmentSuper: true,
		Ban:             false,
		// EntityID:     nil,
		// DepartmentID: nil,
	}
	dao.UserDao.Create(admin)

	UserLogin := define.UserLoginReq{
		UserName: "admin",
		Password: "21232f297a57a5a743894a0e4a801fc3",
	}

	{
		req := GetRequest(http.MethodPost, "/user/login", headerJson, GetJsonBody(UserLogin))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

	b, err := io.ReadAll(res.Result().Body)
	if err != nil {
		log.Fatal(err)
	}

	data := map[string]interface{}{}
	json.Unmarshal(b, &data)
	// fmt.Println(data["data"])
	user := data["data"].(map[string]interface{})
	token := user["token"].(string)
	headerJsonToken["Authorization"] = token
	headerFormToken["Authorization"] = token

	// GET /info/:user_id
	{
		req := GetRequest(http.MethodGet, "/user/info/1", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodGet, "/user/info/1", headerJson, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodGet, "/user/info/2", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	// GET /list
	{
		req := GetRequest(http.MethodGet, "/user/list", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		// assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodGet, "/user/list", headerForm, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		// assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}

	// DELETE /:user_id
	{
		req := GetRequest(http.MethodDelete, "/user/1", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodDelete, "/user/1", headerForm, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodDelete, "/user/2", headerFormToken, nil)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		//assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	// POST /info/:user_id/password
	ChangePassword := define.ChangePasswordReq{
		Password: "newpass",
	}
	{
		req := GetRequest(http.MethodPost, "/user/info/1/password", headerFormToken, GetJsonBody(ChangePassword))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/user/info/1/password", headerForm, GetJsonBody(ChangePassword))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodPost, "/user/info/2/password", headerFormToken, GetJsonBody(ChangePassword))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	// DELETE /info/:user_id/entity
	ChangeEntity := define.ChangeUserEntityReq{
		EntityID: 1,
	}
	{
		req := GetRequest(http.MethodDelete, "/user/info/1/entity", headerFormToken, GetJsonBody(ChangeEntity))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodDelete, "/user/info/1/entity", headerForm, GetJsonBody(ChangeEntity))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodDelete, "/user/info/2/entity", headerFormToken, GetJsonBody(ChangeEntity))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	// DELETE /info/:user_id/department
	ChangeDepartment := define.ChangeUserDepartmentReq{
		DepartmentID: 1,
	}
	{
		req := GetRequest(http.MethodDelete, "/user/info/1/department", headerForm, GetJsonBody(ChangeDepartment))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}
	{
		req := GetRequest(http.MethodDelete, "/user/info/2/department", headerFormToken, GetJsonBody(ChangeDepartment))
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}
}
