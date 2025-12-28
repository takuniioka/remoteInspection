module github.com/inspection-tool/backend

go 1.22

require (
	github.com/aws/aws-sdk-go-v2 v1.24.0
	github.com/aws/aws-sdk-go-v2/config v1.26.1
	github.com/aws/aws-sdk-go-v2/credentials v1.16.12
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.13.0
	github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider v1.32.0
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.30.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.47.0
	github.com/go-chi/chi/v5 v5.0.11
	github.com/go-chi/cors v1.2.1
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/kelseyhightower/envconfig v1.4.0
	golang.org/x/crypto v0.21.0
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.0
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.15.0
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.2.9
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.5.9
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.0
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.0
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.3.0
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.0
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.17.0
	github.com/aws/aws-sdk-go-v2/service/sso v1.18.5
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.21.5
	github.com/aws/aws-sdk-go-v2/service/sts v1.26.5
	github.com/aws/smithy-go v1.19.0
	github.com/jmespath/go-jmespath v0.4.0
)
