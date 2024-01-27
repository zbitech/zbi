
import * as mongo from "./mongodb"

const MONGO_DB_CONNECTION_STRING = process.env.MONGO_DB_CONNECTION_STRING || "mongodb://localhost/zbiRepo";

class RepositoryFactory {

    constructor() {
        
    }

    async connect() {
        await mongo.database.connect(MONGO_DB_CONNECTION_STRING);
    }

    generateId() {
        return mongo.database.generateId();
    }

    getUserRepository() {
        return mongo.userRepository;
    }

    getProjectRepository() {
        return mongo.projectRepository;
    }

    getConfigRepository() {
        return mongo.configRepository;
    }

    async close() {
        await mongo.database.close();
    }
}

export default new RepositoryFactory();