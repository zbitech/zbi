import { Activity, ActivityType, BlockchainType, Instance, KubernetesResource, KubernetesResources, NetworkType, NodeType, Permission, Project, ResourceRequest, ResourceType, StateType, StatusType} from "../../types";
import { activityModel, instanceModel, permissionModel, projectModel, resourceModel } from "./schema";
import { FilterQuery, Types } from "mongoose";
import * as fn from "./fn";
import { getDuration, getLogger } from "../../lib/logger";
import { AppError, ItemNotFoundError } from "../../lib/errors";

const createProject = async (id: string, name: string, owner: string, blockchain: BlockchainType, network: NetworkType, description: string): Promise<Project> => {
    let logger = getLogger('repo-create-project');
    try {
        const proj = new projectModel({
            id, name, owner, blockchain, network, description, status: "new"
        });
        if (proj) {
            await proj.save();
            return fn.createProject(proj);
        }
        throw new AppError("failed to create project");
    } catch (err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const checkProjectId = async (id: string): Promise<boolean> => {
    let logger = getLogger('repo-check-project');
    try {
        const project = await projectModel.findById(id, {_id: 1});
        return project ? true : false;
    } catch (err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const findProjects = async (query: any): Promise<Project[]> => {
    let logger = getLogger('repo-find-projects');
    try {
        const filter = query as FilterQuery<any>;
        const projects = await projectModel.find(query);
        if (projects) {
            return projects.map((project: any) => {
                return fn.createProject(project);
            })
        }
        return [];
    } catch (err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const findProject = async (id: string): Promise<Project | undefined> => {
    let logger = getLogger('repo-find-project');
    try {
        const project = await projectModel.findById(id);
        if (project) {
            return fn.createProject(project);
        }

        return undefined;
    } catch (err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const findProjectByName = async (name: string): Promise<Project | undefined> => {
    let logger = getLogger('repo-find-project-byname');
    try {
        const project = await projectModel.findOne({ name });
        if (project) {
            return fn.createProject(project);
        }

        return undefined;
    } catch (err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const updateProject = async (project: Project): Promise<Project> => {
    let logger = getLogger('repo-update-project');
    try {
        const p = await projectModel.findById(project.id);
        if (p) {
            p.description = project.description;
            await p.save();
            return fn.createProject(p);
        }

        throw new ItemNotFoundError("project not found");
    } catch (err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

// const updateProjectStatus = async (id: string, status: string): Promise<Project> => {
//     let logger = getLogger('repo-update-project-status');
//     try {
//         const p = await projectModel.findById(id);
//         if (p) {
//             p.status = status;
//             await p.save();
//             return fn.createProject(p);
//         }

//         throw new ItemNotFoundError("project not found");
//     } catch (err: any) {
//         throw err;
//     }
// }

const updateProjectState = async (id: string, status: StatusType, state: StateType): Promise<Project> => {
    let logger = getLogger('repo-update-instance-state');
    try {
        const project = await projectModel.findByIdAndUpdate(id, {$set: {status, state}});
        if (project) {
            return fn.createProject(project);
        }
        throw new Error("project not found");
    } catch (err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const getProjectSnapshots = async (project: string): Promise<KubernetesResource[]> => {
    let logger = getLogger('repo-get-instance-resource');
    try {

        const instances = await instanceModel.find({project}, {_id: 1});

        if( instances ) {
            const resources = await resourceModel.find({
                type: ResourceType.volumesnapshot, object: {$in: instances}
            });
            return fn.createResources(instances);
        }

        return []
    } catch(err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const deleteProject = async (id: string): Promise<void> => {
    let logger = getLogger('repo-delete-project');
    try {
        const project = await projectModel.findByIdAndDelete(id);
    } catch (err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const createInstance = async (project: string, id: string, name: string, type: NodeType, request: ResourceRequest): Promise<Instance> => {
    let logger = getLogger('repo-create-instance');
    try {
        const instance = new instanceModel({ id, project, name, type, request, status: StatusType.new});
        await instance.save();
        return fn.createInstance(instance);
    } catch (err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const findInstances = async (query: any): Promise<Instance[]> => {
    let logger = getLogger('repo-find-instances');
    try {
        const filter = query as FilterQuery<any>;
        const instances = await instanceModel.find(query, {_id: 1, name: 1, type: 1, network: 1, description: 1, request: 1, state: 1, createdAt: 1, updatedAt: 1});
        logger.debug(`found instances - ${instances}`);
        if (instances) {
            return instances.map((instance: any) => {
                return fn.createInstance(instance);
            })
        }
        return [];
    } catch (err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const checkInstanceId = async (id: string): Promise<boolean> => {
    let logger = getLogger('repo-check-instance');
    try {
        const instance = await instanceModel.findById(id, {_id: 1});
        return instance ? true : false;
    } catch (err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const findInstance = async (id: string): Promise<Instance | undefined> => {
    let logger = getLogger('repo-find-instance');
    try {
        const instance = await instanceModel.findById(id)
                                    .populate({path: "project",
                                               model: "project",
                                               select: {"_id":1, "name": 1, "blockchain": 1, "network": 1, "status": 1, "owner": 1 }});

        logger.debug(`found instance - ${instance}`);
        if (instance) {
            logger.debug(`project - ${instance.project}`);
            const newInstance = fn.createInstance(instance);
            newInstance.resources = await getResources(instance.id);
            newInstance.activities = await getActivities(instance.id);

            return newInstance;
        }

        return undefined;
    } catch (err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const findInstanceByName = async (project: string, name: string): Promise<Instance | undefined> => {
    let logger = getLogger('repo-find-instance-byname');
    try {
        const instance = await instanceModel.findOne({ project, name })
                                    .populate({path: "project",
                                               model: "project",
                                               select: {"_id":1, "name": 1, "blockchain": 1, "network": 1, "status": 1 }});
        if (instance) {
            const newInstance = fn.createInstance(instance);
            newInstance.resources = await getResources(instance.id);
            newInstance.activities = await getActivities(instance.id);
            return newInstance;
        }

        return undefined;
    } catch (err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const updateInstance = async (instance: Instance): Promise<Instance> => {
    let logger = getLogger('repo-update-instance');
    try {
        const _instance = await instanceModel.findById(instance.id);
        if (_instance) {
            _instance.description = instance.description;
            _instance.state = instance.state;
            _instance.request = {
                cpu: instance.request?.cpu ? instance.request.cpu : undefined,
                memory: instance.request?.memory ? instance.request.memory : undefined,
                peers: instance.request?.peers as string[],
                properties: instance.request?.properties,
            }

            await _instance.save();

            return fn.createInstance(_instance);
        }

        throw new Error("instance not found");
    } catch (err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const updateInstanceState = async (id: string, status: StatusType, state: StateType): Promise<Instance> => {
    let logger = getLogger('repo-update-instance-state');
    try {
        const instance = await instanceModel.findByIdAndUpdate(id, {$set: {status, state}});
        if (instance) {
            return fn.createInstance(instance);
        }
        throw new Error("instance not found");
    } catch (err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const deleteInstance = async (id: string): Promise<void> => {
    let logger = getLogger('repo-delete-instance');
    try {
        const instance = await instanceModel.findByIdAndDelete(id);
        if (!instance) {
            throw new Error("instance not found");
        }
    } catch (err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const deleteVolumeSnapshot = async (object: string, name: string): Promise<KubernetesResources> => {
    let logger = getLogger('repo-delete-volume-snapshot');
    try {

        const resource = await resourceModel.findOneAndDelete({object, type: ResourceType.volumesnapshot, name});
        if( resource ) {
            return getResources(object);
        }

        throw new Error("resource not found");
    } catch (err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const addActivity = async (object: string, operation: ActivityType): Promise<Activity> => {
    let logger = getLogger('repo-add-instance-activity');
    try {

        const activity = await activityModel.create({operation, object, completed: false, success: false});
        return fn.createActivity(activity);
    } catch(err: any) {
        logger.error(err);
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const updateActivity = async (id: string, success: boolean): Promise<Activity> => {
    let logger = getLogger('repo-update-activity');
    try {

        const activity = await activityModel.findByIdAndUpdate(id, {success, completed: true});
        if(activity) {
            return fn.createActivity(activity);
        }

        throw new Error("activity not found");

    } catch(err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const getActivity = async (id: string): Promise<Activity> => {
    let logger = getLogger('repo-get-activity');

    try {
        const activity = await activityModel.findById(id);
        if(activity) {
            return fn.createActivity(activity);
        }
        throw new Error("instance not found");

    } catch(err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }

}

const getActivities = async (object: string): Promise<Activity[]> => {
    let logger = getLogger('repo-get-activities');
    try {
        const activities = await activityModel.find({object}).sort({createdAt: 1});
        logger.debug(`got back activities - ${JSON.stringify(activities)}`);
        if(activities) {
            return fn.createActivities(activities);
        }

        return []
    } catch(err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const getPermissions = async (object: string): Promise<Permission[]> => {
    let logger = getLogger('repo-get-permissions');
    try {
        const permissions = await permissionModel.find({object});
        if(permissions) {
            return fn.createPermissions(permissions);
        }

        throw new Error("instance not found");
    } catch(err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const getPermission = async(id: string): Promise<Permission> => {
    let logger = getLogger('repo-get-permission');
    try {
        const permission = await permissionModel.findById(id);
        if(permission) {
            return fn.createPermission(permission);
        }

        throw new Error("instance not found");
    } catch(err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const getUserPermission = async(object: string, user: string): Promise<Permission> => {
    let logger = getLogger('repo-get-user-permission');
    try {
        const permission = await permissionModel.findOne({object, user});
        if(permission) {
            return fn.createPermission(permission);
        }

        throw new Error("instance not found");
    } catch(err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const setPermission = async (object: string, user: string, permission: Permission): Promise<Permission> => {
    let logger = getLogger('repo-set-permissions');
    try {

        const p = await permissionModel.findOneAndUpdate({object, user: user}, {...permission});
        if(!p) {
            return fn.createPermission( await permissionModel.create({...permission, object}) );
        }
        return fn.createPermission(p);
    } catch(err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }

}

const removePermission = async (object: string, user: string): Promise<void> => {
    let logger = getLogger('repo-remove-permissions');
    try {
        const p = await permissionModel.findOneAndRemove({object, user});
        if(p) {
            throw new Error("permision not found");
        }
    } catch(err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const getResources = async (object: string): Promise<KubernetesResources> => {
    let logger = getLogger('repo-get-resources');
    try {
        const resources = await resourceModel.find({object});
        if(resources) {
            return fn.createKubernetesResources(resources);
        }

        throw new Error("resources not found");
    } catch(err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const getResource = async (id: string): Promise<KubernetesResource> => {
    let logger = getLogger('repo-get-resource');
    try {
        const resource = await resourceModel.findById(id);
        if(resource) {
            return fn.createKubernetesResource(resource);
        }

        throw new Error("resource not found");
    } catch(err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const getResourceByName = async (object: string, type: ResourceType, name: string): Promise<KubernetesResource> => {
    let logger = getLogger('repo-get-resource[name]');
    try {
        const resource = await resourceModel.find({object, type, name});
        if(resource) {
            return fn.createKubernetesResource(resource);
        }

        throw new Error("resource not found");
    } catch(err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const createResource = async (object: string, type: ResourceType, name: string): Promise<KubernetesResource> => {
    let logger = getLogger('repo-create-resource');
    try {
        const resource = await resourceModel.create({object, type, name});
        if(resource) {
            return fn.createKubernetesResource(resource);
        }

        throw new Error("resource not found");
    } catch(err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const updateResource = async (object: string, resource: KubernetesResource): Promise<KubernetesResource> => {
    let logger = getLogger('repo-get-resource');
    try {
        const type = resource.type;
        const name = resource.name;
        const status = resource.status;
        const properties = resource.properties;

        const rsc = await resourceModel.findOneAndUpdate({object, type, name}, {status, properties});
        if(rsc) {
            return fn.createKubernetesResource(resource);
        } else {
            return fn.createKubernetesResource( await resourceModel.create({object, type, name, status, properties}) );
        }

    } catch(err: any) {
        throw err;
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

    
const projectMongoRepository = {
    checkProjectId,
    checkInstanceId,
    createProject,
    findProjects,
    findProject,
    findProjectByName,
    updateProject,
    updateProjectState,
    getProjectSnapshots,
    deleteProject,
    createInstance,
    findInstances,
    findInstance,
    findInstanceByName,
    updateInstance,
    updateInstanceState,
    deleteInstance,
 
    deleteVolumeSnapshot,
    addActivity,
    getActivity,
    updateActivity,
    getActivities,
    getPermissions,
    getPermission,
    getUserPermission,
    setPermission,
    removePermission,
    getResources,
    getResource,
    getResourceByName,
    createResource,
    updateResource
}

export default projectMongoRepository