module github.com/techarm/toolkit

go 1.18

retract (
	// Published v1 too early
	[v1.0.0, v2.0.1]
)

require (
	github.com/go-stack/stack v1.8.1
	github.com/mattn/go-isatty v0.0.16
)

require golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab // indirect
