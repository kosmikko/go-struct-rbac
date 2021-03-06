package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type User struct {
	DisplayName string            `access:"admin:*, owner: read update, *: read create" json:"display_name"`
	Password    string            `access:"admin:*, owner: update, *: create" json:"password"`
	Settings    map[string]string `access:"admin:*, owner: update, *: create" json:"settings"`
	ACL         *ACL              `access:"admin:*, owner: update, *: read create" access_field_name:"acl"`
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
	assert.Equal(t, len(allowedAdminActions), 4)
	allowedAnonActions := user.ACL.AllowedActions([]string{"anon"})
	assert.Equal(t, len(allowedAnonActions), 2)

	allowedAnonReadFields := user.ACL.AllowedFields([]string{"anon"}, "read")
	assert.Equal(t, len(allowedAnonReadFields), 1)
	allowedAdminReadFields := user.ACL.AllowedFields([]string{"admin"}, "read")
	assert.Equal(t, len(allowedAdminReadFields), 3)
	assert.Equal(t, stringInSlice("password", allowedAdminReadFields), true)
	allowedAdminDeleteFields := user.ACL.AllowedFields([]string{"admin"}, "delete")
	assert.Equal(t, len(allowedAdminDeleteFields), 3)

	u1 := &User{DisplayName: "foobar"}
	u2 := &User{DisplayName: "foobarx"}
	anonActor := &Actor{Roles: []string{"anon"}}
	adminActor := &Actor{Roles: []string{"anon", "admin"}}
	err := user.ACL.CheckChangeAccess(anonActor, "update", u1, u2)
	assert.NotNil(t, err)
	err = user.ACL.CheckChangeAccess(adminActor, "update", u1, u2)
	assert.Nil(t, err)
	err = user.ACL.CheckAccess(anonActor, "create", u1)
	assert.Nil(t, err)

}
