import {Router} from "express";
import {userController} from "../controllers";
import * as types from "../types";

const userRoutes = Router();

userRoutes.post("/", userController.createUser);
userRoutes.get("/", userController.findUsers);
userRoutes.delete("/:userid", userController.deleteUser);
userRoutes.get("/:userid", userController.findUser);
userRoutes.put("/:userid/reactivate", userController.reactivateUser);
userRoutes.put("/:userid/deactivate", userController.deactivateUser);

export default userRoutes;