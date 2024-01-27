import {Router, Request, Response} from "express";
import { configController } from "../controllers";
import * as types from "../types";

const configRoutes = Router();

configRoutes.get("/policy", configController.getPolicy);
configRoutes.post("/policy", configController.updatePolicy);
configRoutes.get("/blockchains", configController.getBlockchains);
configRoutes.post("/blockchains", configController.createBlockchain);
configRoutes.get("/blockchains/:blockchain", configController.getBlockchain);
configRoutes.post("/blockchains/:blockchain", configController.updateBlockchainNode);
configRoutes.get("/blockchains/:blockchain/:node", configController.getBlockchainNode);
configRoutes.delete("/blockchains/:blockchain/:node", configController.removeBlockchainNode);

configRoutes.get("/blockchains/:blockchain/:node/template", configController.getBlockchainNodeTemplate);

export default configRoutes;