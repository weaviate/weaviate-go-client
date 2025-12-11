package query

type nearVectorRequest struct {
	CommonOptions
	Target              NearVectorTarget
	Distance, Certainty *float32
}

type nearVectorOption interface {
	Apply(*nearVectorRequest)
}

func (l LimitOption) Apply(r *nearVectorRequest) {
	r.CommonOptions.Limit = (*int)(&l)
}

func (l OffsetOption) Apply(r *nearVectorRequest) {
	r.CommonOptions.Offset = (*int)(&l)
}

func (l AutoLimitOption) Apply(r *nearVectorRequest) {
	r.CommonOptions.AutoLimit = (*int)(&l)
}

// DistanceOption sets the `distance` parameter.
type DistanceOption float32

var _ nearVectorOption = (*DistanceOption)(nil)

func WithDistance(l float32) DistanceOption {
	return DistanceOption(l)
}

func (d DistanceOption) Apply(r *nearVectorRequest) {
	r.Distance = (*float32)(&d)
}

// CertaintyOption sets the `certainty` parameter.
type CertaintyOption float32

var _ nearVectorOption = (*CertaintyOption)(nil)

func WithCertainty(l float32) CertaintyOption {
	return CertaintyOption(l)
}

func (d CertaintyOption) Apply(r *nearVectorRequest) {
	r.Certainty = (*float32)(&d)
}

type NearVectorFunc func(NearVectorTarget, ...nearVectorOption) (any, error)

func nearVector(NearVectorTarget, ...nearVectorOption) (any, error) {
	return nil, nil
}

func (nv NearVectorFunc) GroupBy(NearVectorTarget /* groupBy */, any, ...nearVectorOption) (any, error) {
	return nil, nil
}

type NearVectorOptions []nearVectorOption

type NearVectorTarget interface {
	ToProto()
}
