package db

type DB struct {
	// List of rules unique by (target, source).
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	// unique identifier
	ID int `yaml:"id"`
	// Target path/link to document.
	Target string `yaml:"target"`
	// Path to source, can be glob pattern.
	Source  []string `yaml:"source"`
	Entries []Entry  `yaml:"entries"`
}

type Entry struct {
	Path     string `yaml:"path"`
	Checksum string `yaml:"checksum"`
}

type CompareResult string

const (
	CompareNew     CompareResult = "new"
	CompareDeleted CompareResult = "deleted"
	CompareChanged CompareResult = "changed"
)

type EntryCompareResult struct {
	Path        string
	Result      CompareResult
	OldChecksum string
	NewChecksum string
}
