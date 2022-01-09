package logger

type nop struct{}

var _ Logger = nop{}

// NewNop returns nop (No Operation) logger, which does nothing on calls.
func NewNop() Logger {
	return nop{}
}

func (nop) Debug(args ...interface{})                 {}
func (nop) Debugf(format string, args ...interface{}) {}
func (nop) Info(args ...interface{})                  {}
func (nop) Infof(format string, args ...interface{})  {}
func (nop) Warn(args ...interface{})                  {}
func (nop) Warnf(format string, args ...interface{})  {}
func (nop) Error(args ...interface{})                 {}
func (nop) Errorf(format string, args ...interface{}) {}
func (nop) Fatal(args ...interface{})                 {}
func (nop) Fatalf(format string, args ...interface{}) {}
