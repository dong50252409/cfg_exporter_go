package interfaces

type IRender interface {
	ExportDir() string
	Execute() error
}
