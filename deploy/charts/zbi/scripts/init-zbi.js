const fs = require("fs");
const pathExists = (path) => fs.existsSync(path);
const loadFile = (path) => fs.readFileSync(path);
const loadJSONFile = (path) => JSON.parse( fs.readFileSync(path) );

const data_path = "/etc/zbi/data";
const tmpl_path = "/etc/zbi/templates/blockchain/zcash"

const policy_path = `${data_path}/policies.json`;
const blockchain_path = `${data_path}/blockchains.json`;
const user_path = `${data_path}/users.json`;


const zdb = db.getSiblingDB("zbi");

if(pathExists(policy_path)) {
  console.log(`loading policies from ${policy_path}`);

  zdb.policies.deleteMany({});
  zdb.policies.insertOne(loadJSONFile(policy_path));
}

if(pathExists(user_path)) {
  console.log(`loading users from ${user_path}`);
  
  zdb.users.deleteMany({});
  zdb.users.insertMany(loadJSONFile(user_path));
}
  

if(pathExists(blockchain_path)) {

  let blockchains = loadJSONFile(blockchain_path);

  console.log(`loading blockchains from ${blockchain_path}`);

  const app_template_path = `${tmpl_path}/app_templates.tmpl`;
  const project_template_path = `${tmpl_path}/project_templates.tmpl`;
  const zcash_template_path = `${tmpl_path}/types/zcash_templates.tmpl`;
  const lwd_template_path = `${tmpl_path}/types/lwd_templates.tmpl`;
  
  blockchains = blockchains.map((blockchain) => {
    if(blockchain.name === "zcash") {
        if(pathExists(app_template_path)) {
            console.log(`setting app template from ${app_template_path}`);
            const app = loadFile(app_template_path); //.toString('base64');
            blockchain.templates["app"] = app;
        }

        if(pathExists(project_template_path)) {
            console.log(`setting project template from ${project_template_path}`);
            const project = loadFile(project_template_path); //.toString('base64');
            blockchain.templates["project"] = project;
        }

        for(var index=0; index<blockchain.nodes.length; index++) {
            if(blockchain.nodes[index].type === "zcash") {
                if(pathExists(zcash_template_path)) {
                    console.log(`setting zcash template from ${zcash_template_path}`);
                    const zcash = loadFile(zcash_template_path); //.toString('base64');
                    blockchain.templates["zcash"] = zcash;
                }        
            } else if(blockchain.nodes[index].type === "lwd") {
                if(pathExists(lwd_template_path)) {
                    console.log(`setting lwd template from ${lwd_template_path}`);
                    const lwd = loadFile(lwd_template_path); //.toString('base64');
                    blockchain.templates["lwd"] = lwd;
                }        
            }
        }
    }
    return blockchain;
  });
  
  zdb.blockchains.deleteMany({});
  zdb.blockchains.insertMany(blockchains);
}