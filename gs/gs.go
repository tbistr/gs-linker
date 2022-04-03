package gs

type Linker struct {
	gSecret []byte
}

func New() *Linker {
	return &Linker{}
}
