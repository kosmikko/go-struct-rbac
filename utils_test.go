package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUtils(t *testing.T) {
	u1 := &User{DisplayName: "foobar", Password: "abc", Settings: map[string]string{"foo": "bar"}}
	u2 := &User{DisplayName: "foobarf", Password: "abcd", Settings: map[string]string{"foo": "bar"}}
	u3 := &User{DisplayName: "foobar", Password: "abc", Settings: map[string]string{"foo": "bar"}}
	u4 := &User{DisplayName: "foobar", Password: "abcd", Settings: map[string]string{"foo": "barf"}}
	err, changed := GetChangedFields(u1, u2)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(changed))
	err, changed = GetChangedFields(u1, u3)
	assert.Equal(t, 0, len(changed))
	err, changed = GetChangedFields(u1, u4)
	assert.Equal(t, 2, len(changed))
}
