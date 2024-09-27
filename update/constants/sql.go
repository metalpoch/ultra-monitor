package constants

const (
	TABLE_DEVICES              string = "devices"
	JOIN_TEMPLATES_ON_DEVICES  string = "JOIN templates ON devices.template_id = templates.id"
	GLOBAL_COLUMN_CREATED_AT   string = "created_at"
	GLOBAL_COLUMN_UPDATED_AT   string = "updated_at"
	INTERFACE_COLUMN_IF_INDEX  string = "if_index"
	INTERFACE_COLUMN_IF_NAME   string = "if_name"
	INTERFACE_COLUMN_IF_DESCR  string = "if_descr"
	INTERFACE_COLUMN_IF_ALIAS  string = "if_alias"
	INTERFACE_COLUMN_DEVICE_ID string = "device_id"

	MEASUREMENT_COLUMN_IF_INDEX string = "interface_id"
)
