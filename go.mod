module github.com/techarm/toolkit/v2

go 1.18

retract (
	// Published v2 too early
	[v0.0.1, v2.0.2]
)

require (
	github.com/go-stack/stack v1.8.1
	github.com/mattn/go-isatty v0.0.16
)

require golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab // indirect
