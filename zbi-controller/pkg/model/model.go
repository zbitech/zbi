package model

import (
	"time"
)

var (
	TESTNET = []string{"testnet.z.cash"}
	MAINNET = []string{"mainnet.z.cash"}

	ZCASH_MAX_CONNECTIONS   = "6"
	ZCASH_RPCCLIENT_TIMEOUT = "30"
	ZCASH_SOLVER            = "tromp"
	ZCASH_MAINNET_SOLVER    = "default"
	ZCASH_TESTNET_SOLVER    = "tromp"
)

type SnapshotRequest struct {
	Version       string `json:"version" validate:"required"`
	VolumeName    string `json:"volumeName"`
	Namespace     string `json:"namespace"`
	SnapshotClass string `json:"snapshotClass"`
	Labels        map[string]string
}

type SnapshotScheduleRequest struct {
	Version          string               `json:"version" validate:"required"`
	Schedule         SnapshotScheduleType `json:"schedule" validate:"required"`
	VolumeName       string               `json:"volumeName"`
	Namespace        string               `json:"namespace"`
	SnapshotClass    string
	BackupExpiration string
	MaxBackupCount   int
	Labels           map[string]string
}

type Project struct {
	Id          string               `json:"id"`
	Name        string               `json:"name"`
	Owner       string               `json:"owner"`
	Blockchain  string               `json:"blockchain"`
	Network     string               `json:"network"`
	Status      string               `json:"status"`
	State       string               `json:"state"`
	Description string               `json:"description"`
	Resources   *KubernetesResources `json:"resources,omitempty"`
	CreatedAt   *time.Time           `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time           `json:"updatedAt,omitempty"`
}

type Instance struct {
	Id           string               `json:"id"`
	Name         string               `json:"name"`
	Type         InstanceType         `json:"type"`
	Project      *Project             `json:"project"`
	Status       string               `json:"status"`
	State        string               `json:"state"`
	InstanceType InstanceType         `json:"instanceType"`
	Owner        string               `json:"owner"`
	Network      NetworkType          `json:"network"`
	Request      *ResourceRequest     `json:"request"`
	Resources    *KubernetesResources `json:"resources,omitempty"`
	CreatedAt    *time.Time           `json:"createdAt,omitempty"`
	UpdatedAt    *time.Time           `json:"updatedAt,omitempty"`
}

type ResourceRequest struct {
	Cpu        string                 `json:"cpu"`
	Memory     string                 `json:"memory"`
	Peers      []string               `json:"peers"`
	Properties map[string]interface{} `json:"properties"`
	Volume     struct {
		Type   DataVolumeType `json:"type"`
		Size   string         `json:"size"`
		Source struct {
			Type DataSourceType `json:"type"`
			Ref  string         `json:"ref"`
		}
	}
}

type InstanceRequest struct {
	Name        string                 `json:"name"`
	Type        InstanceType           `json:"type"`
	Description string                 `json:"description"`
	Peers       []string               `json:"peers"`
	Properties  map[string]interface{} `json:"properties"`
	Volume      struct {
		Type   DataVolumeType `json:"type"`
		Size   string         `json:"size"`
		Source DataSourceType `json:"source"`
		Ref    string         `json:"ref"`
	}
}

type KubernetesResources struct {
	Namespace             *KubernetesResource  `json:"namespace,omitempty"`
	Configmap             *KubernetesResource  `json:"configmap,omitempty"`
	Secret                *KubernetesResource  `json:"secret,omitempty"`
	Persistentvolumeclaim *KubernetesResource  `json:"persistentvolumeclaim,omitempty"`
	Deployment            *KubernetesResource  `json:"deployment,omitempty"`
	Service               *KubernetesResource  `json:"service,omitempty"`
	Httpproxy             *KubernetesResource  `json:"httpproxy,omitempty"`
	Volumesnapshot        []KubernetesResource `json:"volumesnapshot,omitempty"`
	Snapshotschedule      *KubernetesResource  `json:"snapshotschedule,omitempty"`
}

type Activity struct {
	Operation string     `json:"operation"`
	Success   bool       `json:"success"`
	Completed bool       `json:"completed"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type KubernetesResource struct {
	Name       string                 `json:"name,omitempty"`
	Namespace  string                 `json:"namespace,omitempty"`
	Type       ResourceObjectType     `json:"type,omitempty"`
	Status     string                 `json:"status,omitempty"`
	Created    *time.Time             `json:"created,omitempty"`
	Updated    *time.Time             `json:"updated,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	CreatedAt  time.Time              `json:"createdAt,omitempty"`
	UpdatedAt  time.Time              `json:"updatedAt,omitempty"`
}

type Metadata struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	Labels            map[string]string `json:"labels,omitempty"`
	ResourceVersion   string            `json:"resourceVersion,omitempty"`
	Uid               string            `json:"uid,omitempty"`
	CreationTimestamp string            `json:"creationTimestamp"`
}

type IngressCondition struct {
	Prefix string `json:"prefix,omitempty"`
}

type IngressService struct {
	Name string `json:"name"`
	Port int32  `json:"port"`
}

type IngressRoute struct {
	Conditions        []IngressCondition       `json:"conditions,omitempty"`
	Services          []IngressService         `json:"services,omitempty"`
	PathRewritePolicy IngressPathRewritePolicy `json:"pathRewritePolicy,omitempty"`
}

type IngressPathRewritePolicy struct {
	ReplacePrefix []struct {
		Replacement string `json:"replacement,omitempty"`
	} `json:"replacePrefix,omitempty"`
}

type IngressInclude struct {
	Name       string             `json:"name,omitempty"`
	Namespace  string             `json:"namespace,omitempty"`
	Conditions []IngressCondition `json:"conditions,omitempty"`
}

type Ingress struct {
	ApiVersion string   `json:"apiVersion"`
	Metadata   Metadata `json:"metadata"`
	Spec       struct {
		Includes    []IngressInclude `json:"includes,omitempty"`
		Virtualhost struct {
			Fqdn string `json:"fqdn"`
			Tls  struct {
				SecretName string `json:"secretName"`
			} `json:"tls"`
		} `json:"virtualhost,omitempty"`
		Routes []IngressRoute `json:"routes,omitempty"`
	} `json:"spec"`
	Status struct {
		Conditions []struct {
			Errors []struct {
				Message string `json:"message,omitempty"`
				Reason  string `json:"reason,omitempty"`
				Status  string `json:"status,omitempty"`
				Type    string `json:"type,omitempty"`
			} `json:"errors,omitempty"`
			LastTransitionTime string `json:"lastTransitionTime,omitempty"`
			Message            string `json:"message,omitempty"`
			ObservedGeneration int    `json:"observedGeneration,omitempty"`
			Reason             string `json:"reason,omitempty"`
			Status             string `json:"status,omitempty"`
			Type               string `json:"type,omitempty"`
		} `json:"conditions,omitempty"`
		CurrentStatus string `json:"currentStatus,omitempty"`
		Description   string `json:"description,omitempty"`
		LoadBalancer  struct {
			Ingress []struct {
				Hostname string `json:"hostname,omitempty"`
			} `json:"ingress,omitempty"`
		} `json:"loadBalancer,omitempty"`
	} `json:"status,omitempty"`
}

type BlockchainInfo struct {
	Name      string               `json:"name"`
	Networks  []string             `json:"networks"`
	Nodes     []BlockchainNodeInfo `json:"nodes"`
	Templates map[string]string    `json:"templates"`
}

type ImageInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Url     string `json:"url"`
}

type KVPair struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type BlockchainNodeInfo struct {
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Endpoints  map[string][]string    `json:"endpoints"`
	Ports      map[string]int32       `json:"ports"`
	Images     []ImageInfo            `json:"images"`
	Settings   map[string][]KVPair    `json:"settings"`
	Properties map[string]interface{} `json:"properties"`
}

type PolicyInfo struct {
	StorageClass          string `json:"storageClass"`
	SnapshotClass         string `json:"snapshotClass"`
	DomainName            string `json:"domainName"`
	CertificateName       string `json:"certificateName"`
	ServiceAccount        string `json:"serviceAccount"`
	InformerResync        int32  `json:"informerResync"`
	EnableMonitor         bool   `json:"enableMonitor"`
	RequireAuthentication bool   `json:"requireAuthentication"`
	EnvoyConfig           struct {
		Image                 string   `json:"image"`
		Command               []string `json:"command"`
		Timeout               float32  `json:"timeout"`
		AccessAuthorization   bool     `json:"accessAuthorization"`
		AuthServerURL         string   `json:"authServerURL"`
		AuthServerPort        int32    `json:"authServerPort"`
		AuthenticationEnabled bool     `json:"authenticationEnabled"`
	} `json:"envoyConfig"`
	Request struct {
		Cpu     string `json:"cpu"`
		Memory  string `json:"memory"`
		Storage string `json:"storage"`
	} `json:"request"`
}
