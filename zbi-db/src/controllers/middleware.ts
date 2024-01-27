import { HttpStatusCode } from "axios";
import { Request, Response, NextFunction } from "express";
import * as types from "../types";
import repoFactory from "../repository";
import context from '../lib/context';
import { getLogger, mainLogger as logger } from "../lib/logger";
import {v4 as uuidv4} from 'uuid';
import { Constants } from '../constants';

const INTERNAL_CLIENT_SECRET = process.env.ZBI_INTERNAL_CLIENT_SECRET || "zbi-internal-client"; 

const initRequest = (request: Request, response: Response, next: NextFunction) => {

    const store = new Map();  

    let requestid = uuidv4();
    const userid = request.headers[Constants.USER_ID] ? request.headers[Constants.USER_ID] : "zbi";

    return context.run(store, () => {
        const child = logger.child({ requestid, userid });

        store.set(Constants.REQUEST_ID, requestid);
        store.set(Constants.USER_ID, userid);
        store.set(Constants.START, process.hrtime());
        store.set('logger', child);
        
        next();
    });
};

// const checkProjectPermissions = (type: types.PermissionType) => {

//     return async (request: Request, response: Response, next: NextFunction) => {
//         const logger = getLogger("check-project-permissions");

//         logger.info("check project permissions");
//         const projectRepository = repoFactory.getUserRepository();

//         const project = response.locals.project as types.Project;
//         const userid = request.headers[Constants.USER_ID] as string;

//         logger.debug(`current user - ${JSON.stringify(userid)}`);
//         logger.debug(`current project - ${JSON.stringify(project)}`);

//         if(project.owner === userid as string) {
//             logger.debug(`allowing project owner`);
//             return next();
//         } else {
//             logger.debug(`user *${userid}* is not owner *${project.owner}*`);
//         }

//         const permission = await projectRepository.getProjectPermission(userid, project.id as string);
//         let permit = false;
//         if( permission ) {
//             if (type === types.PermissionType.read) {
//                 permit = permission.read;
//             } else if (type === types.PermissionType.update) {
//                 permit = permission.update;
//             } else if (type === types.PermissionType.delete) {
//                 permit = permission.delete;
//             }
//         }

//         if(permit) return next();
//         response.sendStatus(HttpStatusCode.Forbidden);
//     }
// }

// const checkInstancePermissions = (type: types.PermissionType) => {
//     return async (request: Request, response: Response, next: NextFunction) => {

//         const userRepository = repoFactory.getUserRepository();
//         const projectRepository = repoFactory.getProjectRepository();

//         const instance = response.locals.instance as types.Instance;
//         const userid = request.headers[Constants.USER_ID] as string;

//         if(!response.locals.project) {
//             response.locals.project = await projectRepository.findProject(instance.project);
//         }

//         if(response.locals.project.owner === userid) return next();

//         const permission = await userRepository.getInstancePermission(userid, instance.id as string);
//         let permit = false;
//         if( permission ) {
//             if (type === types.PermissionType.read) {
//                 permit = permission.read;
//             } else if (type === types.PermissionType.update) {
//                 permit = permission.update;
//             } else if (type === types.PermissionType.delete) {
//                 permit = permission.delete;
//             } else if (type === types.PermissionType.operate) {
//                 permit = permission.operate ? permission.operate : false;
//             } else if (type === types.PermissionType.access) {
//                 permit = permission.access ? permission.access : false;
//             }
//         }

//         if(permit) return next();
//         response.sendStatus(HttpStatusCode.Forbidden);
//     }
// }

// const checkCreateProjectPermission = async (request: Request, response: Response, next: NextFunction) => {

//     try {
//         const userid = request.headers[Constants.USER_ID];

//         const userRepository = repoFactory.getUserRepository();
//         const projectRepository = repoFactory.getProjectRepository();

//         const permissions = await userRepository.findPermissions({userid});
//         if( permissions.limits.createResource ) {

//             const projects = await projectRepository.findProjects({owner: userid});
//             if(projects.length < permissions.limits.projects) {
//                 return next();
//             }

//             response.status(HttpStatusCode.TooManyRequests).json({message: 'project quota exceeded'});
//         } else {
//             response.status(HttpStatusCode.Forbidden);
//         }
        
//     } catch (err: any) {
//         response.status(HttpStatusCode.InternalServerError);
//     }
// }

// const checkCreateInstancePermission = async (request: Request, response: Response, next: NextFunction) => {

//     try {

//         const userid = request.headers[Constants.USER_ID];
//         const userRepository = repoFactory.getUserRepository();
//         const projectRepository = repoFactory.getProjectRepository();

//         const permissions = await userRepository.findPermissions({userid});
//         if( permissions.limits.createResource ) {

//             const instances = await projectRepository.findInstances({owner: userid});
//             if(instances.length < permissions.limits.instances) {
//                 return next();
//             }

//             response.status(HttpStatusCode.TooManyRequests).json({message: 'instance quota exceeded'});
//         } else {
//             response.status(HttpStatusCode.Forbidden);
//         }
        
//     } catch (err: any) {
//         response.status(HttpStatusCode.InternalServerError);
//     }
// }

const validateProject = async (request: Request, response: Response, next: NextFunction) => {
    try {
        const projectRepository = repoFactory.getProjectRepository();
        const valid = await projectRepository.checkProjectId(request.params.project);
        if(valid) return next();
        response.status(HttpStatusCode.NotFound).json({message: 'project not found'});
    } catch (err: any) {
        response.status(HttpStatusCode.BadRequest).json({message: 'valid project is required'});   
    }
}

const validateInstance = async (request: Request, response: Response, next: NextFunction) => {
    const logger = getLogger("set-instance");
    try {
        const projectRepository = repoFactory.getProjectRepository();
        const valid = await projectRepository.checkInstanceId(request.params.instance);
        if(valid) return next();
        response.status(HttpStatusCode.BadRequest).json({message: 'instance not found'});        
    } catch (err: any) {
        logger.info(`failed to setInstance - ${err}`);
        response.status(HttpStatusCode.BadRequest).json({message: 'valid instance is required'});        
    } finally {
        logger.info(`complete`);
    }
}

const logRequest =async (request: Request, response: Response, next: NextFunction) => {
    const logger = getLogger("set-instance");
    next();

}

const middleware = {
    initRequest,
    validateProject,
    validateInstance
}

export default middleware;
