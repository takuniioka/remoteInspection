import * as cdk from 'aws-cdk-lib'
import * as cloudfront from 'aws-cdk-lib/aws-cloudfront'
import * as origins from 'aws-cdk-lib/aws-cloudfront-origins'
import * as cognito from 'aws-cdk-lib/aws-cognito'
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb'
import * as s3 from 'aws-cdk-lib/aws-s3'
import { Construct } from 'constructs'

/**
 * CDK stack defining infrastructure for the inspection tool. It
 * provisions DynamoDB tables, S3 buckets for evidence and UI hosting,
 * a CloudFront distribution for the UI and a Cognito user pool for
 * authentication. The resources use `RemovalPolicy.DESTROY` to simplify
 * development teardown; change to RETAIN for production.
 */
export class InspectionToolStack extends cdk.Stack {
    constructor(scope: Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props)

        // DynamoDB Tables
        const inspectionsTable = new dynamodb.Table(this, 'InspectionsTable', {
            tableName: 'Inspections',
            partitionKey: { name: 'inspectionId', type: dynamodb.AttributeType.STRING },
            billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
            removalPolicy: cdk.RemovalPolicy.DESTROY,
        })

        const checklistItemsTable = new dynamodb.Table(this, 'ChecklistItemsTable', {
            tableName: 'ChecklistItems',
            partitionKey: { name: 'inspectionId', type: dynamodb.AttributeType.STRING },
            sortKey: { name: 'itemId', type: dynamodb.AttributeType.STRING },
            billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
            removalPolicy: cdk.RemovalPolicy.DESTROY,
        })

        const issuesTable = new dynamodb.Table(this, 'IssuesTable', {
            tableName: 'Issues',
            partitionKey: { name: 'inspectionId', type: dynamodb.AttributeType.STRING },
            sortKey: { name: 'issueId', type: dynamodb.AttributeType.STRING },
            billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
            removalPolicy: cdk.RemovalPolicy.DESTROY,
        })

        const captureRequestsTable = new dynamodb.Table(this, 'CaptureRequestsTable', {
            tableName: 'CaptureRequests',
            partitionKey: { name: 'inspectionId', type: dynamodb.AttributeType.STRING },
            sortKey: { name: 'captureRequestId', type: dynamodb.AttributeType.STRING },
            billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
            removalPolicy: cdk.RemovalPolicy.DESTROY,
        })

        const evidencePhotosTable = new dynamodb.Table(this, 'EvidencePhotosTable', {
            tableName: 'EvidencePhotos',
            partitionKey: { name: 'inspectionId', type: dynamodb.AttributeType.STRING },
            sortKey: { name: 'photoId', type: dynamodb.AttributeType.STRING },
            billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
            removalPolicy: cdk.RemovalPolicy.DESTROY,
        })

        const annotationTemplatesTable = new dynamodb.Table(this, 'AnnotationTemplatesTable', {
            tableName: 'AnnotationTemplates',
            partitionKey: { name: 'templateId', type: dynamodb.AttributeType.STRING },
            sortKey: { name: 'version', type: dynamodb.AttributeType.NUMBER },
            billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
            removalPolicy: cdk.RemovalPolicy.DESTROY,
        })

        const viewerLinksTable = new dynamodb.Table(this, 'ViewerLinksTable', {
            tableName: 'ViewerLinks',
            partitionKey: { name: 'viewerAccessToken', type: dynamodb.AttributeType.STRING },
            billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
            removalPolicy: cdk.RemovalPolicy.DESTROY,
            timeToLiveAttribute: 'expiresAt',
        })

        // S3 Buckets
        const evidenceBucket = new s3.Bucket(this, 'EvidenceBucket', {
            bucketName: `inspection-evidence-${cdk.Stack.of(this).account}-${cdk.Stack.of(this).region}`,
            removalPolicy: cdk.RemovalPolicy.DESTROY,
            cors: [
                {
                    allowedMethods: [s3.HttpMethods.GET, s3.HttpMethods.PUT, s3.HttpMethods.POST],
                    allowedOrigins: ['*'],
                    allowedHeaders: ['*'],
                },
            ],
        })

        const uiBucket = new s3.Bucket(this, 'UIBucket', {
            bucketName: `inspection-ui-${cdk.Stack.of(this).account}-${cdk.Stack.of(this).region}`,
            removalPolicy: cdk.RemovalPolicy.DESTROY,
            websiteIndexDocument: 'index.html',
            blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL,
        })

        // CloudFront Distribution for UI
        const distribution = new cloudfront.Distribution(this, 'UIDistribution', {
            defaultBehavior: {
                origin: new origins.S3Origin(uiBucket),
                viewerProtocolPolicy: cloudfront.ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
            },
            defaultRootObject: 'index.html',
            errorResponses: [
                {
                    httpStatus: 404,
                    responseHttpStatus: 200,
                    responsePagePath: '/index.html',
                    ttl: cdk.Duration.minutes(5),
                },
            ],
        })

        // Cognito User Pool
        const userPool = new cognito.UserPool(this, 'InspectionUserPool', {
            userPoolName: 'inspection-tool-pool',
            selfSignUpEnabled: false,
            signInAliases: {
                email: true,
            },
            autoVerifiedAttributes: [cognito.AutoVerifiedAttr.EMAIL],
            passwordPolicy: {
                minLength: 12,
                requireLowercase: true,
                requireUppercase: true,
                requireDigits: true,
                requireSymbols: true,
            },
            accountRecovery: cognito.AccountRecovery.EMAIL_ONLY,
            removalPolicy: cdk.RemovalPolicy.DESTROY,
        })

        const userPoolClient = userPool.addClient('InspectionClient', {
            clientName: 'inspection-tool-client',
            auth: cognito.AuthFlow.ALLOW_USER_PASSWORD_AUTH_FLOW,
            generateSecret: false,
        })

        // Outputs
        new cdk.CfnOutput(this, 'InspectionsTableName', {
            value: inspectionsTable.tableName,
        })
        new cdk.CfnOutput(this, 'EvidenceBucketName', {
            value: evidenceBucket.bucketName,
        })
        new cdk.CfnOutput(this, 'UIBucketName', {
            value: uiBucket.bucketName,
        })
        new cdk.CfnOutput(this, 'CloudFrontDistributionDomain', {
            value: distribution.distributionDomainName,
        })
        new cdk.CfnOutput(this, 'CognitoUserPoolId', {
            value: userPool.userPoolId,
        })
        new cdk.CfnOutput(this, 'CognitoClientId', {
            value: userPoolClient.userPoolClientId,
        })
    }
}
