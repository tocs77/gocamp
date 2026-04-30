package mongodb

import (
	"context"
	"errors"

	"sch-grpc/internals/models"
	"sch-grpc/pkg/utils"
	pb "sch-grpc/proto/gen"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddStudents(ctx context.Context, students []*pb.Student) ([]*pb.Student, error) {
	newStudents := make([]models.Student, 0, len(students))
	for _, student := range students {
		var modelStudent models.Student
		utils.MapStructFields(student, &modelStudent)
		newStudents = append(newStudents, modelStudent)
	}

	pbStudents := make([]*pb.Student, 0, len(newStudents))
	for _, student := range newStudents {
		result, err := MongoClient.Database("sch-db").Collection("students").InsertOne(ctx, student)
		if err != nil {
			return nil, utils.HandleError(err, "failed to add student to MongoDB")
		}

		objID, ok := result.InsertedID.(primitive.ObjectID)
		if !ok {
			return nil, utils.HandleError(errors.New("failed to get object ID"), "failed to get object ID")
		}
		student.ID = objID.Hex()

		pbStudent := &pb.Student{}
		utils.MapStructFields(student, pbStudent)
		pbStudents = append(pbStudents, pbStudent)
	}

	return pbStudents, nil
}

func GetStudents(ctx context.Context, filter bson.M, sort bson.D, pageNumber int32, pageSize int32) ([]*pb.Student, error) {
	collection := MongoClient.Database("sch-db").Collection("students")
	cursor, err := collection.Find(ctx, filter, options.Find().SetSort(sort).SetSkip(int64((pageNumber-1)*pageSize)).SetLimit(int64(pageSize)))
	if err != nil {
		return nil, utils.HandleError(err, "failed to get students from MongoDB")
	}
	defer cursor.Close(ctx)

	students := make([]*pb.Student, 0)
	for cursor.Next(ctx) {
		var student models.Student
		if err := cursor.Decode(&student); err != nil {
			return nil, utils.HandleError(err, "failed to decode student from MongoDB")
		}

		pbStudent := &pb.Student{}
		utils.MapStructFields(student, pbStudent)
		students = append(students, pbStudent)
	}

	return students, nil
}

func UpdateStudents(ctx context.Context, students []*pb.Student) ([]*pb.Student, error) {
	updatedStudents, err := updateInDb[*pb.Student, *pb.Student](ctx, "students", students, bson.M{})
	if err != nil {
		return nil, utils.HandleError(err, "failed to update students in MongoDB")
	}

	return updatedStudents, nil
}

func DeleteStudents(ctx context.Context, studentIDs []string) ([]string, error) {
	deletedIDs, err := deleteInDbByID(ctx, "students", studentIDs)
	if err != nil {
		return nil, utils.HandleError(err, "failed to delete students in MongoDB")
	}

	return deletedIDs, nil
}
