package manager

//Manager represents the manager that collects and processes node data
type Manager struct {
}

//NewManager creates a new Manager
func NewManager() *Manager {
	return new(Manager)
}
