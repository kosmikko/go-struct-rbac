package main

import (
	"fmt"
	"reflect"
)

// ACL defines struct's permissions globally (e.g. can read struct at all)
// and per field (e.g. can read field's value)
type ACL struct {
	global *Permissions
	// each field in a struct may define its permissions
	fields map[string]*Permissions
}

// check global access
func (acl *ACL) HasAccess(roles []string, action string) (isGranted bool) {
	return acl.global.HasAccess(roles, action)
}

// check global & per field access
func (acl *ACL) HasAccessToFields(roles []string, action string, fields []string) (isGranted bool) {
	globalAccess := acl.HasAccess(roles, action)
	if !globalAccess {
		return false
	}
	for _, field := range fields {
		fieldACL, fieldACLDefined := acl.fields[field]
		if fieldACLDefined {
			hasAccess := fieldACL.HasAccess(roles, action)
			if !hasAccess {
				return false
			}
		}
	}
	return true
}

// return allowed global actions for roles
func (acl *ACL) AllowedActions(roles []string) (allowedActions []string) {
	return acl.global.AllowedActions(roles)
}

// return a slice of struct fields given roles have access to do given action
func (acl *ACL) AllowedFields(roles []string, action string) (allowedFields []string) {
	for f, permissions := range acl.fields {
		if permissions.HasAccess(roles, action) {
			allowedFields = append(allowedFields, f)
		}
	}
	return
}

func (acl *ACL) CheckChangeAccess(actor Actor, action string, oldStruct, newStruct interface{}) (err error) {
	err, changes := GetChangedFields(oldStruct, newStruct)
	if err != nil {
		return
	}
	if len(changes) == 0 {
		return
	}
	roles := actor.Roles
	hasAccess := acl.HasAccessToFields(roles, action, changes)
	if !hasAccess {
		err = fmt.Errorf("No access to update fields %v", changes)
	}
	return
}

// read struct s tags & parse its permissions
func NewACL(s interface{}) (acl *ACL) {
	acl = &ACL{fields: make(map[string]*Permissions)}

	tagType := reflect.TypeOf(s)
	for i := 0; i < tagType.NumField(); i++ {
		field := tagType.Field(i)
		access := field.Tag.Get("access")
		fieldName := GetFieldName(field)
		if fieldName == "acl" {
			acl.global = NewPermissions(access)
		} else {
			acl.fields[fieldName] = NewPermissions(access)
		}
	}

	return
}
