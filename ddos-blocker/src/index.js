const AWS = require('aws-sdk'); // eslint-disable-line import/no-extraneous-dependencies

// https://docs.aws.amazon.com/AWSJavaScriptSDK/latest/AWS/APIGateway.html#createApiKey-property
const apigateway = new AWS.APIGateway({apiVersion: '2015-07-09', region: process.env.AWS_REGION});


const handler = (event, context, callback) => {
    const createKeyParams = {
        customerId: 'tooManyCalls',
        description: 'tooManyCalls',
        enabled: true || false,
        generateDistinctId: true || false,
        name: 'tooManyCalls',
        value: 'STRING_VALUE'
    };
    apigateway.createApiKey(createKeyParams, function (err, data) {
        if (err) {
            console.log(err, err.stack);
            callback(err)
        } else {
            callback(null, {
                statusCode: 200,
                body: JSON.stringify({message : "Api key created!"})
            });
        }
    });
};

module.exports = {
    handler,
};
