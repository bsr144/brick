package main

import (
	"brick/internal/adapters/presenters"
	"brick/internal/adapters/repositories"
	"brick/internal/adapters/rest/controllers"
	"brick/internal/adapters/rest/middlewares"
	"brick/internal/adapters/rest/routes"
	"brick/internal/config"
	"brick/internal/drivers/database"
	"brick/internal/drivers/logging"
	"brick/internal/drivers/rest"
	"brick/internal/usecases/transfer"
	"brick/internal/usecases/user"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

func main() {
	driverConfig := config.New()
	RESTServer := rest.NewGinServer(driverConfig.App)
	log := logging.NewLogger(driverConfig.Logging)
	sqlDB := database.NewSQLDatabase(driverConfig.Database, log)

	bootstrapRESTServer(&config.Bootstrap{
		RESTServer: RESTServer,
		Log:        log,
		SqlDB:      sqlDB,
		Driver:     driverConfig,
	})

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", driverConfig.App.Port),
		Handler: RESTServer,
	}

	err := httpServer.ListenAndServe()

	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to start server")
	}

	log.Println(fmt.Sprintf("Server listening on port %d", driverConfig.App.Port))
}

func bootstrapRESTServer(bootstrapConfig *config.Bootstrap) {
	// Init repositories
	userRepository := repositories.NewUserRepository(bootstrapConfig.SqlDB, bootstrapConfig.Log)
	credentialRepository := repositories.NewCredentialRepository(bootstrapConfig.SqlDB, bootstrapConfig.Log)
	transferRepository := repositories.NewTransferRepository(bootstrapConfig.SqlDB, bootstrapConfig.Log)
	recipientAccountRepository := repositories.NewRecipientAccountRepository(bootstrapConfig.SqlDB, bootstrapConfig.Log)

	// Init presenters
	userPresenter := presenters.NewUserPresenter(bootstrapConfig.Log)
	transferPresenter := presenters.NewTransferPresenter(bootstrapConfig.Log)

	// Init usecases
	userUsecase := user.NewUserUsecase(userRepository, credentialRepository, userPresenter, bootstrapConfig.SqlDB, bootstrapConfig.Log)
	transferUsecase := transfer.NewTransferUsecase(transferRepository, recipientAccountRepository, transferPresenter, bootstrapConfig.SqlDB, bootstrapConfig.Log)

	// Init controllers
	transferController := controllers.NewTransferController(transferUsecase, bootstrapConfig.Log)
	userController := controllers.NewUserController(userUsecase, bootstrapConfig.Log)

	// Init middlewares
	authMiddleware := middlewares.NewAuthMiddleware(bootstrapConfig.Log)

	// Init app routes
	restRoute := routes.NewRESTRoute(userController, transferController, authMiddleware)

	// Setup routes
	restRoute.SetupRoutes(bootstrapConfig.RESTServer)
}
