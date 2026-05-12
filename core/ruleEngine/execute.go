package ruleengine

type FileOptions map[string]Options

type Options struct {
	MovePaths  []string
	CopyPaths  []string
	DeleteFlag bool
}

type operations interface {
	Copy()
	Move()
	Delete()
}
