module github.com/godaddy-x/wallet-api-go

go 1.26

require (
	github.com/godaddy-x/freego v1.1.39
	github.com/godaddy-x/wallet-adapter v1.0.13
	github.com/mailru/easyjson v0.9.2
)

require (
	filippo.io/mldsa v0.0.0-20260215214346-43d0283efc3e // indirect
	github.com/godaddy-x/eccrypto v1.1.20 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/lxzan/gws v1.9.1 // indirect
	github.com/valyala/fastjson v1.6.3 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.49.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

//replace github.com/godaddy-x/wallet-adapter => ../wallet-adapter
//
//replace github.com/godaddy-x/freego => ../freego
