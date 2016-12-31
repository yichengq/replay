package replay

type compileType int

const (
	compileToJS compileType = iota + 1
	compileToML
)

func compile(code []byte, typ compileType) []byte {
	return nil
}
