
import * as dotenv from "dotenv";
import express from "express";
import cors from "cors";
import helmet from "helmet";
import morgan from "morgan";

console.log("initializing application ...");
//dotenv.config();

import {mainLogger as logger, morganStream} from "./lib/logger";

import { middleware } from "./controllers";

const app = express();

console.log("initializing logger ...");
app.use(morgan(":method :url :status :res[content-length] - :response-time ms", {stream: morganStream}));

app.use(helmet());
app.use(cors({origin: "*", credentials: true}));
app.use(express.json());
app.use(express.urlencoded({extended: true}));

app.use(middleware.initRequest);

export default app;