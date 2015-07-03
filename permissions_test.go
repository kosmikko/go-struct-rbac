package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPermissions(t *testing.T) {
	permissions := NewPermissions("admin:*, owner: read update, *: read")
	adminPerms := permissions.RolesAccess["admin"]
	assert.Equal(t, 1, len(adminPerms), "Admin permissions had incorrect length")
	assert.Equal(t, "*", adminPerms[0], "Admin permissions was incorrect")
	assert.Equal(t, permissions.HasAccess([]string{"admin"}, "delete"), true, "admin should have access")
	assert.Equal(t, permissions.HasAccess([]string{"anon"}, "delete"), false, "anon should not have access")
	assert.Equal(t, permissions.HasAccess([]string{"anon"}, "read"), true, "anon should have read access")
	assert.Equal(t, permissions.HasAccess([]string{"owner"}, "update"), true, "owner should have update access")
	assert.Equal(t, permissions.HasAccess([]string{"owner"}, "delete"), false, "owner should not have delete access")
	allowedOwnerActions := permissions.AllowedActions([]string{"owner"})
	assert.Equal(t, 2, len(allowedOwnerActions))
	allowedAnonActions := permissions.AllowedActions([]string{"anon"})
	assert.Equal(t, 1, len(allowedAnonActions))
	allowedAdminActions := permissions.AllowedActions([]string{"admin"})
	assert.Equal(t, 2, len(allowedAdminActions))
	assert.Equal(t, true, stringInSlice("*", allowedAdminActions))
}
