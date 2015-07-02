package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type User struct {
	DisplayName string `access:"admin:*, owner: read update, *: read" json:"display_name"`
	Password    string `access:"admin:*, owner: update" json:"password"`
	ACL         *ACL   `access:"admin:*, owner: update, *: read" access_field_name:"acl"`
}

func TestACL(t *testing.T) {
	user := User{}
	user.ACL = NewACL(user)
	assert.NotNil(t, user.ACL.global)
	// check global permissions:
	assert.Equal(t, false, user.ACL.HasAccess([]string{"anon"}, "update"))
	assert.Equal(t, true, user.ACL.HasAccess([]string{"anon"}, "read"))
	assert.Equal(t, true, user.ACL.HasAccess([]string{"admin"}, "delete"))
	assert.Equal(t, true, user.ACL.HasAccess([]string{"owner"}, "update"))
	assert.Equal(t, false, user.ACL.HasAccess([]string{"owner"}, "delete"))

	// check per field access:
	assert.Equal(t, true, user.ACL.HasAccessToFields([]string{"owner"}, "read", []string{"display_name"}))
	assert.Equal(t, false, user.ACL.HasAccessToFields([]string{"owner"}, "read", []string{"password"}))
	assert.Equal(t, true, user.ACL.HasAccessToFields([]string{"admin"}, "read", []string{"password"}))
	assert.Equal(t, true, user.ACL.HasAccessToFields([]string{"anon"}, "read", []string{"display_name"}))
	assert.Equal(t, false, user.ACL.HasAccessToFields([]string{"anon"}, "read", []string{"display_name", "password"}))

	allowedAdminActions := user.ACL.AllowedActions([]string{"admin", "owner"})
	assert.Equal(t, len(allowedAdminActions), 3)
	allowedAnonActions := user.ACL.AllowedActions([]string{"anon"})
	assert.Equal(t, len(allowedAnonActions), 1)

}
