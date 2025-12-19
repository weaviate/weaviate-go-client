package internal

// Optional returns v if it is not nil and a zero value of T otherwise.
//
// Example:
//
//	type SayOptions struct { Word string }
//
//	func say(ctx context.Context, opt *SayOptions) {
//		opt = internal.Optional(opt)
//		print(opt.Word) // opt can be safely dereferenced
//	}
func Optional[T any](v *T) *T {
	if v != nil {
		return v
	}
	return new(T)
}
