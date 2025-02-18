package types

// ContextProfile represents a context profile
type ContextProfile struct {
	Current   bool
	Name      string
	Cluster   string
	User      string
	Server    string
	Namespace string
	Emoji     string
}
