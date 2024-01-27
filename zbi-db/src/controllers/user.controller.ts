import { Request, Response } from 'express';
import repoFactory from "../repository";

import { getLogger } from "../lib/logger";
import { HttpStatusCode } from 'axios';
import { handleError } from '../lib/errors';

const createUser = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger('uc-create-user');
    try {
        const [email, role] = request.body;
        let userRepo = repoFactory.getUserRepository();
        const user = await userRepo.createUser("", email, role);
        response.json({ user });
    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });
    }
}

const findUsers = async (request: Request, response: Response): Promise<void> => {

    let logger = getLogger('uc-find-users');
    try {
        const userRepo = repoFactory.getUserRepository();

        const name = request.query.name;
        const value = request.query.value as string;
        const size = request.query.size ? parseInt(request.query.size as string) : 10;
        const page = request.query.page ? parseInt(request.query.page as string) : 1;
        const param = name ? { name: value } : {};

        const users = await userRepo.findUsers(param);
        response.json({ users, size, page });
    } catch (err: any) {
        logger.error(err);
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });
    }

}

const findUser = async (request: Request, response: Response): Promise<void> => {

    let logger = getLogger('find-user');
    try {
        const userRepo = repoFactory.getUserRepository();

        logger.info(`finding user ... ${JSON.stringify(request.body)}`);
        const param = request.body.param;

        const user = userRepo.findUser(param);
        response.json(user);

    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });
    }
}

const deactivateUser = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger(`uc-deactivate-user`);

    try {
        const userid = response.locals.user;
        const userRepo = repoFactory.getUserRepository();

        const user = await userRepo.updateUserStatus(userid, false);
        response.status(HttpStatusCode.Ok).json(user);

    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });
    }

}

const reactivateUser = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger(`uc-reactivate-user`);

    try {
        const userid = response.locals.user.userid;
        const userRepo = repoFactory.getUserRepository();

        const user = await userRepo.updateUserStatus(userid, true);
        response.status(HttpStatusCode.Ok).json(user);

    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });
    }

}

const deleteUser = async (request: Request, response: Response): Promise<void> => {
    let logger = getLogger(`uc-delete-user`);

    try {
        const userid = response.locals.user.userid;
        const userRepo = repoFactory.getUserRepository();

        await userRepo.deleteUser(userid);
        response.sendStatus(HttpStatusCode.Ok)

    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });
    }
}



const userController = {
    createUser,
    findUsers,
    findUser,
    deactivateUser,
    reactivateUser,
    deleteUser
}

export default userController