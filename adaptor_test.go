/* Create By Jenner.luo */
package db_adaptor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var uri = "mongodb://localhost:27017"

func TestNewDbAdaptor(t *testing.T) {
	adap := NewDbAdaptor(&AdaptorOptions{Uri: uri})
	assert.Equal(t, adap.Opt.Uri, uri)
}

func TestNewDbAdaptor2MgoClient(t *testing.T) {
	adap := NewDbAdaptor(&AdaptorOptions{Uri: uri})
	assert.Equal(t, adap.Opt.Uri, uri)
	assert.Equal(t, adap.Opt.DBType, "mongodb")
	client := adap.DbC
	t.Log(client)
	//uType := reflect.TypeOf(adap.DbC)
	//for i := 0; i < uType.NumField(); i++ {
	//	t.Log(uType.Field(i))
	//}
	t.Log(adap)
}
