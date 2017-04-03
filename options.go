package binder

type Options struct {
	// CallStackSize is call stack size
	CallStackSize int
	// RegistrySize is data stack size
	RegistrySize int
	// SkipOpenLibs controls whether or not libraries are opened by default
	SkipOpenLibs bool
	// IncludeGoStackTrace tells whether a Go stacktrace should be included in a Lua stacktrace when panics occur.
	IncludeGoStackTrace bool
}
