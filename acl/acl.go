package acl

import "github.com/mikespook/gorbac"

// ACL variables
var (
	roleMerchant = gorbac.NewStdRole("merchant")
	roleCustomer = gorbac.NewStdRole("customer")

	PermissionCreateRestaurantAny = gorbac.NewStdPermission("createRestaurant:any")
	PermissionCreateRestaurantOwn = gorbac.NewStdPermission("createRestaurant:own")
	PermissionGetRestaurantAny    = gorbac.NewStdPermission("getRestaurant:any")
	PermissionGetRestaurantOwn    = gorbac.NewStdPermission("getRestaurant:own")
	PermissionDeleteRestaurantAny = gorbac.NewStdPermission("deleteRestaurant:any")
	PermissionDeleteRestaurantOwn = gorbac.NewStdPermission("deleteRestaurant:own")

	PermissionCreateProductAny = gorbac.NewStdPermission("createProduct:any")
	PermissionCreateProductOwn = gorbac.NewStdPermission("createProduct:own")
	PermissionGetProductAny    = gorbac.NewStdPermission("getProduct:any")
	PermissionGetProductOwn    = gorbac.NewStdPermission("getProduct:own")
	PermissionDeleteProductAny = gorbac.NewStdPermission("deleteProduct:any")
	PermissionDeleteProductOwn = gorbac.NewStdPermission("deleteProduct:own")
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
	roleMerchant.Assign(PermissionCreateRestaurantOwn)
	roleMerchant.Assign(PermissionGetRestaurantOwn)
	roleMerchant.Assign(PermissionDeleteRestaurantOwn)
	roleMerchant.Assign(PermissionCreateProductOwn)
	roleMerchant.Assign(PermissionGetProductOwn)
	roleMerchant.Assign(PermissionDeleteProductOwn)

	// customer permissions
	roleCustomer.Assign(PermissionGetRestaurantAny)
	roleCustomer.Assign(PermissionGetProductAny)

	rbacObj.Add(roleMerchant)
	rbacObj.Add(roleCustomer)
	return &rbac{rbac: rbacObj}
}

// Can checks if the role has the Permission
func (rbac *rbac) Can(role string, permission gorbac.Permission) bool {
	return rbac.rbac.IsGranted(role, permission, nil)
}
