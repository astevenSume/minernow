package dao

import "errors"

const (
	InsertMulCount = 100
	DefaultPerPage = 100
)

var (
	ErrParam = errors.New("param error")
	ErrSql   = errors.New("db sql error")
)

type PageInfo struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type RoleInfo struct {
	Id          uint64   `json:"id"`
	Name        string   `json:"name"`
	Desc        string   `json:"desc"`
	Permissions []string `json:"permissions"`
	Ctime       int64    `json:"ctime"`
	Utime       int64    `json:"utime"`
}

type RoleMember struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Status    int8   `json:"status"`
	GrantedBy string `json:"granted_by"`
	GrantedAt int64  `json:"granted_at"`
}

type RoleBase struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}
type AdminMember struct {
	Id     uint64     `json:"id"`
	Name   string     `json:"name"`
	Email  string     `json:"email"`
	Status int8       `json:"status"`
	RBase  []RoleBase `json:"roles"`
	//Permissions []string `json:"permissions"`
	WhitelistIps string `json:"whitelist_ips"`
	CTime        int64  `json:"ctime"`
	UTime        int64  `json:"utime"`
	TimeLogin    int64  `json:"ltime"`
	DTime        int64  `json:"dtime"`
	IsBind       bool   `json:"is_bind"`
}

type LoginInfo struct {
	Id           uint64     `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	Status       int8       `json:"status"`
	RBase        []RoleBase `json:"roles"`
	Permissions  []string   `json:"permissions"`
	WhitelistIps string     `json:"whitelist_ips"`
	CTime        int64      `json:"ctime"`
	UTime        int64      `json:"utime"`
	TimeLogin    int64      `json:"ltime"`
	IsBind       bool       `json:"is_bind"`
	SecretId     string     `json:"secret_id"`
	QrCode       string     `json:"qr_code"`
}

//返佣配置信息
type Commissioncfg struct {
	Min        uint64 `json:"min"`
	Max        uint64 `json:"max"`
	Commission int32  `json:"commission"`
}
type Commissioncfgs []Commissioncfg

func (p Commissioncfgs) Len() int { return len(p) }
func (p Commissioncfgs) Less(i, j int) bool {
	return p[i].Min < p[j].Min
}
func (p Commissioncfgs) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

//月分红档位配置信息
type MonthDividendCfg struct {
	Min           int64 `json:"min"`
	Max           int64 `json:"max"`
	AgentLv       int32 `json:"agent_lv"`
	Position      int32 `json:"position"`
	DividendRatio int32 `json:"dividend_ratio"`
	ActivityNum   int32 `json:"activity_num"`
}
type MonthDividendCfgs []MonthDividendCfg

func (p MonthDividendCfgs) Len() int { return len(p) }
func (p MonthDividendCfgs) Less(i, j int) bool {
	return p[i].DividendRatio < p[j].DividendRatio
}
func (p MonthDividendCfgs) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
