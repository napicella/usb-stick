Add authentication to usbstick.  
It's work in progress and it's currently unused. 
To note that, Cloudformation does not allow to set conformation link, so the operation needs to be performed outside
Cloudformation.
Third party authentication can only be used in the context of a browser, while my intention was to login user from 
the command line.   
An idea would be to use: https://stackoverflow.com/questions/48537867/set-cognito-verification-type-to-link-in-cloudformation
Or https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-lambda-custom-message.html#aws-lambda-triggers-custom-message-example


https://hackernoon.com/you-should-use-ssm-parameter-store-over-lambda-env-variables-5197fc6ea45b