package web

import (
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"log"
	"my-go-basic-study/webook/internal/domain"
	"my-go-basic-study/webook/internal/service"
	"net/http"
)

var _ handler = &UserHandler{}

type UserHandler struct {
	svc                service.UserService
	emailExpression    *regexp.Regexp
	passwordExpression *regexp.Regexp
}

func NewUserHandler(svc service.UserService) *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	emailExpression := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExpression := regexp.MustCompile(passwordRegexPattern, regexp.None)

	return &UserHandler{
		svc:                svc,
		emailExpression:    emailExpression,
		passwordExpression: passwordExpression,
	}

}

func (u *UserHandler) RegisterRouters(server *gin.Engine) {

	group := server.Group("/users")

	group.POST("/signup", u.Signup)
	group.POST("/login", u.Login)
	group.POST("/loginJWT", u.LoginJWT)
	group.POST("/edit", u.Edit)
	group.GET("/profile", u.Profile)
	group.GET("/profileJWT", u.ProfileJWT)
	group.GET("/logout", u.Logout)

}

func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	loginUser, err := u.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "用户名或者密码错误")
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, "系统错误")
		return
	}

	//这里使用jwt而不是session了
	userClaim := UserClaims{
		Uid: loginUser.Id,
		/*RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},*/
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
	signedString, err := token.SignedString([]byte("a59b1c734e8f2d5a967b3c841e5f9a2d6b4"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	ctx.Header("x-jwt-token", signedString)
	ctx.JSON(http.StatusOK, "登录成功～")
	return

}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid int64 `json:"uid"`
}

func (u *UserHandler) Signup(ctx *gin.Context) {

	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	matchString, err := u.emailExpression.MatchString(req.Email)
	if err != nil {
		ctx.JSON(http.StatusOK, "系统错误")
		return
	}
	if !matchString {
		ctx.String(http.StatusOK, "邮件格式错误")
		return
	}

	matchString, err = u.passwordExpression.MatchString(req.Password)
	if err != nil {
		ctx.JSON(http.StatusOK, "系统错误")
		return
	}
	if !matchString {
		ctx.String(http.StatusOK, "密码格式错误")
		return
	}
	err = u.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrUserDuplicateEmail) {
		ctx.String(http.StatusOK, "邮箱已存在")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
	}
	ctx.String(http.StatusOK, "signup success")
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	loginUser, err := u.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "用户名或者密码错误")
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, "系统错误")
		return
	}
	//登录成功 在这里设置session
	sess := sessions.Default(ctx)
	sess.Set("userId", loginUser.Id)
	sess.Options(sessions.Options{
		MaxAge: 30 * 60,
	})
	sess.Save()

	ctx.JSON(http.StatusOK, "登录成功～")
	return

}

func (u *UserHandler) Edit(ctx *gin.Context) {

}

func (u *UserHandler) Profile(ctx *gin.Context) {

	ctx.String(http.StatusOK, "这是你的profile")
}

func (u *UserHandler) ProfileJWT(ctx *gin.Context) {

	value, exists := ctx.Get("claims")
	//可以断定必然有exists
	if !exists {
		log.Printf("ProfileJWT claims not found")
		ctx.String(http.StatusOK, "系统错误～～")
		return
	}
	claims, ok := value.(UserClaims)
	if !ok {
		log.Printf("ProfileJWT claims not found")
		ctx.String(http.StatusOK, "系统错误～～")
		return
	}
	ctx.String(http.StatusOK, "这是你的profile,you ID", claims.Uid)
}
func (u *UserHandler) Logout(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Options(sessions.Options{
		MaxAge: -1,
	})
	sess.Save()
	ctx.String(http.StatusOK, "推出登录成功～")
}
