import mongoose from "mongoose";
import { Activity, BlockchainInfo, Instance, KubernetesResource, KubernetesResources, NodeInfo, Permission, PolicyInfo, Project, ResourceType, User, UserPermissions } from "../../types";

const generateId = () => {
    return (new mongoose.mongo.ObjectId()).toString()
}

const createProject = (project: any): Project => {
    return {
        id: project._id.toString(),
        owner: project.owner ? project.owner.toString() : undefined,
        name: project.name,
        blockchain: project.blockchain,
        network: project.network,
        status: project.status,
        state: project.state,
        description: project.description ? project.description as string : undefined,
        createdAt: project.createdAt ? new Date(project.createdAt) : undefined,
        updatedAt: project.updatedAt ? new Date(project.updatedAt) : undefined
    }
}

const createInstance = (instance: any): Instance => {
    return {
        id: instance._id.toString(),
        name: instance.name,
        type: instance.type,
        project: instance.project ? createProject(instance.project) : undefined,
        description: instance.description,
        request: instance.request,
        status: instance.status,
        state: instance.state,
        createdAt: new Date(instance.createdAt),
        updatedAt: new Date(instance.updatedAt)
    }
}

const createKubernetesResource = (resource: any): KubernetesResource => {
    return {
        name: resource.name,
        type: resource.rtype,
        status: resource.status,
        properties: resource.properties,
        createdAt: resource.createdAt ? new Date(resource.createdAt): undefined,
        updatedAt: resource.updatedAt ? new Date(resource.updatedAt) : undefined
    }
}

const createResources = (resources: any): KubernetesResource[] => {
    return resources.map((resource: any) => createKubernetesResource(resource));
}

const createSnapshotResources = (instances: any): KubernetesResource[] => {
    return instances.map((instance:any) => {
        return instance.resources.volumesnapshot.map((resource: any) => createKubernetesResource(resource))
    }).flat();
}

const createKubernetesResources = (resources: any): KubernetesResources => {

    const namespace = createResources( resources.filter((resource: any) => resource.type === ResourceType.namespace));
    const configmap = createResources( resources.filter((resource: any) => resource.type === ResourceType.configmap));
    const secret = createResources( resources.filter((resource: any) => resource.type === ResourceType.secret));
    const persistentvolumeclaim = createResources( resources.filter((resource: any) => resource.type === ResourceType.persistentvolumeclaim));
    const deployment = createResources( resources.filter((resource: any) => resource.type === ResourceType.deployment));
    const service = createResources( resources.filter((resource: any) => resource.type === ResourceType.service));
    const httpproxy = createResources( resources.filter((resource: any) => resource.type === ResourceType.httpproxy));
    const volumesnapshot = createResources( resources.filter((resource: any) => resource.type === ResourceType.volumesnapshot));
    const snapshotschedule = createResources( resources.filter((resource: any) => resource.type === ResourceType.snapshotschedule));

    return {
        namespace: namespace ? namespace[0] : undefined,
        configmap: configmap ? configmap[0] : undefined,
        secret: secret ? secret[0] : undefined,
        persistentvolumeclaim: persistentvolumeclaim ? persistentvolumeclaim[0] : undefined,
        deployment: deployment ? deployment[0] : undefined,
        service: service ? service[0] : undefined,
        httpproxy: httpproxy ? httpproxy[0] : undefined,
        volumesnapshot: volumesnapshot && volumesnapshot.length > 0 ? volumesnapshot : undefined,
        snapshotschedule: snapshotschedule ? snapshotschedule[0] : undefined,
    }
}

const createActivity = (activity: any): Activity => {
    return {
        id: activity._id.toString(),
        operation: activity.operation,
        completed: activity.completed,
        success: activity.success,
        createdAt: activity.createdAt,
        updatedAt: activity.updatedAt
    }
}

const createActivities = (activities: any): Activity[] => {
    return activities.map((activity: any) => createActivity(activity));
}

const createPermission = (permission: any): Permission => {
    return {
        id: permission.id,
        read: permission.read,
        update: permission.update,
        delete: permission.delete,
        operate: permission.operate,
        access: permission.access
    }
}

const createPermissions = (permissions: any[]): Permission[] => {
    return permissions.map(permission => {
        return createPermission(permission);
    })
}

const createUserPermissions = (perm: any): UserPermissions => {
    return {
        userid: perm.userid,
        projects: createPermissions(perm.projects),
        instances: createPermissions(perm.instances),
        access_token: perm.access_token ? {token: perm.token, date: perm.date}: undefined,
        refresh_token: perm.refreshh_token ? {token: perm.token, date: perm.date} : undefined,
        api_token: perm.api_token ? {token: perm.token, date: perm.date} : undefined,
        limits: {
            createResource: perm.limits.create_resource,
            projects: perm.limits.projects,
            instances: perm.limits.instances,
            updatedAt: perm.limits.updatedAt
        },
        createdAt: perm.createdAt,
        updatedAt: perm.updatedAt
    }
}

const createUser = (user: any): User => {
    return {
        userid: user._id.toString(),
        active: user.active,
        role: user.role,
        email: user.email,
        name: user.name,
        provider: user.provider.name,
        createdAt: new Date(user.createdAt),
        updatedAt: new Date(user.updatedAt)
    }
}

const createBlockchainInfo = (blockchain: any): BlockchainInfo => {
    return {
        name: blockchain.name,
        networks: blockchain.networks,
        templates: blockchain.templates,
        nodes: blockchain.nodes.map((node: any) => createBlockchainNodeInfo(node)),
        createdAt: blockchain.createdAt,
        updatedAt: blockchain.updatedAt
    }
}

const createBlockchainNodeInfo = (node: any): NodeInfo => {
    return {
        name: node.name,
        type: node.type,
        images: node.images,
        endpoints: node.endpoints,
        ports: node.ports,
        settings: node.settings,
        properties: node.properties
    }
}

const createPolicyInfo = (policy: any): PolicyInfo => {
    return {
        storageClass: policy.storageClass,
        snapshotClass: policy.snapshotClass,
        domainName: policy.domainName,
        certificateName: policy.certificateName,
        serviceAccount: policy.serviceAccount,
        envoyConfig: {
            image: policy.envoyConfig.image,
            command: policy.envoyConfig.command,
            timeout: policy.envoyConfig.timeout,
            accessAuthorization: policy.envoyConfig.accessAuthorization,
            authServerURL: policy.envoyConfig.authServerURL,
            authServerPort: policy.envoyConfig.authServerPort,
            authenticationEnabled: policy.envoyConfig.authenticationEnabled
        },
        informerResync: policy.informerResync,
        enableMonitor: policy.enableMonitor,
        requireAuthentication: policy.requireAuthentication,
        request: {
            cpu: policy.request.cpu,
            memory: policy.request.memory,
            storage: policy.request.storage
        }    
    }
}

const decodeMap = (data: any): Map<string, string> => {
    const newData = new Map();
    for( let [key, value] of data) {
        console.log(key);
        const d = Buffer.from(value, 'base64').toString('utf-8');
        console.log(d);
        newData.set(key, d);
    }
    return newData;
}

export {
    generateId, createProject, createUser, createInstance, 
    createKubernetesResource, createResources,
    createKubernetesResources, createActivity, createActivities,
    createPermission, createPermissions, createUserPermissions,
    createSnapshotResources,
    createBlockchainInfo, createBlockchainNodeInfo, createPolicyInfo
}