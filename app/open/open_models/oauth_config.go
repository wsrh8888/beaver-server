package open_models

// OAuthClientConfig OAuth 客户端配置（区分不同平台）
type OAuthClientConfig struct {
	H5      H5OAuthConfig      `json:"h5"`
	Desktop DesktopOAuthConfig `json:"desktop"`
	Mobile  MobileOAuthConfig  `json:"mobile"`
}

// H5OAuthConfig H5 应用 OAuth 配置
type H5OAuthConfig struct {
	Enabled      bool     `json:"enabled"`
	RedirectURIs []string `json:"redirect_uris"`
	JsSdkDomains []string `json:"js_sdk_domains"` // JS-SDK 安全域名
}

// DesktopOAuthConfig 桌面端 OAuth 配置
type DesktopOAuthConfig struct {
	Enabled      bool   `json:"enabled"`
	CustomScheme string `json:"custom_scheme"` // 如: beaver://oauth/callback
}

// MobileOAuthConfig 移动端 OAuth 配置
type MobileOAuthConfig struct {
	Enabled            bool   `json:"enabled"`
	IOSBundleID        string `json:"ios_bundle_id"`
	AndroidPackageName string `json:"android_package_name"`
	UniversalLink      string `json:"universal_link"` // iOS Universal Link
	CustomScheme       string `json:"custom_scheme"`  // 如: beaver://oauth/callback
}
