const AWS = require('aws-sdk'); // eslint-disable-line import/no-extraneous-dependencies

// https://docs.aws.amazon.com/cli/latest/reference/apigateway/update-stage.html
// aws apigateway update-stage --rest-api-id rrtvx1muac --stage-name 'Prod' --patch-operations op=replace,path=/~1*/*/throttling/rateLimit,value=1
const apigateway = new AWS.APIGateway({apiVersion: '2015-07-09', region: process.env.AWS_REGION});

const handler = (event, context, callback) => {

    var params = {
        restApiId: 'rrtvx1muac',
        stageName: 'Prod',
        patchOperations: [{
            op: 'replace',
            path: '/~1*/*/throttling/rateLimit',
            value: '0'
        }]
    };

    apigateway.updateStage(params, function(err, data) {
        if (err) {
            console.log(err, err.stack);
            callback(err);
        } else {
            console.log("All the request will be throttled");
            callback(null, {
                statusCode: 200,
                body: JSON.stringify({message : "All the request will be throttled!!"})
            });
        }
    });

}

module.exports = {
    handler,
};
