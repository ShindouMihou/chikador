package chikador

type options struct {
	recursive bool
	kind      watcherKind
}

type watcherKind func(chismis *Chismis)
type Option func(options *options)

var Recursive = func(options *options) {
	options.recursive = true
}

var WithDedupe = func(options *options) {
	options.kind = withDedupe
}

var WithoutDedupe = func(options *options) {
	options.kind = withoutDedupe
}
