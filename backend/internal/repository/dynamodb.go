package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/inspection-tool/backend/internal/domain"
)

// DynamoDBClient wraps the DynamoDB client
type DynamoDBClient struct {
	client *dynamodb.Client
}

// NewDynamoDBClient creates a new DynamoDB client wrapper
func NewDynamoDBClient(client *dynamodb.Client) *DynamoDBClient {
	return &DynamoDBClient{client: client}
}

// InspectionRepositoryImpl implements domain.InspectionRepository
type InspectionRepositoryImpl struct {
	db *DynamoDBClient
}

// NewInspectionRepository creates a new inspection repository
func NewInspectionRepository(db *DynamoDBClient) domain.InspectionRepository {
	return &InspectionRepositoryImpl{db: db}
}

// Create creates a new inspection
func (r *InspectionRepositoryImpl) Create(ctx context.Context, inspection *domain.Inspection) error {
	if inspection.InspectionID == "" {
		inspection.InspectionID = uuid.New().String()
	}
	inspection.CreatedAt = time.Now()

	item, err := attributevalue.MarshalMap(inspection)
	if err != nil {
		return fmt.Errorf("failed to marshal inspection: %w", err)
	}

	_, err = r.db.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("Inspections"),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to put inspection: %w", err)
	}

	return nil
}

// GetByID retrieves an inspection by ID
func (r *InspectionRepositoryImpl) GetByID(ctx context.Context, inspectionID string) (*domain.Inspection, error) {
	result, err := r.db.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String("Inspections"),
		Key: map[string]types.AttributeValue{
			"inspectionId": &types.AttributeValueMemberS{Value: inspectionID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get inspection: %w", err)
	}

	if result.Item == nil {
		return nil, domain.NewAppError(domain.ErrCodeNotFound, "inspection not found", nil)
	}

	var inspection domain.Inspection
	err = attributevalue.UnmarshalMap(result.Item, &inspection)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal inspection: %w", err)
	}

	return &inspection, nil
}

// ListAll retrieves all inspections
func (r *InspectionRepositoryImpl) ListAll(ctx context.Context) ([]*domain.Inspection, error) {
	result, err := r.db.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String("Inspections"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scan inspections: %w", err)
	}

	var inspections []*domain.Inspection
	err = attributevalue.UnmarshalListOfMaps(result.Items, &inspections)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal inspections: %w", err)
	}

	return inspections, nil
}

// Update updates an inspection
func (r *InspectionRepositoryImpl) Update(ctx context.Context, inspection *domain.Inspection) error {
	item, err := attributevalue.MarshalMap(inspection)
	if err != nil {
		return fmt.Errorf("failed to marshal inspection: %w", err)
	}

	_, err = r.db.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("Inspections"),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to update inspection: %w", err)
	}

	return nil
}

// ChecklistItemRepositoryImpl implements domain.ChecklistItemRepository
type ChecklistItemRepositoryImpl struct {
	db *DynamoDBClient
}

// NewChecklistItemRepository creates a new checklist item repository
func NewChecklistItemRepository(db *DynamoDBClient) domain.ChecklistItemRepository {
	return &ChecklistItemRepositoryImpl{db: db}
}

// Create creates a new checklist item
func (r *ChecklistItemRepositoryImpl) Create(ctx context.Context, item *domain.ChecklistItem) error {
	if item.ItemID == "" {
		item.ItemID = uuid.New().String()
	}
	item.UpdatedAt = time.Now()

	attrMap, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal checklist item: %w", err)
	}

	_, err = r.db.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("ChecklistItems"),
		Item:      attrMap,
	})
	if err != nil {
		return fmt.Errorf("failed to put checklist item: %w", err)
	}

	return nil
}

// GetByID retrieves a checklist item by ID
func (r *ChecklistItemRepositoryImpl) GetByID(ctx context.Context, inspectionID, itemID string) (*domain.ChecklistItem, error) {
	result, err := r.db.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String("ChecklistItems"),
		Key: map[string]types.AttributeValue{
			"inspectionId": &types.AttributeValueMemberS{Value: inspectionID},
			"itemId":       &types.AttributeValueMemberS{Value: itemID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get checklist item: %w", err)
	}

	if result.Item == nil {
		return nil, domain.NewAppError(domain.ErrCodeNotFound, "checklist item not found", nil)
	}

	var item domain.ChecklistItem
	err = attributevalue.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal checklist item: %w", err)
	}

	return &item, nil
}

// ListByInspection retrieves all checklist items for an inspection
func (r *ChecklistItemRepositoryImpl) ListByInspection(ctx context.Context, inspectionID string) ([]*domain.ChecklistItem, error) {
	result, err := r.db.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String("ChecklistItems"),
		KeyConditionExpression: aws.String("inspectionId = :inspectionId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":inspectionId": &types.AttributeValueMemberS{Value: inspectionID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query checklist items: %w", err)
	}

	var items []*domain.ChecklistItem
	err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal checklist items: %w", err)
	}

	return items, nil
}

// Update updates a checklist item
func (r *ChecklistItemRepositoryImpl) Update(ctx context.Context, item *domain.ChecklistItem) error {
	item.UpdatedAt = time.Now()

	attrMap, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal checklist item: %w", err)
	}

	_, err = r.db.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("ChecklistItems"),
		Item:      attrMap,
	})
	if err != nil {
		return fmt.Errorf("failed to update checklist item: %w", err)
	}

	return nil
}
