package chikador

type options struct {
	recursive bool
	kind      watcherKind
}

type watcherKind func(chismis *Chismis)
type Option func(options *options)

// Recursive will make chikador scan through the directory and its subdirectories to add to the file watcher.
var Recursive = func(options *options) {
	options.recursive = true
}

// WithDedupe enables deduping of events. As operating systems tends to duplicate events around and give you
// "unnecessary noise", for most cases, it is recommended to enable this unless you need those "unnecessary noise".
var WithDedupe = func(options *options) {
	options.kind = withDedupe
}

// WithoutDedupe explicitly states that the Chismis instance shouldn't dedupe events, this is the default behavior, and
// this is only added for completeness and when you want to explicitly indicate the behavior.
var WithoutDedupe = func(options *options) {
	options.kind = withoutDedupe
}
