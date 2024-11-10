package db

func NewRule(id int, target, source string) Rule {
	return Rule{
		ID:     id,
		Target: target,
		Source: []string{source},
	}
}

// Compare entries.
func (r Rule) Compare(entries []Entry) []EntryCompareResult {
	ents := make(map[string]string, len(entries))
	for _, e := range entries {
		ents[e.Path] = e.Checksum
	}

	var result []EntryCompareResult
	for _, e := range r.Entries {
		cs, ok := ents[e.Path]
		switch {
		case !ok:
			result = append(result, EntryCompareResult{
				Path:   e.Path,
				Result: CompareDeleted,
			})
		case e.Checksum != cs:
			result = append(result, EntryCompareResult{
				Path:        e.Path,
				Result:      CompareChanged,
				OldChecksum: e.Checksum,
				NewChecksum: cs,
			})
		}
		delete(ents, e.Path)
	}

	for p := range ents {
		result = append(result, EntryCompareResult{
			Path:   p,
			Result: CompareNew,
		})
	}

	return result
}
