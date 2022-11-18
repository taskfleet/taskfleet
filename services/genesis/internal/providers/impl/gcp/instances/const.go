package gcpinstances

const (
	// LabelID captures the internal ID of the instance.
	LabelID string = "id"
	// LabelKeyCreatedBy identifies the component that created an instance.
	LabelKeyCreatedBy string = "created-by"
	// LabelKeyOwnedBy describes the label key used for identifying the owner of an instance.
	LabelKeyOwnedBy string = "owned-by"
	// LabelKeyDeviceName describes the label key used for identifying a disk on an instance.
	LabelKeyDeviceName string = "device-name"
)
