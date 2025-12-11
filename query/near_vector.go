package query

type NearVectorRequest struct {
	CommonOptions
	Target              NearVectorTarget
	Distance, Certainty *float32
}

type NearVectorOption interface {
	Apply(*NearVectorRequest)
}

func (l LimitOption) Apply(r *NearVectorRequest) {
	r.CommonOptions.Limit = (*int)(&l)
}

func (l OffsetOption) Apply(r *NearVectorRequest) {
	r.CommonOptions.Offset = (*int)(&l)
}

func (l AutoLimitOption) Apply(r *NearVectorRequest) {
	r.CommonOptions.AutoLimit = (*int)(&l)
}

// DistanceOption sets the `distance` parameter.
type DistanceOption float32

var _ NearVectorOption = (*DistanceOption)(nil)

func WithDistance(l float32) DistanceOption {
	return DistanceOption(l)
}

func (d DistanceOption) Apply(r *NearVectorRequest) {
	r.Distance = (*float32)(&d)
}

// CertaintyOption sets the `certainty` parameter.
type CertaintyOption float32

var _ NearVectorOption = (*CertaintyOption)(nil)

func WithCertainty(l float32) CertaintyOption {
	return CertaintyOption(l)
}

func (d CertaintyOption) Apply(r *NearVectorRequest) {
	r.Certainty = (*float32)(&d)
}

type NearVectorFunc func(NearVectorTarget, ...NearVectorOption) (any, error)

func nearVector(NearVectorTarget, ...NearVectorOption) (any, error) {
	return nil, nil
}

func (nv NearVectorFunc) GroupBy(NearVectorTarget /* groupBy */, any, ...NearVectorOption) (any, error) {
	return nil, nil
}

type NearVectorOptions []NearVectorOption

type NearVectorTarget interface {
	ToProto()
}
