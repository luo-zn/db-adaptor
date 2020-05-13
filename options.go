/* Create By Jenner.luo */
package db_adaptor

import "time"

type AdaptorOptions struct {
	DBType  string        // DB type
	Uri     string        //connect string
	Timeout time.Duration // connection timeout
}

//NewAdaptorOptions creates a new AdaptorOptions instance.
func NewAdaptorOptions() *AdaptorOptions {
	return new(AdaptorOptions)
}

//MergeAdaptorOptions combines the given *AdaptorOptions into a single *AdaptorOptions in a last one wins fashion.
// The specified options are merged with the existing options on the collection, with the specified options taking
// precedence.
func MergeAdaptorOptions(opts ...*AdaptorOptions) *AdaptorOptions {
	c := NewAdaptorOptions()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.DBType != "" {
			c.DBType = opt.DBType
		}
		if opt.Uri != "" {
			c.Uri = opt.Uri
		}
		if opt.Timeout != 0 {
			c.Timeout = opt.Timeout
		}
	}
	return c
}
