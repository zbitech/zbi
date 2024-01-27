import { createLogger, format, transports, addColors } from "winston";
import context from "./context";
import {v4 as uuidv4} from 'uuid';
import morgan, {StreamOptions} from "morgan";

const { combine, label, timestamp, printf, json, colorize } = format;

addColors({error: 'red', warn: 'yellow', info: 'green', http: 'magenta', debug: 'white'});
const console = new transports.Console();

const log_transports = [];

if(process.env.NODE_ENV === "production" || process.env.NODE_ENV === "dev") {
    log_transports.push(console);
} else {
    log_transports.push(new transports.File({
        filename: "logs/info.log"
    }));
}

const logFormat = printf(({level, message, label: logLabel, timestamp: logTimeStamp, requestid, userid, service}) => {
    return `${level} | ${logTimeStamp} | ${requestid} | ${userid} | ${service} | ${message}`;
});

export const defaultLogger = createLogger({
    level: process.env.LOG_LEVEL || 'debug',
    format: combine(label({label: process.env.NODE_ENV}), timestamp(), logFormat, colorize({all: true})),
    transports: log_transports,
});

export const mainLogger = defaultLogger.child({requestid: uuidv4(), userid: "zbi", service: "main"});

export const morganStream: StreamOptions = {write: (message: string) => mainLogger.http(message)};


export const getLogger = (service: string) => {
    let x = context.getStore()?.get('logger');

    let logger = context.getStore()?.get('logger') || mainLogger;
    return logger.child({service});
}

export const getDuration = () => {
    let start = context.getStore()?.get("start");
    const NS_PER_SEC = 1e9
    const NS_TO_MS = 1e6
    const diff = process.hrtime(start)

    return ((diff[0] * NS_PER_SEC + diff[1]) / NS_TO_MS).toLocaleString();
}