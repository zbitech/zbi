package manager

import (
	"testing"

	"github.com/stretchr/testify/assert"
	fake_zbi "github.com/zbitech/controller/fake-zbi"
	"github.com/zbitech/controller/internal/helper"
	"github.com/zbitech/controller/pkg/file_template"
	"github.com/zbitech/controller/pkg/model"
)

const App_Template_Content = `
{{define "VOLUME"}}
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{.VolumeName}}
  namespace: {{.Namespace}}
  labels:
    {{- range $key, $value := .Labels}}
    {{$key}}: {{$value}}
    {{- end}}
spec:
  accessModes:
    - ReadWriteOnce
{{- if .StorageClass}}
  storageClassName: {{.StorageClass}}
{{- end}}
{{- if eq .VolumeSourceType "pvc" }}
  dataSource:
    name: {{.SourceName}}
    kind: PersistentVolumeClaim
{{- else if eq .VolumeSourceType "snapshot" }}
  dataSource:
    name: {{.SourceName}}
    kind: VolumeSnapshot
    apiGroup: snapshot.storage.k8s.io
{{- end}}
  resources:
    requests:
      storage: {{.Size}}Gi
{{end}}

{{define "SNAPSHOT"}}
apiVersion: snapshot.storage.k8s.io/v1
kind: VolumeSnapshot
metadata:
  name: {{.SnapshotName}}
  namespace: {{.Namespace}}
  labels:
    {{- range $key, $value := .Labels}}
    {{$key}}: {{$value}}
    {{- end}}
spec:
  volumeSnapshotClassName: {{.SnapshotClass}}
  source:
    persistentVolumeClaimName: {{.VolumeName}}
{{end}}

{{define "SCHEDULE_SNAPSHOT"}}
apiVersion: snapscheduler.backube/v1
kind: SnapshotSchedule
metadata:
  name: {{.ScheduleName}}
  namespace: {{.Namespace}}
  labels:
    {{- range $key, $value := .Labels}}
    {{$key}}: {{$value}}
    {{- end}}
spec:
  claimSelector:
    matchLabels:
      {{- range $key, $value := .Labels}}
      {{$key}}: {{$value}}
      {{- end}}
  disabled: false
  retention:
    expires: "{{.BackupExpiration}}"
    maxCount: {{.MaxBackupCount}}
  schedule: "{{.Schedule}}"
  snapshotTemplate:
    labels:
      {{- range $key, $value := .Labels}}
      {{$key}}: {{$value}}
      {{- end}}
      name: {{.Name}}
      schedule: {{.ScheduleType}}
    snapshotClassName: {{.SnapshotClass}}
{{end}}
`

func TestAppResourceManager_NewAppResourceManager(t *testing.T) {
	ctx := fake_zbi.InitContext()

	fileTemplate, err := file_template.NewFileTemplate("app", App_Template_Content, file_template.NO_FUNCS)
	assert.NoError(t, err)
	assert.NotNil(t, fileTemplate)

	helper.LoadConfig(ctx)

	var manager = NewAppResourceManager()
	assert.NotNil(t, manager)
}

func TestAppResourceManager_CreateVolumeResource(t *testing.T) {

	ctx := fake_zbi.InitContext()
	helper.LoadConfig(ctx)

	fileTemplate := helper.GetAppTemplate()
	assert.NotNil(t, fileTemplate)

	var manager = NewAppResourceManager()
	assert.NotNil(t, manager)

	volume := model.VolumeSpec{
		VolumeName:     "main",
		StorageClass:   "csi",
		Namespace:      "test",
		VolumeDataType: "type",
		DataSourceType: model.NoDataSource,
		SourceName:     "",
		Size:           "10Gi",
		Labels:         map[string]string{"instance": "instance"},
	}

	resources, err := manager.CreateVolumeResource(ctx, volume)
	assert.NoError(t, err)
	assert.NotNil(t, resources)
}

func TestAppResourceManager_CreateSnapshotResource(t *testing.T) {
	ctx := fake_zbi.InitContext()

	helper.LoadConfig(ctx)
	//	helper.Config.LoadTemplates(ctx)

	fileTemplate := helper.GetAppTemplate()
	assert.NotNil(t, fileTemplate)

	var manager = NewAppResourceManager()
	assert.NotNil(t, manager)

	req := &model.SnapshotRequest{
		Version:       "v1",
		VolumeName:    "volume",
		Namespace:     "namespace",
		SnapshotClass: "snapshot",
		Labels:        map[string]string{},
	}

	resources, err := manager.CreateSnapshotResource(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resources)
}

func TestAppResourceManager_CreateSnapshotScheduleResource(t *testing.T) {
	ctx := fake_zbi.InitContext()

	helper.LoadConfig(ctx)

	//	policy := helper.GetPolicyInfo()
	fileTemplate := helper.GetAppTemplate()
	assert.NotNil(t, fileTemplate)

	var manager = NewAppResourceManager()
	assert.NotNil(t, manager)

	req := &model.SnapshotScheduleRequest{
		Version:          "v1",
		Schedule:         "",
		VolumeName:       "volume",
		Namespace:        "namespace",
		SnapshotClass:    "snapshot",
		BackupExpiration: "24h",
		MaxBackupCount:   5,
		Labels:           map[string]string{},
	}

	resources, err := manager.CreateSnapshotScheduleResource(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resources)
}
