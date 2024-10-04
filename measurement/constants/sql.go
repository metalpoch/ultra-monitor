package constants

const (
	TABLE_DEVICES                   string = "devices"
	JOIN_TEMPLATES_ON_DEVICES       string = "JOIN templates ON devices.template_id = templates.id"
	SELECT_TEMPLATES_ON_DEVICES     string = "devices.id as id, devices.ip, devices.sys_name, devices.sys_location, devices.community, devices.is_alive, devices.template_id, devices.created_at, devices.updated_at, templates.oid_bw, templates.oid_in, templates.oid_out, devices.last_check"
	GLOBAL_COLUMN_CREATED_AT        string = "created_at"
	GLOBAL_COLUMN_UPDATED_AT        string = "updated_at"
	INTERFACE_COLUMN_IF_INDEX       string = "if_index"
	INTERFACE_COLUMN_IF_NAME        string = "if_name"
	INTERFACE_COLUMN_IF_DESCR       string = "if_descr"
	INTERFACE_COLUMN_IF_ALIAS       string = "if_alias"
	INTERFACE_COLUMN_DEVICE_ID      string = "device_id"
	MEASUREMENT_COLUMN_INTERFACE_ID string = "interface_id"
)
