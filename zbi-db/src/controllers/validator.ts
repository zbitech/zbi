import { HttpStatusCode } from "axios";
import { Request, Response, NextFunction } from "express";
import repoFactory from "../repository";
import Joi, { ObjectSchema } from "joi";
import { getLogger } from "../lib/logger";

const userSchema = Joi.object({
    body: Joi.object({
        email: Joi.string().required().email().label("email").messages({"any.required": "email is required"})
    }),
    query: {}, params: {}
});

const projectSchema = Joi.object({
    body: Joi.object({
        name: Joi.string().required().label("name"),
        blockchain: Joi.string().required().label("blockchain"),
        network: Joi.string().required().label("network"),
        description: Joi.string().allow("").label("description")
    }),
});

const instanceSchema = Joi.object({
    body: Joi.object({
        name: Joi.string().required().label("name"),
        type: Joi.string().required().label("type"),
        description: Joi.string().allow("").label("description")
    })
});


const validateSchema = (schema: ObjectSchema, data: any) => {

    const logger = getLogger('validate-schema');

    try {
        const result = schema.validate(data, {abortEarly: false});
        if (result.error && result.error.details) {
            const fieldErrors = result.error.details.map((error:any) => {
                return [error.context.label, error.message];
            })

            const fields = Object.fromEntries(fieldErrors);
            logger.info(`validation error: ${JSON.stringify(fields)}`);
            return {success: false, fields};
        }   
        return {success: true};
    } catch (err: any) {
        logger.error(`validation error: ${JSON.stringify(err)}`);
        throw err;
    }
}


const userEmailExists = async (request: Request, response: Response, next: NextFunction) => {
    try {
        const email = request.body.email;
        const userRepository = repoFactory.getUserRepository();

        const user = await userRepository.findUser({email});
        if(user) {
            response.status(HttpStatusCode.BadRequest).json({message: 'user already exists with email address'});
        } else {
            return next();
        }
    } catch (err: any) {
        response.sendStatus(HttpStatusCode.InternalServerError);
    }
}

const projectNameExists = async (request: Request, response: Response, next: NextFunction) => {
    try {
        const name = request.body.name;
        const projectRepository = repoFactory.getProjectRepository();

        const project = await projectRepository.findProjects({name});
        if(project.length>0) {
            response.status(HttpStatusCode.BadRequest).json({message: 'project already exists with name'});
        } else {
            return next();
        }
    } catch (err: any) {
        response.sendStatus(HttpStatusCode.InternalServerError);
    }
}

const instanceNameExists = async (request: Request, response: Response, next: NextFunction) => {
    try {
        const project = response.locals.project;
        const name = request.body.name;
        const projectRepository = repoFactory.getProjectRepository();

        const instance = await projectRepository.findInstances({project: project.id, name});
        if(instance.length>0) {
            response.status(HttpStatusCode.BadRequest).json({message: 'instance already exists with name'});
        } else {
            return next();
        }
    } catch (err: any) {
        response.sendStatus(HttpStatusCode.InternalServerError);
    }
}

const validateRequest = (message: string, schema: ObjectSchema, payload: any, response: Response, next: NextFunction) => {
    const logger = getLogger('validate-request');
    try {
        const result = validateSchema(schema, payload);
        if(!result.success) {
            logger.error(`validation error = ${JSON.stringify(result.fields)}`);
            response.status(HttpStatusCode.BadRequest).json({message, errors: result.fields})
            return;
        }
    } catch (err: any) {
        logger.error(`validation error: ${JSON.stringify(err)}`);
        response.status(HttpStatusCode.InternalServerError).json({message: err.message});
        return;
    }

    next();
}

const validateNewUser = async (request: Request, response: Response, next: NextFunction) => {
    const payload = {body: request.body};
    validateRequest("", userSchema, payload, response, next);
}

const validateNewProject = async (request: Request, response: Response, next: NextFunction) => {
    const logger = getLogger('validate-new-project');
    try {

        const payload = {body: request.body};
        const result = await validateSchema(projectSchema, payload);

        if(!result.success) {
            logger.error(`validation error = ${JSON.stringify(result.fields)}`);
            response.status(HttpStatusCode.BadRequest).json({message: "Project not created", errors: result.fields})
            return;
        }

    } catch (err: any) {
        logger.error(`validation error: ${JSON.stringify(err)}`);
        response.status(HttpStatusCode.InternalServerError).json({message: err.message});
        return;
    }

    next();

}

const validateNewInstance = async (request: Request, response: Response, next: NextFunction) => {
    const logger = getLogger('validate-new-instance');
    try {

        const blockchain = "zcash";
        const nodeType = request.body.type;

        const configRepository = repoFactory.getConfigRepository();

        const nodeInfo = await configRepository.getBlockchainNode(blockchain, nodeType);
        if( nodeInfo ) {
            const payload = {body: request.body};
            return validateRequest(`${blockchain} instance request could not be processed`, instanceSchema, payload, response, next);
        }

        response.status(HttpStatusCode.BadRequest).json({message: "invalid instance request"})

    } catch (err: any) {
        logger.error(`validation error: ${JSON.stringify(err)}`);
        response.status(HttpStatusCode.InternalServerError).json({message: err.message});
        return;
    }
}
 
const validateUpdateInstance = async (request: Request, response: Response, next: NextFunction) => {
    const logger = getLogger('validate-new-instance');
    try {

        const blockchain = "zcash";
        const nodeType = request.body.type;

        const configRepository = repoFactory.getConfigRepository();

        const nodeInfo = await configRepository.getBlockchainNode(blockchain, nodeType);
        if( nodeInfo ) {
            const payload = {body: request.body};
            return validateRequest(`${blockchain} instance request could not be processed`, instanceSchema, payload, response, next);
        }

        response.status(HttpStatusCode.BadRequest).json({message: "invalid instance request"})

    } catch (err: any) {
        logger.error(`validation error: ${JSON.stringify(err)}`);
        response.status(HttpStatusCode.InternalServerError).json({message: err.message});
        return;
    }
}


const validator = {
    userEmailExists,
    projectNameExists,
    instanceNameExists,
    validateNewProject,
    validateNewInstance,
    validateUpdateInstance
}

export default validator;