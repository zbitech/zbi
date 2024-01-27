
export enum RoleType {
    admin = "admin",
    owner = "owner",
    user = "user",
}

export enum BlockchainType {
    zcash = "zcash"
}

export enum NetworkType {
    testnet = 'testnet',
    regnet = 'regnet',
    mainnet = 'mainnet'
}

export enum NodeType {
    zcash = 'zcash',
    lwd = 'lwd',
    zebra = 'zebra'
}

export enum VolumeType {
    ephemeral = 'ephemeral',
    pvc = 'pvc'
};

export enum VolumeSourceType {
    new = 'new',
    volume = 'pvc',
    snapshot = 'snapshot'
}

export enum ResourceType {
    namespace = 'Namespace',
    configmap = 'Configmap',
    secret = 'Secret',
    persistentvolumeclaim = 'PersistentVolumeClaim',
    deployment = 'Deployment',
    service = 'Service',
    httpproxy = 'HTTPproxy',
    volumesnapshot = 'VolumeSnapshot',
    snapshotschedule = 'SnapshotSchedule',
    pod = 'Pod'
}

export enum PermissionType {
    create = "create",
    read = "read",
    update = "update",
    delete = "delete",
    operate = "operate",
    access = "access"
}

export enum ActivityType {
    create = "create",
    start = "start",
    stop = "stop",
    snapshot = "snapshot",
    schedule = "schedule",
    rotate = "rotate",
    delete = "delete",
    purge = "purge",
    repair = "repair",
    update = "update"
}

export enum StateType {
    available = "available",
    creating = "creating",
    updating = "updating",
    starting = "starting",
    stopping = "stopping",
    repairing = "repairing",
    rotating = "rotating",
    creating_snapshot = "creating_snapshot",
    creating_schedule = "creating_schedule",
    deleting = "deleting"
}

export enum StatusType {
    active = "active",
    new = "new",
    running = "running",
    stopped = "stopped",
    deleted = "deleted",
    pending = "pending",
    failed = 'failed',
}

export enum SnapshotScheduleType {
    daily = 'daily',
    weekly = 'weekly',
    monthly = 'monthly'
}

export interface Project {
    id?: string;
    name: string;
    owner?: string;
    blockchain: BlockchainType;
    network: NetworkType;
    status?: string;
    readonly state?: string;
    description?: string;
    createdAt?: Date;
    updatedAt?: Date;
}


export interface Instance {
    id?: string;
    name: string;
    type: NodeType;
    project?: Project;
    description?: string;
    request?: ResourceRequest;
    status?: string;
    readonly state?: string;
    resources?: KubernetesResources;
    activities?: Activity[];
    permissions?: Permission[];
    createdAt?: Date;
    updatedAt?: Date;
}

export interface Activity {
    id?: string;
    operation: ActivityType;
    success: boolean;
    completed: boolean;
    createdAt?: Date;
    updatedAt?: Date;
}

export interface ResourceRequest {
    cpu?: string;
    memory?: string;
    peers?: string[];
    properties: Map<string, any>;
    volume: {
        type: VolumeType;
        size?: string;
        source: {
            type: VolumeSourceType;
            ref?: string;
        }
    }
}

export interface KubernetesResource {
    name: string;
    type: ResourceType;
    status: string;
    properties: Map<String, any>;
    createdAt?: Date;
    updatedAt?: Date;
}

export interface KubernetesResources {
    namespace?: KubernetesResource,
    configmap?: KubernetesResource,
    secret?: KubernetesResource,
    persistentvolumeclaim?: KubernetesResource,
    deployment?: KubernetesResource,
    service?: KubernetesResource,
    httpproxy?: KubernetesResource,
    volumesnapshot?: KubernetesResource[];
    snapshotschedule?: KubernetesResource;
}

export interface Permission {
    id?: string;
    read: boolean;
    update: boolean;
    delete: boolean;
    operate?: boolean;
    access?: boolean
}

export interface UserPermissions {
    userid: string;
    projects: Permission[];
    instances: Permission[];
    limits: UserLimits,
    access_token?: {
        token: string;
        date: Date;
    },
    refresh_token?: {
        token?: string;
        date: Date;
    }
    api_token?: {
        token: string;
        date: Date;
    }
    createdAt?: Date;
    updatedAt?: Date;
}

export interface UserLimits {
    createResource: boolean;
    projects: number;
    instances: number;
    updatedAt?: Date;
}

export interface User {
    userid: string;
    active: boolean;
    role: RoleType;
    email: string;
    name: string;
    picture?: string;
    provider?: string;
    createdAt?: Date;
    updatedAt?: Date;
}

export interface BlockchainInfo {
    name: string;
    networks: string[];
    nodes: NodeInfo[];
    templates?: Map<string, string>;
    createdAt?: Date;
    updatedAt?: Date;
}

export interface NodeInfo {
    name: string;
    type: string;
    images: [{
        name: string;
        version: string;
        url: string;
    }];
    endpoints: any;
    ports: any;
    settings: any;
    properties: any;
    createdAt?: Date;
    updatedAt?: Date;
}

export interface PolicyInfo {
    storageClass: string,
    snapshotClass: string,
    domainName: string,
    certificateName: string,
    serviceAccount: string,
    envoyConfig: {
        image: string,
        command: string[],
        timeout: number,
        accessAuthorization: boolean,
        authServerURL: string,
        authServerPort: number,
        authenticationEnabled: boolean
    },
    informerResync: number,
    enableMonitor: boolean,
    requireAuthentication: boolean,
    request?: {
        cpu: string,
        memory: string,
        storage: string
    }
}

export interface ProjectRequest {
    name: string;
    owner?: string;
    blockchain: BlockchainType;
    network: NetworkType;
    description?: string;
}


export interface InstanceRequest {
    name: string;
    type: NodeType;
    description: string;
    peers?: Array<string>;
    volume?: {
        type: VolumeType;
        size?: string;
        source: VolumeSourceType;
        ref: string;
    };
    properties: any;
}


export interface SnapshotScheduleRequest {
    schedule: SnapshotScheduleType;
    hourOfDay?: number;
    dayOfWeek?: number;
    dayOfMonth?: number;
}
