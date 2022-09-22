package main

import (
	"context"
	movieService "github.com/sujayhub/go-gRPC-app/modules/movie"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"time"
)

const (
	address = "localhost:8088"
)

var _ironMan = movieService.MovieInfo{
	Title: "Iron Man",
	Director: &movieService.Director{
		FirstName: "Jon",
		LastName:  "Favreau",
	},
	Imdb: 7.9,
}
var _endGame = movieService.MovieInfo{
	Title: "Avengers Endgame",
	Director: &movieService.Director{
		FirstName: "Ruso",
		LastName:  "Brothers",
	},
	Imdb: 8.4,
}
var _spiderMan = movieService.MovieInfo{
	Title: "SpiderMan No Way Home",
	Director: &movieService.Director{
		FirstName: "Jon",
		LastName:  "Wats",
	},
	Imdb: 8.3,
}
var _doctorStrange = movieService.MovieInfo{
	Id:    "Doctor Strange Multiverse of Madness",
	Title: "Doctor Strange Multiverse of Madness",
	Director: &movieService.Director{
		FirstName: "Sam",
		LastName:  "Raimi",
	},
	Imdb: 7.0,
}
var updateToIronMan3 = "Iron Man 3"
var updateIronManDirector = movieService.Director{
	FirstName: "Shane",
	LastName:  "Blake",
}
var movieIds = make([]string, 0)

func main() {
	log.Println("starting client")

	conn, err := grpc.Dial(address, grpc.WithCredentialsBundle(insecure.NewBundle()), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect %v", err)
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("uanble to deffer connection %v", err)
		}
	}(conn)

	client := movieService.NewMovieClient(conn)

	runGetMovies(client)
	runCreateMovies(client, &_ironMan)
	runCreateMovies(client, &_endGame)
	runCreateMovies(client, &_spiderMan)
	runCreateMovies(client, &_doctorStrange)
	runGetMovies(client)
	runGetMovieById(client, movieIds[1])
	updateMovie()
	runUpdateMovies(client, &_ironMan)
	runGetMovies(client)
	runDeleteMovieById(client, movieIds[1])
	runGetMovies(client)

}

func runGetMovies(client movieService.MovieClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &movieService.Empty{}
	stream, err := client.GetMovies(ctx, req)
	if err != nil {
		log.Fatalf("%v.GetMovies(_) = %v", client, err)
	}

	for {
		row, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("%v.GetMovies(_) = %v", client, err)
		}

		log.Printf("MovieInfo: %v", row)
	}
}

func runGetMovieById(client movieService.MovieClient, movieId string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &movieService.Id{Value: movieId}
	movieInfo, err := client.GetMovie(ctx, req)
	if err != nil {
		log.Fatalf("%v.GetMovie(_) = %v", client, err)
	}
	log.Printf("MovieInfo: %v", movieInfo)
}

func runCreateMovies(client movieService.MovieClient, movieInfo *movieService.MovieInfo) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &movieService.MovieInfo{
		Title:    movieInfo.GetTitle(),
		Director: movieInfo.GetDirector(),
		Imdb:     movieInfo.GetImdb(),
	}

	movieId, err := client.CreateMovie(ctx, req)
	if err != nil {
		log.Fatalf("%v.CreateMovie(_) = %v", client, err)
	}

	if movieId.GetValue() != "" {
		movieIds = append(movieIds, movieId.GetValue())
		log.Printf("Created Movie with Id: %v", movieId)
	}
}

func runUpdateMovies(client movieService.MovieClient, movieInfo *movieService.MovieInfo) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &movieService.MovieInfo{
		Title:    movieInfo.GetTitle(),
		Director: movieInfo.GetDirector(),
		Imdb:     movieInfo.GetImdb(),
	}
	updateStatus, err := client.UpdateMovie(ctx, req)
	if err != nil {
		log.Fatalf("%v.UpdateMovie(_) = %v", client, err)
	}

	if updateStatus.GetValue() == 1 {
		log.Printf("Updated Movie with Id: %v", req.GetId())
	}
}

func runDeleteMovieById(client movieService.MovieClient, movieId string) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &movieService.Id{Value: movieId}
	deleteStatus, err := client.DeleteMovie(ctx, req)
	if err != nil {
		log.Fatalf("%v.DeleteMovie(_) = %v", client, err)
	}

	if deleteStatus.GetValue() == 1 {
		log.Printf("deleted Movie with Id: %v", movieId)
	}

}

func updateMovie() {
	_ironMan.Title = updateToIronMan3
	_ironMan.Director = &updateIronManDirector
}
