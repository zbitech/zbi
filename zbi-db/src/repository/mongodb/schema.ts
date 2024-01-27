import { Schema, model, Document, Types, Model } from "mongoose";
import { ResourceType, RoleType } from "../../types";

const userSchema = new Schema({
    active: {type: Boolean, default: false},
    role: {type: String, enum: [RoleType.admin, RoleType.user]},
    password: {type: String},
    email: {type: String},
    name: {type: String},
    avatar: {type: String},
    provider: {
        name: {type: String},
        userid: {type: String}
    },
}, {timestamps: true, strict: false});

const activitySchema = new Schema({
    type: {type: String, enum: ['project', 'instance']},
    operation: {type: String, required: true, immutable: true},
    object: {type: Schema.Types.ObjectId},
    completed: {type: Boolean},
    success: {type: Boolean},
}, {
    timestamps: true,
    index: {
        keyPattern: {createdAt: 1},
        expireAtSeconds: 604800
    }
});

const permissionSchema = new Schema({
    type: {type: String, enum: ['project', 'instance']},
    object: {type: Schema.Types.ObjectId},
    user: {type: Schema.Types.ObjectId, ref: "user", required: true, immutable: true},
    read: {type: Boolean},
    update: {type: Boolean},
    delete: {type: Boolean},
    operate: {type: Boolean},
    access: {type: Boolean},
},{
    timestamps: true
});

const resourceSchema = new Schema({
    object: {type: Schema.Types.ObjectId, required: true, immutable: true},
    type: {type: String, enum: [ResourceType.namespace, ResourceType.configmap,
        ResourceType.secret, ResourceType.persistentvolumeclaim,
        ResourceType.deployment, ResourceType.service,
        ResourceType.httpproxy, ResourceType.volumesnapshot,
        ResourceType.snapshotschedule,
    ]},
    name: {type: String},
    status: {type: String},
    properties: {type: Schema.Types.Mixed},
}, {
    timestamps: true
});

const projectSchema = new Schema({
    name: {type: String, required: true, immutable: true},
    owner: {type: Schema.Types.ObjectId, ref: "user", immutable: true, required: true},
    blockchain: {type: String, required: true, immutable: true, enum: ['zcash']},
    network: {type: String, required: true, immutable: true, enum: ['testnet', 'regnet', 'mainnet']},
    status: {type: String, default: 'new', /*enum: ['new', 'pending', 'active', 'inactive']*/},
    description: {type: String},
    state: {type: String}

}, {timestamps: true});

const instanceSchema = new Schema({
    name: {type: String, required: true, immutable: true},
    type: {type: String, required: true, immutable: true, enum:['zcash', 'lwd', 'zebra']},
    description: {type: String},
    status: {type: String, default: 'new'},
    project: {type: Schema.Types.ObjectId, ref: "project", immutable: true},
    request: {
        cpu: {type: String},
        memory: {type: String},
        peers: {type: [String]},
        properties: {type: Schema.Types.Mixed},
        volume: {
            type: {type: String},
            size: {type: String},
            source: {
                type: {type: String},
                ref: {type: String}
            }
        }
    },
    state: {type: String}
}, {timestamps: true});

const policySchema = new Schema({
    storageClass: {type: String},
    snapshotClass: {type: String},
    domainName: {type: String},
    certificateName: {type: String},
    serviceAccount: {type: String},
    envoyConfig: {
        image: {type: String},
        command: [{type: String}],
        timeout: {type: Number},
        accessAuthorization: {type: Boolean},
        authServerURL: {type: String},
        authServerPort: {type: Number},
        authenticationEnabled: {type: Boolean}
    },
    informerResync: {type: Number},
    enableMonitor: {type: Boolean},
    requireAuthentication: {type: Boolean},
    request: {
        cpu: {type: String},
        memory: {type: String},
        storage: {type: String}
    },
});

const blockchainSchema = new Schema({
    name: {type: String},
    networks: [{type: String}],
    templates: {type: Schema.Types.Map, of: Buffer},
    nodes: [{
        name: {type: String},
        type: {type: String},
        images: [{
            name: {type: String},
            version: {type: String},
            url: {type: String}
        }],
        endpoints: {type: Schema.Types.Mixed},
        ports: {type: Schema.Types.Map, of: Number},
        settings: {type: Schema.Types.Mixed},
        properties: {type: Schema.Types.Mixed},
        createdAt: {type: Date},        
        updatedAt: {type: Date}
    }]
}, {timestamps: true });

const userModel = model("users", userSchema);
const projectModel = model("project", projectSchema);
const instanceModel = model("instance", instanceSchema);
const policyModel = model("policy", policySchema);
const blockchainModel = model("blockchain", blockchainSchema);
const activityModel = model("activity", activitySchema);
const permissionModel = model("permission", permissionSchema);
const resourceModel = model("resource", resourceSchema);

export {
    userModel, projectModel, instanceModel, policyModel, blockchainModel,
    activityModel, permissionModel, resourceModel
}