import * as fn from "./fn";
import { BlockchainInfo, NodeInfo, PolicyInfo } from "../../types";
import { blockchainModel, policyModel } from "./schema";
import { getLogger } from "../../lib/logger";
import { Types } from "mongoose";


const updatePolicy = async (policy: PolicyInfo): Promise<PolicyInfo> => {
    let logger = getLogger('repo-update-policy');
    try {
        let info = await policyModel.findOne({});
        if(!info) {
            info = new policyModel(policy);
        }
        await info.save();
        return fn.createPolicyInfo(info);
    } catch (err: any) {
        throw err;
    }
}

const getPolicy = async (): Promise<PolicyInfo|undefined> => {
    let logger = getLogger('repo-get-policy');
    try {
        const info = await policyModel.findOne({});
        if(info) {
            return fn.createPolicyInfo(info);
        }
        return undefined;
    } catch (err: any) {
        throw err;
    }
}

const createBlockchain = async (name: string, networks: string[]): Promise<BlockchainInfo> => {
    let logger = getLogger('repo-create-blockchain');
    try {
        const info = await blockchainModel.create({name, networks, nodes: []});
        return fn.createBlockchainInfo(info);
    } catch (err: any) {
        throw err;
    }
}

const getBlockchains = async (): Promise<BlockchainInfo[]> => {
    let logger = getLogger('repo-get-blockchains');
    try {
        const items = await blockchainModel.find({});
        return items.map((info: any) => fn.createBlockchainInfo(info));
    } catch (err: any) {
        throw err;
    }
}

const getBlockchain = async (name: string): Promise<BlockchainInfo> => {
    let logger = getLogger('repo-get-blockchain');
    try {
        const info = await blockchainModel.findOne({name});
        return fn.createBlockchainInfo(info);
    } catch (err: any) {
        throw err;
    }
}

const updateBlockchainNode = async (blockchain: string, node: NodeInfo): Promise<BlockchainInfo> => {
    let logger = getLogger('repo-update-blockchain-node');
    try {
        const date = new Date();
        const info = await blockchainModel.findOne({name: blockchain});
        if(info) {
            let found = false;
            for(var index=0; index<info.nodes.length; index++) {
                if(info.nodes[index].name === node.name) {
                    info.nodes[index].images = node.images;
                    info.nodes[index].endpoints = node.endpoints;
                    info.nodes[index].ports = node.ports;
                    info.nodes[index].settings = node.settings;
                    info.nodes[index].properties = node.properties;
                    info.nodes[index].updatedAt = new Date();
                                    
                    found = true;
                }
            }
            if(!found) info.nodes.push({...node, createdAt: date, updatedAt: date});
            await info.save();
            return fn.createBlockchainInfo(info);
        }

        throw new Error("blockchain not found");
    } catch (err: any) {
        throw err;        
    }
}

const removeBlockchainNode = async (blockchain: string, node: string): Promise<BlockchainInfo> => {
    let logger = getLogger('repo-remove-blockchain-node');
    try {
        const info = await blockchainModel.findOne({name: blockchain});
        if(info) {

            info.nodes = info.nodes.filter((n: any) => n.name !== node);
            // info.nodes = new Types.DocumentArray(info.nodes.filter((n: any) => n.name !== node));
            await info.save();
            return fn.createBlockchainInfo(info);
        }

        throw new Error("blockchain not found");
    } catch (err: any) {
        throw err;        
    }
}


const getBlockchainNode = async (blockchain: string, node: string): Promise<NodeInfo|undefined> => {
    let logger = getLogger('repo-get-blockchain-node');
    try {
        const info = await blockchainModel.findOne({name: blockchain});
        if(info) {
//            logger.debug(`nodes - ${JSON.stringify(info.nodes)}`);
            const nodes = info.nodes.filter((n: any) => n.type === node);
//            logger.debug(`nodes - ${JSON.stringify(nodes)}`);
            if(nodes.length > 0) return fn.createBlockchainNodeInfo(nodes[0]);
        } else {
            console.log(`could not find blockchain ${blockchain}`);
        }

        return undefined;
    } catch (err: any) {
        throw err;
    }
}

const getBlockchainNodeTemplate = async (blockchain: string, node: string): Promise<String|undefined> => {
    let logger = getLogger('repo-get-blockchain-node');
    try {
        const info = await blockchainModel.findOne({name: blockchain}, {templates: 1});
        if(info) {
            if(info.templates) {
                
            }
            // const template = info?.templates.get(node);
            // return template?.toString();
            return "Hello";
        }

        throw new Error("node not found");
    } catch (err: any) {
        throw err;
    }

}

const configMongoRepository = {
    updatePolicy,
    getPolicy,
    createBlockchain,
    getBlockchains,
    getBlockchain,
    updateBlockchainNode,
    removeBlockchainNode,
    getBlockchainNode,
    getBlockchainNodeTemplate
}

export default configMongoRepository;