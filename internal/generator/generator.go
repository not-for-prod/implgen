package generator

// GenerateCommand holds configuration for generating a Go implementation
// of an interface, including destination, naming, and output structure.
type GenerateCommand struct {
	// dst is the base directory where generated files will be written.
	dst string

	// interfaceName is the name of the interface to generate.
	// If empty, all interfaces in the package will be processed.
	interfaceName string

	// implementationName is the name of the struct that implements the interface.
	implementationName string

	// implementationPackageName overrides the default package name for the generated code.
	// If empty, the package name is derived from the interface name.
	implementationPackageName string

	// singleFile determines whether all methods should be generated into a single file.
	// If false, each method will be written into its own file.
	singleFile bool

	// tracerName is the name used in otel.Tracer(tracerName).Start(...) for tracing spans.
	tracerName string
}

// NewGenerateCommand creates a new GenerateCommand with the given parameters.
// This struct is typically passed into a code generator to drive its behavior.
func NewGenerateCommand(
	dst string,
	interfaceName string,
	implementationName string,
	implementationPackageName string,
	tracerName string,
	singleFile bool,
) *GenerateCommand {
	return &GenerateCommand{
		dst:                       dst,
		interfaceName:             interfaceName,
		implementationName:        implementationName,
		implementationPackageName: implementationPackageName,
		tracerName:                tracerName,
		singleFile:                singleFile,
	}
}
