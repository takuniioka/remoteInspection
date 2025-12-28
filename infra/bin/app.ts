import * as cdk from 'aws-cdk-lib'
import { InspectionToolStack } from './lib/inspection-stack'

const app = new cdk.App()

new InspectionToolStack(app, 'InspectionToolStack', {
    env: {
        account: process.env.CDK_DEFAULT_ACCOUNT,
        region: process.env.CDK_DEFAULT_REGION || 'ap-northeast-1',
    },
})

app.synth()
