package web

import (
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"learn-geektime-basic-go/webook/internal/domain"
	"learn-geektime-basic-go/webook/internal/service"
	"net/http"
	"net/mail"
	"time"
)

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

type UserHandler struct {
	svc        *service.UserService
	passwordRe *regexp.Regexp
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid int64
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	re := regexp.MustCompile(passwordRegexPattern, 0)

	return &UserHandler{
		svc:        svc,
		passwordRe: re,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	group := server.Group("/users")
	group.POST("/signup", u.SignUp)
	//group.POST("/login", u.Login)
	group.POST("/login", u.LoginJWT)
	group.POST("/edit", u.Edit)
	group.GET("/profile", u.Profile)
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}

	var req SignUpReq

	if err := ctx.Bind(&req); err != nil {
		return
	}

	if !isValidEmail(req.Email) {
		ctx.String(http.StatusOK, "你的邮箱格式不对")
		return
	}

	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
	}

	ok, err := u.passwordRe.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8位，包含特殊字符")
	}
	err = u.svc.Signup(ctx.Request.Context(), domain.User{Email: req.Email, Password: req.ConfirmPassword})

	if errors.Is(err, service.ErrUserDuplicateEmail) {
		ctx.String(http.StatusOK, "重复邮箱，请换一个邮箱")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	ctx.String(http.StatusOK, "hello, 注册成功")
	return
}

func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if ctx.Bind(&req) != nil {
		return
	}

	user, err := u.svc.Login(ctx.Request.Context(), req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	userClaims := UserClaims{
		Uid: user.Id,
	}

	userClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour))

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, userClaims)

	signedString, err := token.SignedString([]byte("zonUXJUU5%XmP6wkH^X%W7l%sNM0dPvI"))

	//key, err := rsa.GenerateKey(rand.Reader, 2048)
	//
	//token := jwt.NewWithClaims(jwt.SigningMethodRS512, userClaims)

	//signedString, err := token.SignedString(key)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
	}

	ctx.Header("x-jwt-token", signedString)

	fmt.Println(signedString)
	fmt.Println(user)

	ctx.String(http.StatusOK, "登录成功")

	return
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if ctx.Bind(&req) != nil {
		return
	}

	user, err := u.svc.Login(ctx.Request.Context(), req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	sess := sessions.Default(ctx)

	sess.Set("userId", user.Id)

	err = sess.Save()
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.String(http.StatusOK, "登录成功")

	return
}

func (u *UserHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		Nickname string    `json:"nickname"`
		Phone    string    `json:"phone"`
		AboutMe  string    `json:"aboutMe"`
		Birthday time.Time `json:"birthday"`
	}

	var req EditReq

	if err := ctx.Bind(&req); err != nil {
		return
	}

	// TODO
}

func (u *UserHandler) Profile(ctx *gin.Context) {
	value, exists := ctx.Get("claims")
	if !exists {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	claims, ok := value.(*UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
	}

	fmt.Println("userId: ", claims.Uid)

	sess := sessions.Default(ctx)
	userId := sess.Get("userId")
	if userId == nil {
		return
	}
	userid := userId.(int64)
	fmt.Println(userid)
	ctx.String(http.StatusOK, "%d", userid)
}
