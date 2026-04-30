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

func AddTeachers(ctx context.Context, teachers []*pb.Teacher) ([]*pb.Teacher, error) {
	newTeachers := make([]models.Teacher, 0, len(teachers))
	for _, teacher := range teachers {
		var modelTeacher models.Teacher
		utils.MapStructFields(teacher, &modelTeacher)
		newTeachers = append(newTeachers, modelTeacher)
	}

	pbTeachers := make([]*pb.Teacher, 0, len(newTeachers))
	for _, teacher := range newTeachers {
		result, err := MongoClient.Database("sch-db").Collection("teachers").InsertOne(ctx, teacher)
		if err != nil {
			return nil, utils.HandleError(err, "failed to add teacher to MongoDB")
		}
		objId, ok := result.InsertedID.(primitive.ObjectID)
		if !ok {
			return nil, utils.HandleError(errors.New("failed to get object ID"), "failed to get object ID")
		}
		teacher.ID = objId.Hex()
		pbTeacher := &pb.Teacher{}
		utils.MapStructFields(teacher, pbTeacher)
		pbTeachers = append(pbTeachers, pbTeacher)
	}
	return pbTeachers, nil
}

func GetTeachers(ctx context.Context, filter bson.M, sort bson.D) ([]*pb.Teacher, error) {
	collection := MongoClient.Database("sch-db").Collection("teachers")
	cursor, err := collection.Find(ctx, filter, options.Find().SetSort(sort))
	if err != nil {
		return nil, utils.HandleError(err, "failed to get teachers from MongoDB")
	}
	defer cursor.Close(ctx)
	teachers := make([]*pb.Teacher, 0)
	for cursor.Next(ctx) {
		var teacher models.Teacher
		err := cursor.Decode(&teacher)
		if err != nil {
			return nil, utils.HandleError(err, "failed to decode teacher from MongoDB")
		}
		pbTeacher := &pb.Teacher{}
		utils.MapStructFields(teacher, pbTeacher)
		teachers = append(teachers, pbTeacher)
	}
	return teachers, nil
}

func UpdateTeachers(ctx context.Context, teachers []*pb.UpdateTeacher) ([]*pb.Teacher, error) {
	updatedTeachers, err := updateInDb[*pb.UpdateTeacher, *pb.Teacher](ctx, "teachers", teachers, bson.M{})
	if err != nil {
		return nil, utils.HandleError(err, "failed to update teachers in MongoDB")
	}
	return updatedTeachers, nil
}

func DeleteTeachers(ctx context.Context, teachersIds []string) ([]string, error) {
	deletedIds, err := deleteInDbByID(ctx, "teachers", teachersIds)
	if err != nil {
		return nil, utils.HandleError(err, "failed to delete teachers in MongoDB")
	}
	return deletedIds, nil
}

func GetTeacherStudents(ctx context.Context, teacherID string) ([]*pb.Student, error) {
	objectID, err := primitive.ObjectIDFromHex(teacherID)
	if err != nil {
		return nil, utils.HandleError(err, "failed to convert teacher ID to ObjectID")
	}

	teacherCollection := MongoClient.Database("sch-db").Collection("teachers")
	teacher := models.Teacher{}
	err = teacherCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&teacher)
	if err != nil {
		return nil, utils.HandleError(err, "failed to get teacher from MongoDB")
	}

	studentsCollection := MongoClient.Database("sch-db").Collection("students")
	cursor, err := studentsCollection.Find(ctx, bson.M{"class": teacher.Class})
	if err != nil {
		return nil, utils.HandleError(err, "failed to get students from MongoDB")
	}
	defer cursor.Close(ctx)
	students := make([]*pb.Student, 0)
	for cursor.Next(ctx) {
		var student models.Student
		err := cursor.Decode(&student)
		if err != nil {
			return nil, utils.HandleError(err, "failed to decode student from MongoDB")
		}
		pbStudent := &pb.Student{}
		utils.MapStructFields(student, pbStudent)
		students = append(students, pbStudent)
	}
	return students, nil
}

func GetTeacherStudentsCount(ctx context.Context, teacherID string) (int, error) {
	objectID, err := primitive.ObjectIDFromHex(teacherID)
	if err != nil {
		return 0, utils.HandleError(err, "failed to convert teacher ID to ObjectID")
	}

	teacherCollection := MongoClient.Database("sch-db").Collection("teachers")
	teacher := models.Teacher{}
	err = teacherCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&teacher)
	if err != nil {
		return 0, utils.HandleError(err, "failed to get teacher from MongoDB")
	}

	studentsCollection := MongoClient.Database("sch-db").Collection("students")
	count, err := studentsCollection.CountDocuments(ctx, bson.M{"class": teacher.Class})
	if err != nil {
		return 0, utils.HandleError(err, "failed to get students count from MongoDB")
	}
	return int(count), nil
}
