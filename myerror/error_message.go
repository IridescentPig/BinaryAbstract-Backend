package myerror

const (
	INVALID_BODY_INFO                    = "Invalid request body"
	USER_NOT_FOUND_INFO                  = "User not found"
	PERMISSION_DENIED_INFO               = "Permission denied"
	DUPLICATED_NAME_INFO                 = "Dupilicated name"
	INVALID_PARAM_INFO                   = "Invalid param in router"
	TOKEN_EMPTY_INFO                     = "Cannot find token in request header"
	TOKEN_EXPIRED_INFO                   = "Token has expired"
	TOKEN_INVALID_INFO                   = "Invaild token"
	ENTITY_NOT_FOUND_INFO                = "Entity not found"
	USER_HAS_EXISTED_INFO                = "User has existed"
	USER_NOT_IN_ENTITY_INFO              = "User does not exist in this entity"
	NAME_CANNOT_EMPTY_INFO               = "Name cannot be empty"
	DEPARTMENT_NOT_FOUND_INFO            = "department not found"
	DEPARTMENT_NOT_IN_ENTITY_INFO        = "Department not in entity"
	USER_NOT_IN_DEPARTMENT_INFO          = "User not in department"
	ENTITY_HAS_USERS_INFO                = "Entity has users, cannot be deleted"
	DELETE_USER_SELF_INFO                = "You cannot delete yourself"
	DEPARTMENT_HAS_USERS_INFO            = "Department has users, cannot be deleted"
	ASSET_CLASS_NOT_FOUND_INFO           = "Asset class not found"
	PARENT_ASSET_CLASS_NOT_FOUND_INFO    = "Parent asset class not found"
	INVALID_TYPE_OF_CLASS_INFO           = "Invalid type of asset class"
	PARENT_CANNOOT_BE_SUCCESSOR_INFO     = "Cannot set current's parent to it's successor"
	CLASS_HAS_ASSET_INFO                 = "There exist some assets belong to this asset class, cannot delete it"
	ASSET_NOT_FOUND_INFO                 = "Asset not found"
	ASSET_NOT_IN_DEPARTMENT_INFO         = "Asset not in department"
	PARENT_ASSET_NOT_FOUND_INFO          = "Parent asset not found"
	TARGET_USER_NOT_FOUND_INFO           = "Target user not found"
	NOT_IN_SAME_ENTITY_INFO              = "Not in the same entity"
	TARGET_NOT_DEPARTMENT_SUPER_INFO     = "Target user is not department manager"
	CLASS_HAS_SUB_CLASS_INFO             = "Class has sub class"
	PRICE_OUT_OF_RANGE_INFO              = "Price out of range"
	ENTITY_NAME_CANNOT_BE_EMPTY_INFO     = "Entity name cannot be empty"
	DEPARTMENT_NAME_CANNOT_BE_EMPTY_INFO = "Department name cannot be empty"
	TARGET_EMPTY_INFO                    = "Target user cannot be empty"
	ASSET_LIST_INVALID_INFO              = "Asset list invalid"
	TASK_NOT_FOUND_INFO                  = "Task not found"
	TASK_NOT_IN_DEPARTMENT_INFO          = "Task not in department"
)
