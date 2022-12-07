package oauth

type CachePrefix string

const (
	CachePrefixAuthorizationCode CachePrefix = "oauth:auth-code:"
	// CachePrefixOauthClient should be followed with clientID.
	CachePrefixOauthClient = "oauth:client:"
)

type Scope string

const (
	ScopeUserinfoRead Scope = "ur"
)

var Scopes = []Scope{
	ScopeUserinfoRead,
}

type AuthorizationCode struct {
	Code     string  `json:"code"`
	Scopes   []Scope `json:"scopes"`
	ClientID string  `json:"client_id"`
}
