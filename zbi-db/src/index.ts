import app from "./app";
import {mainLogger as logger, morganStream} from "./lib/logger";
import repositoryFactory from "./repository";
import routes from "./routes";

const PORT: number = parseInt(process.env.PORT || "4000" as string, 10);

console.log(`running app ...`);

repositoryFactory.connect().then(async () => {

    routes(app);

    app.listen(PORT, async () => {
        logger.info(`Running Node.js version ${process.version}`);
        logger.info(`App environment: ${process.env.NODE_ENV}`);
    
        logger.info(`App is running on port ${PORT}`, `{request: '55'}`);
    });        
}).catch(error => {
   logger.error(`failed to initialize database: ${error}`);
});
