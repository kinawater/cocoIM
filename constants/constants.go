package constants

type EnvType string

const (
	EnvDebug EnvType = "debug"
	EnvStg   EnvType = "stg"
	EnvPre   EnvType = "pre"
	EnvPrd   EnvType = "prd"
)

// 便于访问
var Env = struct {
	Debug EnvType
	Stg   EnvType
	Pre   EnvType
	Prd   EnvType
}{
	Debug: EnvDebug,
	Stg:   EnvStg,
	Pre:   EnvPre,
	Prd:   EnvPrd,
}
