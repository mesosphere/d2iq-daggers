package daggers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	defaultOutputDir = ".daggers"
	defaultFileName  = "dagger.log"
)

var (
	_ io.Writer = new(Logger)
	_ io.Closer = new(Logger)
)

// Logger is a simple interface wrapper around the standard io.Writer interface.
type Logger struct {
	w io.Writer
}

// NewLogger returns simple logger that writes logs to a log file under default output directory(.daggers). If verbose
// is true, logs will be shown in os.Stdout at the sametime.
func NewLogger(verbose bool) (Logger, error) {
	logger := Logger{}

	wd, err := os.Getwd()
	if err != nil {
		return logger, err
	}

	if err := os.MkdirAll(filepath.Join(wd, defaultOutputDir), os.ModePerm); err != nil {
		return logger, err
	}

	filePath := filepath.Join(wd, defaultOutputDir, defaultFileName)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return logger, err
	}

	if verbose {
		logger.w = io.MultiWriter(file, os.Stdout)
	} else {
		logger.w = file
	}

	return logger, nil
}

// Write writes the given bytes to the underlying writer.
func (l Logger) Write(p []byte) (n int, err error) {
	return fmt.Fprint(l.w, string(p))
}

// Close closes the underlying writer.
func (l Logger) Close() error {
	if closer, ok := l.w.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}
