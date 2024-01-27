import { Request, Response } from "express";
import { HttpStatusCode } from "axios";
import repoFactory from "../repository";
import { handleError } from "../lib/errors";
import { getLogger } from "../lib/logger";


const getPolicy = async (request: Request, response: Response) => {
    let logger = getLogger('get-policy');
    try {
        
        const config = repoFactory.getConfigRepository();
        const policy = await config.getPolicy();
        if(policy) {
            response.status(HttpStatusCode.Ok).json(policy);
        } else {
            response.status(HttpStatusCode.Ok).json();
        }

    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });        
    }
}

const updatePolicy = async (request: Request, response: Response) => {
    let logger = getLogger('update-policy');
    try {
        
        let policy = request.body;
        const config = repoFactory.getConfigRepository();
        policy = await config.updatePolicy(policy);
        if(policy) {
            response.status(HttpStatusCode.Ok).json(policy);
        } else {
            response.status(HttpStatusCode.Ok).json();
        }

    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });        
    }
}

const createBlockchain = async (request: Request, response: Response) => {
    let logger = getLogger('get-blockchain');
    try {
        const {name, networks} = request.body;

        const config = repoFactory.getConfigRepository();
        const blockchain = await config.createBlockchain(name, networks);
        response.status(HttpStatusCode.Ok).json(blockchain);
    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });        
    }
}

const getBlockchains = async (request: Request, response: Response) => {
    let logger = getLogger('get-blockchains');
    try {
        
        const config = repoFactory.getConfigRepository();
        const blockchains = await config.getBlockchains();
        response.status(HttpStatusCode.Ok).json(blockchains)
    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });        
    }
}

const getBlockchain = async (request: Request, response: Response) => {
    let logger = getLogger('get-blockchain');
    try {
        const name = request.params.blockchain;

        const config = repoFactory.getConfigRepository();
        const blockchain = await config.getBlockchain(name);
        response.status(HttpStatusCode.Ok).json(blockchain);
    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });        
    }
}

const getBlockchainNode = async (request: Request, response: Response) => {
    let logger = getLogger('get-blockchain-node');
    try {
        const blockchain = request.params.blockchain;
        const node = request.params.node;

        logger.debug(`searching for node ${node} in blockchain ${blockchain}`);
        const config = repoFactory.getConfigRepository();
        const blockchainNode = await config.getBlockchainNode(blockchain, node);
        response.status(HttpStatusCode.Ok).json(blockchainNode);
    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });        
    }
}

const updateBlockchainNode = async (request: Request, response: Response) => {
    let logger = getLogger('update-blockchain-node');
    try {
        let blockchain = request.params.blockchain;
        let node = request.body;        
        const config = repoFactory.getConfigRepository();
        node = await config.updateBlockchainNode(blockchain, node);
        response.status(HttpStatusCode.Ok).json({node});
    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });        
    }
}

const removeBlockchainNode = async (request: Request, response: Response) => {
    let logger = getLogger('remove-blockchain-node');
    try {
        
        let blockchain = request.params.blockchain;
        let node = request.params.node;

        const config = repoFactory.getConfigRepository();
        const result = await config.removeBlockchainNode(blockchain, node);
        response.status(HttpStatusCode.Ok).json(result);

    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });        
    }
}

const getBlockchainNodeTemplate = async(request: Request, response: Response) => {
    let logger = getLogger('get-blockchain-node-template');

    try {
        let blockchain = request.params.blockchain;
        let node = request.params.node;

        logger.info(`retrieving ${node} template for ${blockchain} node`);
        const config = repoFactory.getConfigRepository();
        const template = config.getBlockchainNodeTemplate(blockchain, node)

        response.set('Content-Type', 'text/plain');
        response.status(HttpStatusCode.Ok);
        response.send(template);

    } catch (err: any) {
        const result = handleError(err);
        logger.error(`response - ${JSON.stringify(result)}`);
        response.status(result.code).json({ message: result.message });        
    }
}

const configController = {
    getPolicy,
    updatePolicy,
    createBlockchain,
    getBlockchains,
    getBlockchain,
    getBlockchainNode,
    updateBlockchainNode,
    removeBlockchainNode,
    getBlockchainNodeTemplate
}

export default configController;