import {Express} from "express";
import configRoutes from "./config.routes";
import userRoutes from "./users.routes";
import projectRoutes from "./projects.routes";
import instanceRoutes from "./instances.routes";

const routes = (app: Express) => {

    app.use("/api/config", configRoutes);
    app.use("/api/users", userRoutes);
    app.use("/api/projects", projectRoutes);
    app.use("/api/instances", instanceRoutes);
}

export default routes;