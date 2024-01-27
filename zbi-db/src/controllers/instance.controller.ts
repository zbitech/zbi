import {Request, Response} from 'express';
import repoFactory from "../repository";
import { getLogger, getDuration } from '../lib/logger';
import { HttpStatusCode } from 'axios';
import { handleError } from '../lib/errors';
import * as types from '../types';
import { Constants } from '../constants';

const findInstance = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-find-instance');

    try {

        let instanceid = request.params.instance;

        const projectRepository = repoFactory.getProjectRepository();
        let instance = await projectRepository.findInstance(instanceid);
        if(instance) {
            logger.info(`found instance - ${JSON.stringify(instance)}`);
            response.status(HttpStatusCode.Ok).json(instance);
        } else {
            response.status(HttpStatusCode.NotFound).json({message: "instance not found"});
        }
    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const updateInstance = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-update-instance');

    try {
        let instanceid = request.params.instance;

        const projectRepository = repoFactory.getProjectRepository();
        let instance = await projectRepository.findInstance(instanceid);
        
        const instanceRequest: types.InstanceRequest = request.body;
        if(instance) {
            logger.info(`project name = ${instance.project?.name}, instance = ${instance.name}, request = ${JSON.stringify(instanceRequest)}`);
            if(instance.request) {
                instance.request.peers = instanceRequest.peers;
            }
            instance = await projectRepository.updateInstance(instance);
            response.status(HttpStatusCode.Ok).json( instance );
        } else {
            response.status(HttpStatusCode.NotFound).json({message: "instance not found"});
        }
    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const deleteInstance = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-delete-instance');

    try {

        let instanceid = request.params.instance;

        const projectRepository = repoFactory.getProjectRepository();
        await projectRepository.updateInstanceState(instanceid, types.StatusType.deleted, types.StateType.deleting);

        response.sendStatus(HttpStatusCode.NoContent);
    } catch (err: any) {
        logger.error(`failed to create project: ${err}`)
        response.status(500).json({ message: err.message });
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const purgeInstance = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-purge-instance');

    try {

        let instanceid = request.params.instance;

        const projectRepository = repoFactory.getProjectRepository();

        await projectRepository.deleteInstance(instanceid);
        response.sendStatus(HttpStatusCode.NoContent);

    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const getInstanceResources = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-get-instance-resources');
    try {

        const instanceid = request.params.instance;

        const projectRepository = repoFactory.getProjectRepository();

        const { type, name } = request.query;
        if (!type && !name) {
            logger.info(`get resources for instance - ${instanceid}`);
            const resources = await projectRepository.getResources(instanceid as string);
            response.status(HttpStatusCode.Ok).json(resources);
        } else {
            const resource = await projectRepository.getResourceByName(instanceid, type as types.ResourceType, name as string);
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

const updateInstanceResource = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-update-instance-resource');

    try {

        const instanceid = request.params.instance;
        let resource = request.body as types.KubernetesResource;
 
        const projectRepository = repoFactory.getProjectRepository();

        resource = await projectRepository.updateResource(instanceid as string, resource);

        if( resource.type == types.ResourceType.deployment ) {
            const instance = await projectRepository.findInstance(instanceid);
            if( instance ) {
                var state: types.StateType;
                var status: types.StatusType;

                if( resource.status === types.StatusType.running ) {
                    status = types.StatusType.running
                    state = types.StateType.available;
                } else if (resource.status === types.StatusType.deleted ) {
                    status = types.StatusType.deleted;
                    state = types.StateType.available;
                } else if ( resource.status === types.StatusType.pending ) {
                    status = types.StatusType.pending;
                    state = types.StateType.available;
                } else {
                    status = types.StatusType.failed;
                    state = types.StateType.available;
                }
                await projectRepository.updateInstanceState(instanceid, status, state);
            }
        }

        response.status(HttpStatusCode.Ok).json( resource );

    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }

}

const addInstanceActivity = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-add-instance-activity');

    try {

        const instanceid = request.params.instance;
        const {op} = request.body;

        const projectRepository = repoFactory.getProjectRepository();

        logger.info(`add instance $(instanceid) activity operation - ${op}`);
        await projectRepository.addActivity(instanceid, op as types.ActivityType);
        const activities = await projectRepository.getActivities(instanceid);
        response.status(HttpStatusCode.Ok).json(activities);

    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });        
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const getInstanceActivities = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-get-instance-activities');

    try {

        const instanceid = request.params.instance;

        const projectRepository = repoFactory.getProjectRepository();

        const activities = await projectRepository.getActivities(instanceid);
        response.status(HttpStatusCode.Ok).json(activities);

    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });        
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const setInstancePermission = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-set-instance-permission');

    try {
        const instanceid = request.params.instance;
        const userid = request.params.user;
        const permission = request.body;

        const projectRepository = repoFactory.getProjectRepository();

        await projectRepository.setPermission(instanceid, userid, permission);
        const permissions = await projectRepository.getPermissions(instanceid);
        response.status(HttpStatusCode.Ok).json(permissions);
        
    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });        
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const removeInstancePermission = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-remove-instance-permission');

    try {
        const instanceid = request.params.instance;
        const userid = request.params.user;

        const projectRepository = repoFactory.getProjectRepository();

        await projectRepository.removePermission(instanceid, userid);
        const permissions = await projectRepository.getPermissions(instanceid);
        response.status(HttpStatusCode.Ok).json(permissions);
         
    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });        
    } finally {
        logger.info(`completed in ${getDuration()} ms`);
    }
}

const getInstancePermisions = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-get-instance-permissions');

    try {

        const instanceid = request.params.instance;
        const projectRepository = repoFactory.getProjectRepository();

        const permissions = await projectRepository.getPermissions(instanceid);
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

const getInstancePermision = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-get-instance-permission');

    try {
        const instanceid = request.params.instance;
        const pid = request.params.permission;

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

const getInstanceUserPermision = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('pctrl-get-instance-permission');

    try {
        const instanceid = request.params.instance;
        const user = request.params.user;

        const projectRepository = repoFactory.getProjectRepository();

        const permission = await projectRepository.getUserPermission(instanceid, user);
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


const instanceController = {
    findInstance,
    updateInstance,
    deleteInstance,
    purgeInstance,
    getInstanceResources,
    updateInstanceResource,
    addInstanceActivity,
    getInstanceActivities,
    setInstancePermission,
    removeInstancePermission,
    getInstancePermisions,
    getInstancePermision,
    getInstanceUserPermision
}

export default instanceController;
