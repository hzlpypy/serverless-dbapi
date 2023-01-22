package valueobject

type Params struct {
	QueryParams map[string][]string
	Body        []byte
}

type Cursor struct {
	Continue string
	Limit    int
}
