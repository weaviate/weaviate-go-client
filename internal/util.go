package internal

// Last returns the Last element of the vararg or a zero value for type T,
// if the len of vararg was 0. The second return is false in the latter case.
//
// This is useful for extracting the optional argument from a variadic.
//
// Usage:
//
//	type SayOption struct { Volume int }
//
//	func say(hello string, option ...SayOption) {
//		opt := SayOption{ Volume: 60 } // default options
//		if option, ok := internal.Last(options...); ok {
//			opt = option
//		}
//		sayWithVolume(opt)
//	}
//
// Ignore the ok check if the value should be set unconditionally:
//
//	func say(hello string, option ...SayOption) {
//		opt, _ := internal.Last(option...)
//		opt.Volume *= 2 // Always say twice as loud
//		sayWithVolume(opt)
//	}
func Last[T any](s ...T) (T, bool) {
	if len(s) == 0 {
		return *new(T), false
	}
	return s[len(s)-1], true
}
