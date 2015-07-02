package main

import (
	"strings"
)

/**
 * Permissions defines allowed actions per role,
 * e.g. [admin][update, delete]
 *
 * This may be used to define global permissions to an object
 * or per field permissions.
 */
type Permissions struct {
	RolesAccess map[string][]string
}

func isAllowed(allowedActions []string, action string) bool {
	for _, a := range allowedActions {
		if action == a || a == "*" {
			return true
		}
	}
	return false
}

// Check if given roles have access to do the given action
func (p *Permissions) HasAccess(roles []string, action string) (isGranted bool) {
	roles = append(roles, "*") // check everyone aka '*' role
	for _, role := range roles {
		roleAccess, roleDefined := p.RolesAccess[role]
		if roleDefined && isAllowed(roleAccess, action) {
			return true
		}
	}

	return false
}

// Return a list of actions given roles have access to
func (p *Permissions) AllowedActions(roles []string) (allowedActions []string) {
	roles = append(roles, "*")
	// keep track of actions already in slice to avoid duplicates
	already := make(map[string]bool)

	for _, role := range roles {
		roleAccess, roleDefined := p.RolesAccess[role]
		if roleDefined {
			for _, action := range roleAccess {
				if !already[action] {
					allowedActions = append(allowedActions, action)
					already[action] = true
				}
			}
		}
	}
	return
}

func NewPermissions(permissionsConfig string) (p *Permissions) {
	p = &Permissions{RolesAccess: make(map[string][]string)}
	rolesConfig := strings.Split(permissionsConfig, ",")
	for _, role := range rolesConfig {
		roleConfig := strings.Split(role, ":")
		roleName, allowedActions := strings.Trim(roleConfig[0], " "), strings.Fields(roleConfig[1])
		p.RolesAccess[roleName] = allowedActions
	}
	return
}
