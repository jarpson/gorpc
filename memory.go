package gorpc


// TODO Get buf from memory 

// resize buffer
// input: s: src buffer, newlen: new buffer min length(abrove len s)
// return new buffer 
func ResizeBuf(s []byte, newlen int) []byte {
	d := make([]byte, newlen)
	copy(d, s)
	return d
}

func GetBufN(n int) []byte {
	return make([]byte, n)
}
