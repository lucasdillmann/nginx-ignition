package valuerange

type ValueRange struct {
	Min int
	Max int
}

func New(min, max int) *ValueRange {
	return &ValueRange{Min: min, Max: max}
}

func (vr *ValueRange) Contains(value int) bool {
	return value >= vr.Min && value <= vr.Max
}
