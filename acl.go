package main

import (
	"reflect"
	"strings"
)

// as struct tags can continue multiple values, this gets the first
func GetFirstTagValue(tag reflect.StructTag, attr string) (value string) {
	tagValue := tag.Get(attr)
	if len(tagValue) == 0 {
		return
	}
	value = strings.Split(tagValue, ",")[0]
	return
}

// get field name based on struct tags
func GetFieldName(field reflect.StructField) (value string) {
	value = field.Tag.Get("access_field_name")
	if len(value) > 0 {
		return
	}
	value = GetFirstTagValue(field.Tag, "json")
	if len(value) > 0 {
		return
	}
	value = field.Name
	return
}

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
