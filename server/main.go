package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/Mir-Labib-Hossain/grpc-crud/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func main() {
	fmt.Println("gRPC server running")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMovieServiceServer(s, &server{})
	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func init() {
	DatabaseConnection()
}

type Movie struct {
	ID        string `gorm:"primarykey"`
	Title     string
	Genre     string
	CreatedAt time.Time `gorm:"autoCreateTime:false"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:false"`
}

func DatabaseConnection() {
	host := "localhost"
	port := "5432"
	dbName := "postgres"
	dbUser := "postgres"
	dbPassword := "1234"
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, dbUser, dbName, dbPassword)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to db ", err)
		return
	}

	// Call AutoMigrate after successful connection
	err = DB.AutoMigrate(&Movie{})
	if err != nil {
		log.Fatal("Error during migration ", err)
		return
	}

	fmt.Println("DB connected successfully")
	fmt.Println("DSN:", dsn)
}

var (
	port = flag.Int("port", 50051, "gRPC server port")
)

type server struct {
	pb.UnimplementedMovieServiceServer
}

func (*server) CreateMovie(ctx context.Context, req *pb.CreateMovieRequest) (*pb.CreateMovieResponse, error) {
	fmt.Println("Create Movie")
	movie := req.GetMovie()
	movie.Id = uuid.New().String()
	data := Movie{
		ID:    movie.GetId(),
		Title: movie.GetTitle(),
		Genre: movie.GetGenre(),
	}
	res := DB.Create(&data)
	if res.RowsAffected == 0 {
		return nil, errors.New("failed to create movie")
	}
	return &pb.CreateMovieResponse{
		Movie: &pb.Movie{
			Id:    movie.GetId(),
			Title: movie.GetTitle(),
			Genre: movie.GetGenre(),
		},
	}, nil
}

func (*server) GetMovie(ctx context.Context, req *pb.ReadMovieRequest) (*pb.ReadMovieResponse, error) {
	fmt.Println("Get Movie", req.GetId())
	var movie Movie
	res := DB.Find(&movie, "id = ?", req.GetId())
	if res.RowsAffected == 0 {
		return nil, errors.New("movie not found")
	}
	return &pb.ReadMovieResponse{
		Movie: &pb.Movie{
			Id:    movie.ID,
			Title: movie.Title,
			Genre: movie.Genre,
		},
	}, nil

}

func (*server) GetMovies(ctx context.Context, req *pb.ReadMoviesRequest) (*pb.ReadMoviesResponse, error) {
	fmt.Println("Get Movies")
	movies := []*pb.Movie{}
	res := DB.Find(&movies)
	if res.RowsAffected == 0 {
		return nil, errors.New("movie not found")
	}
	return &pb.ReadMoviesResponse{
		Movies: movies,
	}, nil
}

func (*server) UpdateMovie(ctx context.Context, req *pb.UpdateMovieRequest) (*pb.UpdateMovieResponse, error) {
	fmt.Println("Update Movie ", req.Movie.GetId())

	var movie Movie
	reqMovie := req.GetMovie()

	res := DB.Model(&movie).Where("id=?", reqMovie.Id).Updates(Movie{
		Title: reqMovie.Title,
		Genre: reqMovie.Genre,
	})
	if res.RowsAffected == 0 {
		return nil, errors.New("movie not found to update")
	}
	return &pb.UpdateMovieResponse{
		Movie: &pb.Movie{
			Id:    movie.ID,
			Title: movie.Title,
			Genre: movie.Genre,
		},
	}, nil
}

func (*server) DeleteMovie(ctx context.Context, req *pb.DeleteMovieRequest) (*pb.DeleteMovieResponse, error) {
	fmt.Println("Update Movie ", req.GetId())
	var movie Movie
	res := DB.Where("id=?", req.GetId()).Delete(&movie)
	if res.RowsAffected == 0 {
		return nil, errors.New("movie not found to delete")
	}
	return &pb.DeleteMovieResponse{
		Success: true,
	}, nil
}
