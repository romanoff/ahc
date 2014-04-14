package component

type Preprocessor interface {
	GetCss([]byte) []byte
}
