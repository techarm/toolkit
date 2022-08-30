module github.com/techarm/toolkit

go 1.18

retract (
	// Published v1 too early
	[v1.0.0, v2.0.2]
)

require (
	github.com/go-stack/stack v1.8.1
	github.com/mattn/go-isatty v0.0.16
)

require (
	github.com/elliotchance/orderedmap v1.4.0 // indirect
	golang.org/x/exp v0.0.0-20220321173239-a90fa8a75705 // indirect
	golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab // indirect
)
