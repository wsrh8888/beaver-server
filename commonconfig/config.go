// Package commonconfig 出站立收费模块的跨环境域名与凭据注册表（仓级共享）。
//
// 对应 cloud-server/commonconfig：yaml 是部署事实，本包是它在进程里的唯一视图。
// 各服务在装配依赖前经 remote/<模块>.Init 把域名与凭据灌进来，remote/<模块> 的调用
// 封装只读这里，不在 80 个调用行里重复写地址。
//
// 只存原值：不去空白、不削末尾斜杠，拼地址的规则由调用方自己保证（多一个斜杠就是 404）。
//
// 包名 commonconfig 是为了和「各 app 内部的 internal/config」区分开——后者是服务私有配置，
// 本包是仓级共享的跨环境注册表，二者都叫 config 容易混。
package commonconfig

// Endpoint 出站立模块（被调方）S2S 注册表项：域名 + client_credentials。
type Endpoint struct {
	BaseURL      string `json:",optional"`
	ClientID     string `json:",optional"`
	ClientSecret string `json:",optional"`
}

// AgentEndpoint 智能体（beaver-agent，py）S2S 注册项：静态 X-Internal-Token，无 OAuth 换票。
type AgentEndpoint struct {
	BaseURL       string `json:",optional"`
	InternalToken string `json:",optional"`
}
