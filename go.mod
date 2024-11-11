module github.com/jeffque/teecp

go 1.22.5

replace github.com/jeffque/teecp/teecp_cli => ./teecp_cli

replace github.com/jeffque/teecp/teecp_server => ./teecp_server

replace github.com/jeffque/teecp/teecp_client => ./teecp_client

replace github.com/jeffque/teecp/teecp => ./teecp

require (
	github.com/jeffque/teecp/teecp_client v0.0.0-00010101000000-000000000000
	github.com/jeffque/teecp/teecp_server v0.0.0-00010101000000-000000000000
)

require github.com/jeffque/teecp/teecp v0.0.0-00010101000000-000000000000 // indirect
