package model

type ResourceObjectType string

const (
	ResourceNamespace             ResourceObjectType = "Namespace"
	ResourceDeployment            ResourceObjectType = "Deployment"
	ResourceService               ResourceObjectType = "Service"
	ResourceConfigMap             ResourceObjectType = "ConfigMap"
	ResourceSecret                ResourceObjectType = "Secret"
	ResourcePod                   ResourceObjectType = "Pod"
	ResourcePersistentVolume      ResourceObjectType = "PersistentVolume"
	ResourcePersistentVolumeClaim ResourceObjectType = "PersistentVolumeClaim"
	ResourceVolumeSnapshot        ResourceObjectType = "VolumeSnapshot"
	ResourceVolumeSnapshotClass   ResourceObjectType = "VolumeSnapshotClass"
	ResourceSnapshotSchedule      ResourceObjectType = "SnapshotSchedule"
	ResourceHTTPProxy             ResourceObjectType = "HTTPProxy"
)

type EventAction string

const (
	EventActionCreate         EventAction = "create"
	EventActionDelete         EventAction = "delete"
	EventActionUpdate         EventAction = "update"
	EventActionRepair         EventAction = "repair"
	EventActionSnapshot       EventAction = "snapshot"
	EventActionSchedule       EventAction = "schedule"
	EventActionStopInstance   EventAction = "stop"
	EventActionStartInstance  EventAction = "start"
	EventActionRotate         EventAction = "rotate"
	EventActionDeleteResource EventAction = "delete_resource"
)

type NetworkType string

const (
	NetworkTypeMain NetworkType = "mainnet"
	NetworkTypeTest NetworkType = "testnet"
)

type InstanceType string

const (
	InstanceTypeZCASH InstanceType = "zcash"
	InstanceTypeLWD   InstanceType = "lwd"
)

type StatusType string

const (
	StatusNew         StatusType = "new"
	StatusActive      StatusType = "active"
	StatusInActive    StatusType = "inactive"
	StatusFailed      StatusType = "failed"
	StatusPending     StatusType = "pending"
	StatusProgressing StatusType = "progressing"
	StatusBound       StatusType = "bound"
	StatusRunning     StatusType = "running"
	StatusStopped     StatusType = "stopped"
	StatusValid       StatusType = "valid"
	StatusReady       StatusType = "ready"
)

type SnapshotScheduleType string

const (
	HourlySnapshotScheduled SnapshotScheduleType = "hourly"
	DailySnapshotSchedule   SnapshotScheduleType = "daily"
	WeeklySnapshotSchedule  SnapshotScheduleType = "weekly"
	MonthlySnapshotSchedule SnapshotScheduleType = "monthly"
)

type DataSourceType string

const (
	NoDataSource       DataSourceType = "none"
	NewDataSource      DataSourceType = "new"
	VolumeDataSource   DataSourceType = "pvc"
	SnapshotDataSource DataSourceType = "snapshot"
)

type DataVolumeType string

const (
	EphemeralDataVolume  DataVolumeType = "ephemeral"
	PersistentDataVolume DataVolumeType = "pvc"
)
