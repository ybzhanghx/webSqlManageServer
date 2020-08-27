module bailun.com/CT4_quote_server/WebManageSvr

go 1.14

replace (
	bailun.com/CT4_quote_server/WebManageSvr => ./
	bailun.com/CT4_quote_server/common => ../common
	bailun.com/CT4_quote_server/lib => ../lib
	bailun.com/CT4_quote_server/protocol => ../protocol
)

require (
	bailun.com/CT4_quote_server/common v0.0.0-00010101000000-000000000000
	bailun.com/CT4_quote_server/lib v0.0.0-00010101000000-000000000000
	bailun.com/CT4_quote_server/protocol v0.0.0-00010101000000-000000000000
	github.com/BurntSushi/toml v0.3.1
	github.com/astaxie/beego v1.12.1
	github.com/jmoiron/sqlx v1.2.0
)
