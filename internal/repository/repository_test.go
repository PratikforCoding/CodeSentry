package repository_test

import (
	"testing"
	"time"

	"github.com/PratikforCoding/CodeSentry/internal/models"
	"github.com/PratikforCoding/CodeSentry/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestSaveAnalysis(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := repository.NewAnalysisRepository(mt.Client.Database("codesentry"))
		req := models.AnalyzeRequest{Code: "package main"}
		res := models.AnalysisResponse{Language: "go"}

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		err := repo.SaveAnalysis(req, res)
		assert.NoError(t, err)
	})
}

func TestGetAllAnalyses(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success with language filter", func(mt *mtest.T) {
		repo := repository.NewAnalysisRepository(mt.Client.Database("codesentry"))
		testTime := time.Now()
		expected := []models.Analysis{
			{
				ID:        primitive.NewObjectID(),
				Code:      "package main",
				Language:  "go",
				Timestamp: testTime,
			},
		}

		first := mtest.CreateCursorResponse(1, "test.analyses", mtest.FirstBatch, bson.D{
			{"_id", expected[0].ID},
			{"code", expected[0].Code},
			{"language", expected[0].Language},
			{"timestamp", expected[0].Timestamp},
		})
		killCursors := mtest.CreateCursorResponse(0, "test.analyses", mtest.NextBatch)
		mt.AddMockResponses(first, killCursors)

		results, err := repo.GetAllAnalyses("go")
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, expected[0].ID, results[0].ID)
	})
}

func TestGetAnalysisByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := repository.NewAnalysisRepository(mt.Client.Database("codesentry"))
		oid := primitive.NewObjectID()
		expected := models.Analysis{ID: oid}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.analyses", mtest.FirstBatch, bson.D{
			{"_id", oid},
		}))

		result, err := repo.GetAnalysisByID(oid.Hex())
		require.NoError(t, err)
		assert.Equal(t, expected.ID, result.ID)
	})

	mt.Run("invalid id", func(mt *mtest.T) {
		repo := repository.NewAnalysisRepository(mt.Client.Database("codesentry"))
		_, err := repo.GetAnalysisByID("invalid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid ID format")
	})

	mt.Run("not found", func(mt *mtest.T) {
		repo := repository.NewAnalysisRepository(mt.Client.Database("codesentry"))
		oid := primitive.NewObjectID()

		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    0,
			Message: "not found",
		}))

		_, err := repo.GetAnalysisByID(oid.Hex())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestUpdateAnalysis(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := repository.NewAnalysisRepository(mt.Client.Database("codesentry"))
		oid := primitive.NewObjectID()
		updateReq := models.UpdateAnalysisRequest{
			Title: "New Title",
			Tags:  []string{"security"},
		}

		mt.AddMockResponses(bson.D{
			{"ok", 1},
			{"n", 1},
			{"matchedCount", 1},
			{"modifiedCount", 1},
		})

		err := repo.UpdateAnalysis(oid.Hex(), updateReq)
		assert.NoError(t, err)
	})

	mt.Run("not found", func(mt *mtest.T) {
		repo := repository.NewAnalysisRepository(mt.Client.Database("codesentry"))
		oid := primitive.NewObjectID()

		mt.AddMockResponses(bson.D{
			{"ok", 1},
			{"n", 0},
			{"matchedCount", 0},
			{"modifiedCount", 0},
		})

		err := repo.UpdateAnalysis(oid.Hex(), models.UpdateAnalysisRequest{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestDeleteAnalysis(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := repository.NewAnalysisRepository(mt.Client.Database("codesentry"))
		oid := primitive.NewObjectID()

		mt.AddMockResponses(bson.D{
			{"ok", 1},
			{"n", 1},
		})

		err := repo.DeleteAnalysis(oid.Hex())
		assert.NoError(t, err)
	})

	mt.Run("not found", func(mt *mtest.T) {
		repo := repository.NewAnalysisRepository(mt.Client.Database("codesentry"))
		oid := primitive.NewObjectID()

		mt.AddMockResponses(bson.D{
			{"ok", 1},
			{"n", 0},
		})

		err := repo.DeleteAnalysis(oid.Hex())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
