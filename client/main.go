package main

import (
	"flag"
	"log"
	"net/http"

	pb "github.com/Mir-Labib-Hossain/grpc-crud/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

type Movie struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Genre string `json:"genre"`
}

func main() {
	flag.Parse()
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewMovieServiceClient(conn)
	r := gin.Default()

	r.GET("/movies", func(ctx *gin.Context) {
		res, err := client.GetMovies(ctx, &pb.ReadMoviesRequest{})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"movies": res.Movies,
		})
	})

	r.GET("/movie/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		res, err := client.GetMovie(ctx, &pb.ReadMovieRequest{Id: id})
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"movie": res.Movie,
		})
	})

	r.POST("/movie", func(ctx *gin.Context) {
		var movie Movie
		err := ctx.ShouldBind(&movie)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{
				"error": err.Error(),
			})
			return
		}
		data := &pb.Movie{
			Title: movie.Title,
			Genre: movie.Genre,
		}
		res, err := client.CreateMovie(ctx, &pb.CreateMovieRequest{
			Movie: data,
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"movie": res.Movie,
		})

	})

	r.PUT("/movie/:id", func(ctx *gin.Context) {
		var movie Movie
		err := ctx.ShouldBind(&movie)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		res, err := client.UpdateMovie(ctx, &pb.UpdateMovieRequest{
			Movie: &pb.Movie{
				Id:    movie.Id,
				Title: movie.Title,
				Genre: movie.Genre,
			},
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"movie": res.Movie,
		})

	})

	r.DELETE("/movie/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		res, err := client.DeleteMovie(ctx, &pb.DeleteMovieRequest{Id: id})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if res.Success == true {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Successfully deleted movie",
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error deleting movie",
			})
			return
		}
	})
	r.Run(":3000")

}
