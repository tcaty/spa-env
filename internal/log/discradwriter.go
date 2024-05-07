package log

type discardWriter struct {
}

func newDiscardWriter() *discardWriter {
	return &discardWriter{}
}

func (dw *discardWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}
