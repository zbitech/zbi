import mongoose from "mongoose";
import projectMongoRepository from "./project.repository";
import userMongoRepository from "./user.repository";
import configMongoRepository from "./config.repository";
import * as schema from "./schema";
import * as fn from "./fn";
import { getLogger } from "../../lib/logger";

const MONGO_DB_CONNECTION_STRING = process.env.MONGO_DB_CONNECTION_STRING || "mongodb://localhost/zbiRepo";
const MONGOOSE_SERVER_SELECTION_TIMEOUT = parseInt(process.env.MONGO_SERVER_SELECTION_TIMEOUT || "5000");
const MONGOOSE_USE_NEW_URL_PARSER = Boolean(process.env.MONGO_USE_NEW_URL_PARSER || "true");

const MONGOOSE_OPTIONS = {
    useNewUrlParser: MONGOOSE_USE_NEW_URL_PARSER,
    serverSelectionTimeoutMS: MONGOOSE_SERVER_SELECTION_TIMEOUT
}

const database = {
    connect: async(url: string) => {
        const logger = getLogger("database");
        try {
           logger.info(`connecting to: ${url}`);
            await mongoose.connect(url, MONGOOSE_OPTIONS)            
            const db = mongoose.connection;
            db.on("error", () => {
                logger.info("could not connect");
            });
            db.once("open", () => {
                logger.info("successfully connected to database");
            });
            logger.info("connected to database");
        } catch (err: any) {
            throw err;
        }
    },

    close: async() => {
        try {
            await mongoose.connection.close();            
        } catch (err: any) {
            throw err;
        }
    },

    generateId: () => {
        return fn.generateId();
    }
}


const projectRepository = projectMongoRepository;
const userRepository = userMongoRepository;
const configRepository = configMongoRepository;

export {
    database, projectRepository, userRepository, configRepository, 
    schema, fn
}