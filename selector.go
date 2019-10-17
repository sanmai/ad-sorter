package bannersorter

type selector struct {
	inputs   []Input
	selector func(in Input) bool
}

func (s selector) earliestExpiringInput() *Input {
	var result *Input

	for _, current := range s.inputs {
		input := current

		if !s.selector(input) {
			continue
		}

		if result == nil {
			result = &input
			continue
		}

		if input.EndTime().Before((*result).EndTime()) {
			result = &input
		}
	}

	return result
}
