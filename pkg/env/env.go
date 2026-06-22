package envutil

type Env string

const (
	EnvLocal Env = "local"
	EnvDev   Env = "dev"
	EnvProd  Env = "prod"
)

func (e Env) IsLocal() bool { return e == EnvLocal }
func (e Env) IsDev() bool   { return e == EnvDev }
func (e Env) IsProd() bool  { return e == EnvProd }
