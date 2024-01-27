import {Request, Response} from 'express';
import repoFactory from "../repository";
import { getLogger, getDuration } from '../lib/logger';
import { HttpStatusCode } from 'axios';
import { handleError } from '../lib/errors';
import * as types from '../types';
import { Constants } from '../constants';


const createProject = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-create-project');

    try {

        const projectRequest: types.ProjectRequest = request.body;
        const owner = request.headers[Constants.OWNER_ID] as string;

        const projectRepository = repoFactory.getProjectRepository();

        //const user = request.user as types.User;
        projectRequest.owner = owner;

        logger.info(`project request: ${JSON.stringify(projectRequest)}`);

        const id = repoFactory.generateId();
        const name = projectRequest.name;
        const blockchain = projectRequest.blockchain;
        const network = projectRequest.network;
        const description = projectRequest.description as string;

        const project = await projectRepository.createProject(id, name, owner, blockchain, network, description);
        logger.info(`created project: ${JSON.stringify(project)}`);
        response.status(HttpStatusCode.Created).json(project);

    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const findProjects = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-find-projects');

    try {

        const name = request.query.name;
        const value = request.query.value as string;
        const param = name ? { name: value } : {};

        logger.info(`request - ${JSON.stringify(param)}`);
        const projectRepository = repoFactory.getProjectRepository();
        const projects = await projectRepository.findProjects(param);

        response.status(HttpStatusCode.Ok).json(projects);
    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }

}

const findProject = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-find-project');

    try {
        const projectid = request.params.project;

        const projectRepository = repoFactory.getProjectRepository();

        const project = await projectRepository.findProject(projectid);
        response.status(HttpStatusCode.Ok).json(project);

    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
   }

}

const updateProject = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-update-project');

    try {

        const projectid = request.params.project;

        const projectRepository = repoFactory.getProjectRepository();
        let project = await projectRepository.findProject(projectid) as types.Project;

        project = await projectRepository.updateProject(project);
        response.status(HttpStatusCode.Ok).json(project);
        
    } catch (err: any) {
        
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const deleteProject = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-delete-project');

    try {

        const projectid = request.params.project;

        const projectRepository = repoFactory.getProjectRepository();

        await projectRepository.updateInstanceState(projectid, types.StatusType.deleted, types.StateType.deleting);
        response.sendStatus(HttpStatusCode.NoContent);

    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }

}

const purgeProject = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-purge-project');

    try {

        const projectid = request.params.project;

        const projectRepository = repoFactory.getProjectRepository();
        await projectRepository.deleteProject(projectid as string);
        response.sendStatus(HttpStatusCode.NoContent);

    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }

}

const getProjectSnapshots = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctr-get-project-snapshots');

    try {
        const projectid = request.params.project;

        const projectRepository = repoFactory.getProjectRepository();

        const snapshots = await projectRepository.getProjectSnapshots(projectid);
        response.status(HttpStatusCode.Ok).json(snapshots);
    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(err.code).json({message: result.message});
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const createInstance = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-create-instance');

    try {

        const projectid = request.params.project;

        const projectRepository = repoFactory.getProjectRepository();

        const instanceRequest: types.InstanceRequest = request.body;

        let sourceName = "";
        if( instanceRequest.volume?.source === types.VolumeSourceType.volume) {
            const instance = await projectRepository.findInstance(instanceRequest.volume?.ref as string);
            sourceName = instance?.id as string;
        } else if( instanceRequest.volume?.source === types.VolumeSourceType.snapshot ) {
            sourceName = instanceRequest.volume.ref;
        }

        const volumeType = instanceRequest.volume?.type ? instanceRequest.volume.type as types.VolumeType : types.VolumeType.pvc;
        const volumeSource = instanceRequest.volume?.source as types.VolumeSourceType;
        const peers = instanceRequest.peers as string[];
        const properties = instanceRequest.properties;

        const resourceRequest = {peers, properties,
            volume: {
                type: volumeType, size: "15Gi", // add to config
                source: {type: volumeSource,ref: sourceName}
            }     
        };
        
        const instanceId = repoFactory.generateId();
        const name = instanceRequest.name;
        const type = instanceRequest.type;
        const description = instanceRequest.description;

        const newInstance = await projectRepository.createInstance(projectid, instanceId, name, type, resourceRequest);
        response.status(HttpStatusCode.Ok).json(newInstance);
        
    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const findInstances = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-find-instances');

    try {

        const projectid = request.params.project;

        const projectRepository = repoFactory.getProjectRepository();

        const param = { "project": projectid };
        const instances = await projectRepository.findInstances(param);

        response.status(HttpStatusCode.Ok).json(instances);

    } catch (err: any) {
        logger.error(err);
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const getProjectResources = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-get-project-resources');
    try {

        const projectid = request.params.project;
        const { type, name } = request.query;

        const projectRepository = repoFactory.getProjectRepository();

        if (!type && !name) {
            const resources = await projectRepository.getResources(projectid as string);
            response.status(HttpStatusCode.Ok).json(resources);
        } else {
            const resource = await projectRepository.getResourceByName(projectid as string, type as types.ResourceType, name as string);

            if (resource) {
                response.status(HttpStatusCode.Ok).json(resource);
            } else {
                response.status(HttpStatusCode.NotFound).json({ message: `invalid resource parameter` });
            }
        }

    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const updateProjectResource = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-update-instance-resource');

    try {

        const projectid = request.params.project;
        let resource = request.body;

        const projectRepository = repoFactory.getProjectRepository(); 

        resource = await projectRepository.updateResource(projectid as string, resource);

        if(resource.type === types.ResourceType.namespace) {

            if( resource.status === types.StatusType.active ) {
                await projectRepository.updateInstanceState(projectid, types.StatusType.active, types.StateType.available );
            } else {
                await projectRepository.updateInstanceState(projectid, types.StatusType.failed, types.StateType.available );
            }
        }

        response.sendStatus(HttpStatusCode.NoContent);

    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }

}

const addProjectActivity = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-add-project-activity');

    try {

        const projectid = request.params.project;
        const {op} = request.body

        const projectRepository = repoFactory.getProjectRepository();

        logger.info(`operation - ${op}`);
        const activity = await projectRepository.addActivity(projectid, op as types.ActivityType);
        const activities = await projectRepository.getActivities(projectid);
        response.status(HttpStatusCode.Ok).json(activities);

    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });        
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const getProjectActivities = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-get-project-activities');

    try {

        const projectid = request.params.project;

        const projectRepository = repoFactory.getProjectRepository();

        const activities = await projectRepository.getActivities(projectid);
        response.status(HttpStatusCode.Ok).json(activities);

    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });        
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const setProjectPermission = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-set-project-permission');

    try {
        const projectid = request.params.project;
        const userid = request.params.user as string;
        const permission = request.body as types.Permission;

        const projectRepository = repoFactory.getProjectRepository(); 

        await projectRepository.setPermission(projectid, userid, permission);
        const permissions = await projectRepository.getPermissions(projectid);
        response.status(HttpStatusCode.Ok).json(permissions);
        
    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });        
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const removeProjectPermission = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-remove-project-permission');

    try {
        const projectid = request.params.project;
        const userid = request.params.user;

        const projectRepository = repoFactory.getProjectRepository();

        await projectRepository.removePermission(projectid, userid);
        const permissions = await projectRepository.getPermissions(projectid);
        response.status(HttpStatusCode.Ok).json(permissions);
         
    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });        
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const getProjectPermisions = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-get-project-permissions');

    try {
        const projectid = request.params.project;

        const projectRepository = repoFactory.getProjectRepository();

        const permissions = await projectRepository.getPermissions(projectid);
        response.status(HttpStatusCode.Ok).json(permissions);
    } catch (err: any) {
        logger.error(err);
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });        
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const getProjectPermision = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-get-project-permission');

    try {
        const pid = request.params.permision;

        const projectRepository = repoFactory.getProjectRepository();

        const permission = await projectRepository.getPermission(pid);
        response.status(HttpStatusCode.Ok).json(permission);    

    } catch (err: any) {
        logger.error(err);
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });        
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const getProjectUserPermision = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-get-project-permission');

    try {

        const projectid = request.params.project;
        const user = request.params.user;

        const projectRepository = repoFactory.getProjectRepository();

        const permission = await projectRepository.getUserPermission(projectid as string, user);
        response.status(HttpStatusCode.Ok).json(permission);    

    } catch (err: any) {
        logger.error(err);
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });        
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const projectController = {
    createProject,
    findProjects,
    findProject,
    updateProject,
    deleteProject,
    purgeProject,

    getProjectSnapshots,

    createInstance,
    findInstances,

    getProjectResources,
    updateProjectResource,

    addProjectActivity,
    getProjectActivities,

    setProjectPermission,
    removeProjectPermission,
    getProjectPermisions,
    getProjectPermision,
    getProjectUserPermision
}

export default projectController;