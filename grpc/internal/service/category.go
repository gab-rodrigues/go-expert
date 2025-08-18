package service

import (
	"context"
	"google.golang.org/grpc"
	"grpc-study/internal/database"
	"grpc-study/internal/pb"
	"io"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDB database.Category
}

func NewCategoryService(categoryDB database.Category) *CategoryService {
	return &CategoryService{
		CategoryDB: categoryDB,
	}
}

func (s *CategoryService) CreateCategory(ctx context.Context, input *pb.CreateCategoryRequest) (*pb.CategoryResponse, error) {
	category, err := s.CategoryDB.Create(input.Name, input.Description)
	if err != nil {
		return nil, err
	}
	return &pb.CategoryResponse{
		Category: &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		},
	}, nil
}

func (s *CategoryService) ListCategories(ctx context.Context, in *pb.Blank) (*pb.CategoryListResponse, error) {
	categories, err := s.CategoryDB.FindAll()
	if err != nil {
		return nil, err
	}

	var pbCategories []*pb.Category
	for _, category := range categories {
		pbCategories = append(pbCategories, &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})
	}

	return &pb.CategoryListResponse{
		Categories: pbCategories,
	}, nil
}

func (s *CategoryService) GetCategoryByID(ctx context.Context, input *pb.GetCategoryRequest) (*pb.Category, error) {
	category, err := s.CategoryDB.FindByID(input.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}, nil
}

func (s *CategoryService) CreateCategoryStream(stream grpc.ClientStreamingServer[pb.CreateCategoryRequest, pb.CategoryListResponse]) error {
	categories := &pb.CategoryListResponse{}
	for {
		category, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(categories)
		}
		if err != nil {
			return err
		}
		newCategory, err := s.CategoryDB.Create(category.Name, category.Description)
		if err != nil {
			return err
		}
		categories.Categories = append(categories.Categories, &pb.Category{
			Id:          newCategory.ID,
			Name:        newCategory.Name,
			Description: newCategory.Description,
		})
	}
}

func (s *CategoryService) CreateCategoryStreamBidirectional(stream grpc.BidiStreamingServer[pb.CreateCategoryRequest, pb.Category]) error {
	for {
		category, err := stream.Recv()
		if err == io.EOF {
			return nil // End of stream
		}
		if err != nil {
			return err
		}

		newCategory, err := s.CategoryDB.Create(category.Name, category.Description)
		if err != nil {
			return err
		}

		if err := stream.Send(&pb.Category{
			Id:          newCategory.ID,
			Name:        newCategory.Name,
			Description: newCategory.Description,
		}); err != nil {
			return err
		}
	}
}
