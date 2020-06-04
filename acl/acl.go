package acl

import "github.com/mikespook/gorbac"

// ACL variables
var (
	roleMerchant = gorbac.NewStdRole("merchant")
	roleCustomer = gorbac.NewStdRole("customer")

	PermissionCatalogWriteAny = gorbac.NewStdPermission("catalog:write:any")
	PermissionCatalogWriteOwn = gorbac.NewStdPermission("catalog:write:own")
	PermissionCatalogReadAny  = gorbac.NewStdPermission("catalog:read:any")
	PermissionCatalogReadOwn  = gorbac.NewStdPermission("catalog:read:own")
)

// RBAC provides interface Role bases access control list
type RBAC interface {
	Can(role string, permission gorbac.Permission) bool
}

type rbac struct {
	rbac *gorbac.RBAC
}

// New returns the role based access control list
func New() RBAC {
	rbacObj := gorbac.New()

	// merchant permissions
	roleMerchant.Assign(PermissionCatalogWriteOwn)
	roleMerchant.Assign(PermissionCatalogReadOwn)

	// customer permissions
	roleCustomer.Assign(PermissionCatalogReadAny)

	rbacObj.Add(roleMerchant)
	rbacObj.Add(roleCustomer)
	return &rbac{rbac: rbacObj}
}

// Can checks if the role has the Permission
func (rbac *rbac) Can(role string, permission gorbac.Permission) bool {
	return rbac.rbac.IsGranted(role, permission, nil)
}
