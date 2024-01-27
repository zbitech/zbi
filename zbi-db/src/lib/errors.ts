import { HttpStatusCode } from "axios";
// import { Http } from "winston/lib/winston/transports";

// export enum ItemType {
//     user = 0,
//     team = 1,
//     membership = 2,
//     project = 3,
//     instance = 4,
//     resource = 5
// }

// export enum ServiceType {
//     database = 0,
//     controller = 1,
//     file = 2,
//     network = 3,
// }

// export class AppResponse {
//     public code: number;
//     public message: string;

//     constructor(code: number, message: string) {
//         this.code = code;
//         this.message = message;
//     }
// }

export class AppError extends Error {
    status: number;

    constructor(message: string, status: number = HttpStatusCode.InternalServerError) {
        super(message);

        Error.captureStackTrace(this, this.constructor);

        this.name = this.constructor.name;
        this.status = status;
    }
}

export class ItemNotFoundError extends AppError {
    constructor(message: string) {
        super(message, HttpStatusCode.NotFound);
    }
}

export class ItemAlreadyExistsError extends AppError {
    constructor(message: string) {
        super(message, HttpStatusCode.BadRequest);
    }
}

export function handleError(err: Error) {
    if(err instanceof ItemNotFoundError) {
        return {code: HttpStatusCode.NotFound, message: err.message};
    } 
    
    else if(err instanceof ItemAlreadyExistsError) {
        return {code: HttpStatusCode.BadRequest, message: err.message};
    } 
    

    // else if(err instanceof ServiceError) {
    //     return {code: HttpStatusCode.InternalServerError, message: err.message};
    // }

    return {code: HttpStatusCode.InternalServerError, message: "Internal ServerError"};
} 


// export class ApplicationError extends Error {
//     public code: number = 0;

//     constructor(code: number, message: string, ...args: any) {
//         super(...args);
//         this.code = code;
//         this.message = message;
//     }
// }

// export class ServiceError extends Error {
//     public service: ServiceType;

//     constructor(service: ServiceType, message: string) {
//         super(message);
//         this.service = service;        
//     }

// }


