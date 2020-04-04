package server

import "strings"

type LoggerPrefix interface {
	SetupLoggerPrefix() (prefix strings.Builder)
}

type Setup interface {
	Setup()
}

