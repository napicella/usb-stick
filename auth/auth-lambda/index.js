// https://medium.com/@prasadjay/amazon-cognito-user-pools-in-nodejs-as-fast-as-possible-22d586c5c8ec
// https://stackoverflow.com/questions/48537867/set-cognito-verification-type-to-link-in-cloudformation
const AmazonCognitoIdentity = require('amazon-cognito-identity-js');
const CognitoUserPool = AmazonCognitoIdentity.CognitoUserPool;
const AWS = require('aws-sdk');
const request = require('request');
const jwkToPem = require('jwk-to-pem');
const jwt = require('jsonwebtoken');
global.fetch = require('node-fetch');
const router = require('aws-lambda-router');

const poolData = {
    UserPoolId : process.env.USER_POOL_ID,
    ClientId : process.env.USER_POOL_CLIENT_ID
};

const userPool = new AmazonCognitoIdentity.CognitoUserPool(poolData);

exports.handler = router.handler({
    proxyIntegration: {
        debug: true,
        routes: [
            {
                path: '/auth/login',
                method: 'POST',
                action: (request, context) => {
                    return login();
                }
            },
            {
                path: '/auth/register',
                method: 'POST',
                action: (request, context) => {
                    return registerUser(request.body);
                }
            }
            ]}
});

function registerUser({email, password}){
    var attributeList = [];
    attributeList.push(new AmazonCognitoIdentity.CognitoUserAttribute({Name:"name",Value:email}));
    attributeList.push(new AmazonCognitoIdentity.CognitoUserAttribute({Name:"email",Value:email}));

    return new Promise((resolve, reject) => {
        userPool.signUp(email, password, attributeList, null, function(err, result) {
            if (err) {
                console.log(err);
                return reject(err);
            }
            console.log(result);
            return resolve(result);
        });
    })
}