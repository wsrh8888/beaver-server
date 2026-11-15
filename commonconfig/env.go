package commonconfig

// 域名清单，由各服务装配依赖前经 remote/<模块>.Init 灌入；该服务用不到的对端保持空。
//
// 一个进程只跑一个服务，不存在多份 Init 交叉覆盖。Init 直接写对应字段（baseURL），
// 这里只是它在进程里的唯一视图。
var (
	CircleUrl string // 圈子服务（beaver-server-social/app/circle，网关 /internal/Circle.*）
	DocUrl    string // 在线文档服务（beaver-doc 独立仓，网关 /internal/Doc.*）
	MomentUrl string // 圈子动态服务（beaver-server-social/app/moment，网关 /internal/Moment.*）
	AiUrl     string // AI 服务（beaver-ai 独立仓，网关 /internal/Ai.*）
	OpenUrl   string // 开放接口服务（beaver-open 独立仓，网关 /internal/Open.*）
	AgentUrl  string // 智能体服务（beaver-agent，py，网关 /internal/agent.*，静态 X-Internal-Token）
)
