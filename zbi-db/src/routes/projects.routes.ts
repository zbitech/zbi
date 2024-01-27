import {Router} from "express";
import {projectController, validator, middleware} from "../controllers";
import * as types from "../types";

const projectRoutes = Router();

projectRoutes.get("/", projectController.findProjects)
projectRoutes.post("/", validator.validateNewProject, validator.projectNameExists, projectController.createProject);
projectRoutes.get("/:project", middleware.validateProject, projectController.findProject)
projectRoutes.put("/:project", middleware.validateProject, projectController.updateProject)
projectRoutes.delete("/:project", middleware.validateProject, projectController.deleteProject)
projectRoutes.purge("/:project", middleware.validateProject, projectController.purgeProject);

projectRoutes.get("/:project/instances", middleware.validateProject, projectController.findInstances)
projectRoutes.post("/:project/instances", middleware.validateProject, validator.validateNewInstance, validator.instanceNameExists, projectController.createInstance);

projectRoutes.get("/:project/snapshots", middleware.validateProject, projectController.getProjectSnapshots);

projectRoutes.get("/:project/resources", middleware.validateProject, projectController.getProjectResources)
projectRoutes.put("/:project/resources", middleware.validateProject, projectController.updateProjectResource)

projectRoutes.post("/:project/activities", middleware.validateProject, projectController.addProjectActivity);
projectRoutes.get("/:project/activities", middleware.validateProject, projectController.getProjectActivities);

projectRoutes.get("/:project/permissions", middleware.validateProject, projectController.getProjectPermisions);
projectRoutes.get("/:project/permissions/:user", middleware.validateProject, projectController.getProjectUserPermision);
projectRoutes.post("/:project/permissions/:user", middleware.validateProject, projectController.setProjectPermission);
projectRoutes.delete("/:project/permissions/:user", middleware.validateProject, projectController.removeProjectPermission);

export default projectRoutes;