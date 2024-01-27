import {Router} from "express";
import {instanceController, validator, middleware} from "../controllers";
import * as types from "../types";

const instanceRoutes = Router();

instanceRoutes.get("/:instance", middleware.validateInstance, instanceController.findInstance)
instanceRoutes.put("/:instance", middleware.validateInstance, instanceController.updateInstance)
instanceRoutes.delete("/:instance", middleware.validateInstance, instanceController.deleteInstance)
instanceRoutes.purge("/:instance", middleware.validateInstance, instanceController.purgeInstance);

instanceRoutes.get("/:instance/resources", middleware.validateInstance, instanceController.getInstanceResources)
instanceRoutes.post("/:instance/resources", middleware.validateInstance, instanceController.updateInstanceResource)

instanceRoutes.post("/:instance/activities", middleware.validateInstance, instanceController.addInstanceActivity);
instanceRoutes.get("/:instance/activities", middleware.validateInstance, instanceController.getInstanceActivities);

instanceRoutes.get("/:instance/permissions", middleware.validateInstance, instanceController.getInstancePermisions);
instanceRoutes.get("/:instance/permissions/:user", middleware.validateInstance, instanceController.getInstanceUserPermision);
instanceRoutes.post("/:instance/permissions/:user", middleware.validateInstance, instanceController.setInstancePermission);
instanceRoutes.delete("/:instance/permissions/:user", middleware.validateInstance, instanceController.removeInstancePermission);

export default instanceRoutes;