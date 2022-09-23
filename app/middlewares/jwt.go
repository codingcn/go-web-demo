package middlewares

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go-web-demo/kernel/tconfig"
	"go-web-demo/kernel/zlog"
	"go.uber.org/zap"
	"net/http"
	"regexp"
	"time"
)

type UserInfoStruct struct {
	UserId   uint64 `json:"userId"`
	Nickname string `json:"nickName"`
}
type CustomClaims struct {
	UserInfo UserInfoStruct `json:"userInfo"`
	jwt.StandardClaims
}

// 一些常量
var (
	TokenExpired     error = errors.New("Token is expired")
	TokenNotValidYet error = errors.New("Token not active yet")
	TokenMalformed   error = errors.New("That's not even a token")
	TokenInvalid     error = errors.New("Couldn't handle this token")
)

func NoAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		zlog.Logger.WithGinContext(ctx).Info("建立请求")
		ctx.Next()
	}
}

// JWTAuth 中间件，检查token
func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.Request.Header.Get("Authorization")
		if tokenString == "" {
			zlog.Logger.WithContext(ctx).Warn("请求未携带token，无权限访问", zap.Any("data", map[string]interface{}{
				"url":    ctx.Request.URL,
				"params": ctx.Params,
			}))
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error_code": 1,
				"message":    "请求未携带token，无权限访问",
				"data":       map[string]interface{}{},
			})
			ctx.Abort()
			return
		}
		j := &JWT{
			[]byte(GetSignKey()),
		}

		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(tokenString)
		if err != nil {
			if err == TokenExpired {
				zlog.Logger.WithContext(ctx).Warn("token 过期", zap.Any("data", map[string]interface{}{
					"url":         ctx.Request.URL,
					"params":      ctx.Params,
					"tokenString": tokenString,
				}))
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"error_code": 1,
					"message":    "授权已过期",
					"data":       map[string]interface{}{},
				})
				ctx.Abort()
				return
			}
			zlog.Logger.WithContext(ctx).Error("token 错误", zap.Any("data", map[string]interface{}{
				"url":    ctx.Request.URL,
				"params": ctx.Params,
			}))
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error_code": 1,
				"message":    err.Error(),
				"data":       map[string]interface{}{},
			})
			ctx.Abort()
			return
		}

		// 继续交由下一个路由处理,并将解析出的信息传递下去
		ctx.Set("claims", claims)
		recoveryLoggerFunc(ctx)
		ctx.Next()
	}
}

// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

// 载荷，可以加一些自己需要的信息

// 获取signKey
func GetSignKey() string {
	return tconfig.C.GetString("jwt_sign")
}

// CreateToken 生成一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析Tokne
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	// 需要从tokenString移除 Bearer，jwt-go这个包为了避免冗余，不会帮我们处理
	re, _ := regexp.Compile(`(?i)Bearer `)
	tokenString = re.ReplaceAllString(tokenString, "")
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

// 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(2 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}

func recoveryLoggerFunc(ctx *gin.Context) {
	userInfo := ctx.MustGet("claims").(*CustomClaims).UserInfo
	zlog.Logger.NewContext(ctx, zap.Uint64("userId", userInfo.UserId), zap.String("nickname", userInfo.Nickname))
	zlog.Logger.WithContext(ctx).Info("建立请求")
}
