package integrations_test

import (
	"github.com/magiconair/properties/assert"
	"sso/tests"
	"testing"
)

func TestPing(t *testing.T) {
	w := tests.GetJson("/ping", nil, "")
	assert.Equal(t, `{"success":true}`, w.Body.String())
}

func TestNotFound(t *testing.T) {
	w := tests.GetJson("/not_found", nil, "")
	assert.Equal(t, `{"code":404,"message":"Page not found"}`, w.Body.String())
}
