package httpserver

import (
	"context"
	"net/http"

	"github.com/RozmiDan/gameReviewHub/internal/config"
	"github.com/RozmiDan/gameReviewHub/internal/entity"

	addcomment "github.com/RozmiDan/gameReviewHub/internal/controller/http/handlers/addcomment"
	gametopic "github.com/RozmiDan/gameReviewHub/internal/controller/http/handlers/gametopic"
	listcomments "github.com/RozmiDan/gameReviewHub/internal/controller/http/handlers/listcomments"
	mainpage "github.com/RozmiDan/gameReviewHub/internal/controller/http/handlers/mainpage"
	postrating "github.com/RozmiDan/gameReviewHub/internal/controller/http/handlers/postrating"
	middleware_main "github.com/RozmiDan/gameReviewHub/internal/controller/http/middleware/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

type GameUseCase interface {
	GetListGames(ctx context.Context, limit, offset int32) ([]entity.GameInList, error)
	GetTopicGame(ctx context.Context, gameID string) (*entity.Game, error)

	PostRating(ctx context.Context, gameID, userID string, rating int32) error

	GetListComments(ctx context.Context, gameID string, limit, offset int32) ([]entity.Comment, error)
	AddComment(ctx context.Context, gameID, userID, text string) (string, error)
}

func InitServer(cnfg *config.Config, logger *zap.Logger, uc GameUseCase) *http.Server {
	logger = logger.With(zap.String("layer", "mainController"))

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware_main.MyLogger(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	// router.Use(cors.Handler(cors.Options{
	// 	AllowedOrigins:   []string{"*"},
	// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowedHeaders:   []string{"*"},
	// 	AllowCredentials: true,
	// 	MaxAge:           300,
	// }))

	router.Route("/games", func(r chi.Router) {
		// GET  /games?limit=&offset=
		r.Get("/", mainpage.NewMainpageHandler(logger, uc))

		// для game_id
		r.Route("/{game_id}", func(r chi.Router) {
			// GET   /games/{game_id}
			r.Get("/", gametopic.NewGameTopicHandler(logger, uc))

			// POST  /games/{game_id}/rating
			r.Post("/rating", postrating.NewRatingPostHandler(logger, uc))

			r.Route("/comments", func(r chi.Router) {
				// GET  /games/{game_id}/comments?limit=&offset=
				r.Get("/", listcomments.NewListCommentsHandler(logger, uc))
				// POST /games/{game_id}/comments
				r.Post("/", addcomment.NewAddCommentHandler(logger, uc))
			})
		})
	})
	//router.Get("/swagger/*", httpSwagger.WrapHandler)
	//router.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:         cnfg.HttpInfo.Port,
		Handler:      router,
		ReadTimeout:  cnfg.HttpInfo.Timeout,
		WriteTimeout: cnfg.HttpInfo.Timeout,
		IdleTimeout:  cnfg.HttpInfo.IdleTimeout,
	}

	return server
}
