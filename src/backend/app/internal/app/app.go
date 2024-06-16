package app

import (
	"business-logic/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"ppo/app/internal/config"
	"ppo/internal/interactors/user_activity_field"
	"ppo/internal/services/activity_field"
	"ppo/internal/services/auth"
	"ppo/internal/services/company"
	"ppo/internal/services/contact"
	"ppo/internal/services/fin_report"
	"ppo/internal/services/skill"
	"ppo/internal/services/user"
	"ppo/internal/services/user_skill"
	"ppo/internal/storage/postgres"
	"ppo/pkg/base"
)

type App struct {
	Logger       *zap.SugaredLogger
	AuthSvc      domain.IAuthService
	UserSvc      domain.IUserService
	FinSvc       domain.IFinancialReportService
	ConSvc       domain.IContactsService
	SkillSvc     domain.ISkillService
	UserSkillSvc domain.IUserSkillService
	ActFieldSvc  domain.IActivityFieldService
	CompSvc      domain.ICompanyService
	Interactor   domain.IInteractor
	Config       config.Config
}

func NewApp(db *pgxpool.Pool, cfg *config.Config, logger *zap.SugaredLogger) *App {
	authRepo := postgres.NewAuthRepository(db)
	userRepo := postgres.NewUserRepository(db)
	finRepo := postgres.NewFinReportRepository(db)
	conRepo := postgres.NewContactRepository(db)
	skillRepo := postgres.NewSkillRepository(db)
	userSkillRepo := postgres.NewUserSkillRepository(db)
	actFieldRepo := postgres.NewActivityFieldRepository(db)
	compRepo := postgres.NewCompanyRepository(db)

	crypto := base.NewHashCrypto()

	authSvc := auth.NewService(authRepo, crypto, cfg.Server.JwtKey)
	userSvc := user.NewService(userRepo, compRepo, actFieldRepo)
	finSvc := fin_report.NewService(finRepo)
	conSvc := contact.NewService(conRepo)
	skillSvc := skill.NewService(skillRepo)
	userSkillSvc := user_skill.NewService(userSkillRepo, userRepo, skillRepo)
	actFieldSvc := activity_field.NewService(actFieldRepo, compRepo)
	compSvc := company.NewService(compRepo)
	interactor := user_activity_field.NewInteractor(userSvc, actFieldSvc, compSvc, finSvc)

	return &App{
		Logger:       logger,
		AuthSvc:      authSvc,
		UserSvc:      userSvc,
		FinSvc:       finSvc,
		ConSvc:       conSvc,
		SkillSvc:     skillSvc,
		UserSkillSvc: userSkillSvc,
		ActFieldSvc:  actFieldSvc,
		CompSvc:      compSvc,
		Interactor:   interactor,
		Config:       *cfg,
	}
}
