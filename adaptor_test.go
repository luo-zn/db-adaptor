package db_adaptor

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var uri = "mongodb://localhost:27017"

func TestNewDbAdaptor(t *testing.T) {
	adap := NewDbAdaptor(&AdaptorOptions{Uri: uri})
	assert.Equal(t, adap.Opt.Uri, uri)
}

func TestNewDbAdaptor2MgoClient(t *testing.T) {
	t.Run("MgoClient", func(t *testing.T) {
		adap := NewDbAdaptor(&AdaptorOptions{Uri: uri})
		assert.Equal(t, adap.Opt.Uri, uri)
		assert.Equal(t, adap.Opt.DBType, "mongodb")
		uType := reflect.TypeOf(adap.DbC)
		assert.Implements(t,(*DBClient)(nil), adap.DbC,"DbAdaptor.Dbc does not implement DBClient!")
		assert.Equal(t,"*mongodb.MgoClient", uType.String(), "DbAdaptor.Dbc does not *mongodb.MgoClient!")
	})


}
