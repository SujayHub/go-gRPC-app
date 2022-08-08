package main

import (
	"context"
	"github.com/google/uuid"
	movieService "github.com/sujayhub/go-gRPC-app/modules/movie"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":8088"
)

var movies = make([]*movieService.MovieInfo, 0)

type MovieServer struct {
	movieService.UnimplementedMovieServer
}

func main() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	movieService.RegisterMovieServer(grpcServer, &MovieServer{})

	log.Printf("server is listening at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start the server %v", err)
	}

}

func (s *MovieServer) GetMovies(in *movieService.Empty, stream movieService.Movie_GetMoviesServer) error {
	log.Printf("received: %v", in)

	for _, movie := range movies {
		if err := stream.Send(movie); err != nil {
			return err
		}
	}

	return nil
}

func (s *MovieServer) GetMovie(ctx context.Context, in *movieService.Id) (*movieService.MovieInfo, error) {
	log.Printf("Received %v", in)

	res := &movieService.MovieInfo{}
	for _, movie := range movies {
		if movie.GetId() == in.GetValue() {
			res = movie
			break
		}
	}
	return res, nil
}

func (s *MovieServer) CreateMovie(ctx context.Context, in *movieService.MovieInfo) (*movieService.Id, error) {
	log.Printf("Received %v", in)
	res := movieService.Id{}
	res.Value = uuid.New().String()
	in.Id = res.GetValue()
	movies = append(movies, in)
	return &res, nil
}

func (s *MovieServer) UpdateMovie(ctx context.Context, in *movieService.MovieInfo) (*movieService.Status, error) {
	log.Printf("Received %v", in)
	res := movieService.Status{}
	for index, movie := range movies {
		if movie.GetId() == in.GetId() {
			movies = append(movies[:index], movies[:index+1]...)
			in.Id = movie.GetId()
			movies = append(movies, in)
			res.Value = 1
			break
		}
	}
	return &res, nil
}

func (s *MovieServer) DeleteMovie(ctx context.Context, in *movieService.Id) (*movieService.Status, error) {
	log.Printf("Received %v", in)
	res := movieService.Status{}
	for index, movie := range movies {
		if movie.GetId() == in.GetValue() {
			movies = append(movies[:index], movies[index+1:]...)
			res.Value = 1
			break
		}
	}
	return &res, nil
}
