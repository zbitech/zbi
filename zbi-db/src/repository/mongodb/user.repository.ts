import { Permission, RoleType, User, UserLimits, UserPermissions } from "../../types";
import { userModel } from "./schema";
import * as fn from "./fn";
import mongoose, { FilterQuery } from "mongoose";
import { getLogger } from "../../lib/logger";

const DEFAULT_LIMIT = {createResource: false, projects: 0, instances: 0};

const createUser = async (email: string, role: RoleType, limits: UserLimits = DEFAULT_LIMIT): Promise<User> => {
    let logger = getLogger('repo-create-user');
    try {
        const obj = { email, role, permissions: { projects: [], instances: [], limits } };
        const user = new userModel(obj);
        if (user) {
            await user.save();
            return fn.createUser(user);
        }

        throw new Error("failed to create user");
    } catch (error: any) {
        throw error;
    }
}

const getUser = async (userid: string): Promise<User | undefined> => {
    let logger = getLogger('repo-get-user');
    try {
        const user = await userModel.findById(userid);
        user?.toJSON()
        if (user) {
            return fn.createUser(user);
        }
        
        return undefined;
    } catch (error: any) {
        throw error;
    }
}

const findUsers = async (query: any): Promise<User[]> => {
    let logger = getLogger('repo-find-users');
    try {
        const filter = query as FilterQuery<any>;
        const users = await userModel.find(query);
        if (users) {
            return users.map((user: any) => {
                return fn.createUser(user);
            })
        }
        return [];
    } catch (err: any) {
        throw err;
    }
}

const findUser = async (query: any): Promise<User | undefined> => {
    let logger = getLogger('repo-find-user');
    try {
        const user = await userModel.findOne(query);
        if (user) {
            return fn.createUser(user);
        }
        return undefined;
    } catch (err: any) {
        throw err;
    }
}

const updateUserStatus = async (userid: string, active: boolean): Promise<User> => {
    let logger = getLogger('repo-update-user-status');
    try {
        const user = await userModel.findByIdAndUpdate(userid, { active });
        return fn.createUser(user);
    } catch (err: any) {
        throw err;
    }
}

const updateUserInfo = async (userid: string, info: any, provider: string): Promise<User> => {
    let logger = getLogger('repo-update-user-info');
    try {
        console.log(`searching for user with id = ${userid}`);
        const user = await userModel.findById(userid);
        console.log(`updating existing user ${userid} with  ${JSON.stringify(info)}`);
        if (user) {
            user.name = info.displayName;
            user.avatar = info.profileUrl;
            user.provider = { name: provider, userid: info.id };
            await user.save();
            return fn.createUser(user);
        }

        throw new Error("failed to find user");
    } catch (error: any) {
        throw error;
    }
}

const deleteUser = async (userid: string): Promise<void> => {
    let logger = getLogger('repo-delete-user');
    try {

    } catch (error: any) {
        throw error;
    }
}

const userMongoRepository = {
  createUser,
  getUser,
  findUsers,
  findUser,
  updateUserStatus,
  updateUserInfo,
  deleteUser
}

export default userMongoRepository