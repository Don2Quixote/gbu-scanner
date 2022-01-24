// Package app provides function Run which is like main, but returns error that can
// be handled in real main function. This approach allows to call os.Exit or log.Fatal(...)
// once in main instead of calling it on each error.
package app
