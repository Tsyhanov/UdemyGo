package forms

type errors map[string][]string

//add error msg for a given form field
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

//return first error msg
func (e errors) Get(field string) string {
	es := e[field]

	if len(es) == 0 {
		return ""
	}

	return es[0]
}
