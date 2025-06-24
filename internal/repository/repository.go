package repository

import (
	"context"
	"errors"
	"github.com/PratikforCoding/CodeSentry/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type AnalysisRepositoryInterface interface {
	SaveAnalysis(req models.AnalyzeRequest, res models.AnalysisResponse) error
	GetAllAnalyses(language string) ([]models.Analysis, error)
	GetAnalysisByID(id string) (models.Analysis, error)
	UpdateAnalysis(id string, updateReq models.UpdateAnalysisRequest) error
	DeleteAnalysis(id string) error
}
type AnalysisRepository struct {
	db *mongo.Database
}

func NewAnalysisRepository(db *mongo.Database) *AnalysisRepository {
	return &AnalysisRepository{
		db: db,
	}
}

func (ar *AnalysisRepository) SaveAnalysis(req models.AnalyzeRequest, res models.AnalysisResponse) error {
	collection := ar.db.Collection("analyses")

	doc := bson.M{
		"code":      req.Code,
		"language":  res.Language,
		"options":   req.Options,
		"response":  res,
		"timestamp": time.Now(),
	}

	_, err := collection.InsertOne(context.Background(), doc)
	return err
}

func (ar *AnalysisRepository) GetAllAnalyses(language string) ([]models.Analysis, error) {
	collection := ar.db.Collection("analyses")

	// Build filter
	filter := bson.M{}
	if language != "" {
		filter["language"] = language
	}

	// Sort by timestamp (newest first)
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "timestamp", Value: -1}})

	cursor, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var analyses []models.Analysis
	if err = cursor.All(context.Background(), &analyses); err != nil {
		return nil, err
	}

	return analyses, nil
}

func (ar *AnalysisRepository) GetAnalysisByID(id string) (models.Analysis, error) {
	collection := ar.db.Collection("analyses")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Analysis{}, errors.New("invalid ID format")
	}

	var analysis models.Analysis
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&analysis)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Analysis{}, errors.New("not found")
		}
		return models.Analysis{}, err
	}

	return analysis, nil
}

func (ar *AnalysisRepository) UpdateAnalysis(id string, updateReq models.UpdateAnalysisRequest) error {
	collection := ar.db.Collection("analyses")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	// Build update document
	updateDoc := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	// Only update fields that are provided
	if updateReq.Title != "" {
		updateDoc["$set"].(bson.M)["title"] = updateReq.Title
	}
	if updateReq.Tags != nil {
		updateDoc["$set"].(bson.M)["tags"] = updateReq.Tags
	}

	result, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		updateDoc,
	)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("not found")
	}

	return nil
}

func (ar *AnalysisRepository) DeleteAnalysis(id string) error {
	collection := ar.db.Collection("analyses")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("not found")
	}

	return nil
}
