package token

import (
	"github.com/go-playground/validator/v10"
)

// GrantType is the type for OAuth2 param `grant_type`
type GrantType string

// oauth2 Authorization Grant: https://tools.ietf.org/html/rfc6749#section-1.3
const (
	// AUTHCODE oauth2 Authorization Code Grant
	AUTHCODE GrantType = "authorization_code"
	// IMPLICIT oauth2 Implicit Grant
	IMPLICIT GrantType = "implicit"
	// PASSWORD oauth2 Resource Owner Password Credentials Grant
	PASSWORD GrantType = "password"
	// CLIENT oauth2 Client Credentials Grant
	CLIENT GrantType = "client_credentials"
	// REFRESH oauth2 Refreshing an Access Token
	REFRESH GrantType = "refresh_token"
	// UPSCOPE is an custom grant for use unscope token acquire scope token
	UPSCOPE GrantType = "upgrade_scope"
	// WeChat is an custom grant for use unscope token acquire scope token
	WeChat GrantType = "wechat"
)

// Type token type
type Type string

// oauth2 Token Type: https://tools.ietf.org/html/rfc6749#section-7.1
const (
	// Bearer detail: https://tools.ietf.org/html/rfc6750
	Bearer Type = "bearer"
	// MAC detail: https://tools.ietf.org/html/rfc6749#ref-OAuth-HTTP-MAC
	MAC Type = "mac"
	// JWT detail:  https://tools.ietf.org/html/rfc7519
	JWT Type = "jwt"
)

// use a single instance of Validate, it caches struct info
var (
	validate = validator.New()
)

// Token is user's access resource token
type Token struct {
	AccessToken   string `bson:"access_token" json:"access_token"`               // 服务访问令牌
	RefreshToken  string `bson:"refresh_token" json:"refresh_token,omitempty"`   // 用于刷新访问令牌的凭证, 刷新过后, 原先令牌将会被删除
	Name          string `bson:"name" json:"name,omitempty"`                     // 独立颁发给SDK使用时, 命名token
	Description   string `bson:"description" json:"description,omitempty"`       // 独立颁发给SDK使用时, 令牌的描述信息, 方便定位与取消
	ApplicationID string `bson:"application_id" json:"application_id,omitempty"` // 用户应用ID, 如果凭证是颁发给应用的, 应用在删除时需要删除所有的令牌, 应用禁用时, 该应用令牌验证会不通过
	CreatedAt     int64  `bson:"create_at" json:"create_at,omitempty"`           // 凭证创建时间
	ExpiresIn     int64  `bson:"-" json:"ttl,omitempty"`                         // 凭证过期的时间
	ExpiresAt     int64  `bson:"expires_at" json:"expires_at,omitempty"`         // 还有多久过期

	ClientID  string    `bson:"client_id" json:"client_id,omitempty"`   // 客户端ID
	GrantType GrantType `bson:"grant_type" json:"grant_type,omitempty"` // 授权的类型
	Type      Type      `bson:"type" json:"type,omitempty"`             // 令牌的类型 类型包含: bearer/jwt  (默认为bearer)
	Scope     string    `bson:"scope" json:"scope,omitempty"`           // 令牌的作用范围: detail https://tools.ietf.org/html/rfc6749#section-3.3

	UserID    string `bson:"user_id" json:"user_id,omitempty"`       // 用户ID
	ProjectID string `bson:"project_id" json:"project_id,omitempty"` // 当前所在项目
	DomainID  string `bson:"domain_id" json:"domain_id,omitempty"`   // 用户所在的域的ID, 用户可以切换域(如果用户加入了多个域)
	ServiceID string `bson:"service_id" json:"service_id,omitempty"` // 服务ID, 如果凭证是颁发给内部服务使用时, 服务删除时,颁发给它的令牌需要删除, 服务禁用时, 令牌验证不通过
}
