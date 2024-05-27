package messages

type Group struct {
	ID       int64
	ParentID int64
	Children []Group
}

type Device struct {
	DeviceUID  string
	DeviceName string
	Group
}
