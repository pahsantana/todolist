package repositories

import (
	"context"
	"time"

	"github.com/pahsantana/todolist/internal/domain/entities"
	"github.com/pahsantana/todolist/internal/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/pahsantana/todolist/internal/dto"
)

const (
	idField         = "_id"
	setOperator     = "$set"
	filterStatus    = "status"
	filterPriority  = "priority"
	tasksCollection = "tasks"
)

type TaskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository(db *mongo.Database) repository.TaskRepository {
	return &TaskRepository{
		collection: db.Collection(tasksCollection),
	}
}

func (r *TaskRepository) ListCountByStatus(ctx context.Context) (*dto.TaskSummary, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: idField, Value: "$status"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	summary := &dto.TaskSummary{}

	for cursor.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Count int64  `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}

		switch entities.Status(result.ID) {
		case entities.Pending:
			summary.Pending = result.Count
		case entities.InProgress:
			summary.InProgress = result.Count
		case entities.Completed:
			summary.Completed = result.Count
		case entities.Cancelled:
			summary.Cancelled = result.Count
		}
	}

	return summary, nil
}

func (r *TaskRepository) Create(ctx context.Context, task *entities.Task) error {
	_, err := r.collection.InsertOne(ctx, task)
	return err
}

func (r *TaskRepository) FindAll(ctx context.Context, filters map[string]string) ([]entities.Task, error) {
	filter := bson.M{}
	if status, ok := filters[filterStatus]; ok && status != "" {
		filter[filterStatus] = entities.Status(status)
	}
	if priority, ok := filters[filterPriority]; ok && priority != "" {
		filter[filterPriority] = entities.Priority(priority)
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []entities.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepository) FindByID(ctx context.Context, id string) (*entities.Task, error) {
	var task entities.Task
	err := r.collection.FindOne(ctx, bson.M{idField: id}).Decode(&task)
	if err == mongo.ErrNoDocuments {
		return nil, entities.TaskNotFound
	}
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) Update(ctx context.Context, id string, task *entities.Task) error {
	task.UpdatedAt = time.Now()
	result, err := r.collection.UpdateOne(ctx, bson.M{idField: id}, bson.M{setOperator: task})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return entities.TaskNotFound
	}
	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{idField: id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return entities.TaskNotFound
	}
	return nil
}