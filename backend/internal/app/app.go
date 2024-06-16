package app

import (
	"ppo/domain"
	"ppo/internal/config"
	"ppo/internal/interactors/user_activity_field"
	"ppo/internal/services/activity_field"
	"ppo/internal/services/auth"
	"ppo/internal/services/company"
	"ppo/internal/services/contact"
	"ppo/internal/services/fin_report"
	"ppo/internal/services/review"
	"ppo/internal/services/skill"
	"ppo/internal/services/user"
	"ppo/internal/services/user_skill"
	"ppo/internal/storage/postgres"
	"ppo/pkg/base"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	AuthSvc      domain.IAuthService
	UserSvc      domain.IUserService
	FinSvc       domain.IFinancialReportService
	ConSvc       domain.IContactsService
	SkillSvc     domain.ISkillService
	UserSkillSvc domain.IUserSkillService
	ActFieldSvc  domain.IActivityFieldService
	CompSvc      domain.ICompanyService
	RevSvc       domain.IReviewService
	Interactor   domain.IInteractor
	Config       config.Config
}

func NewApp(db *pgxpool.Pool, cfg *config.Config) *App {
	authRepo := postgres.NewAuthRepository(db)
	userRepo := postgres.NewUserRepository(db)
	finRepo := postgres.NewFinReportRepository(db)
	conRepo := postgres.NewContactRepository(db)
	skillRepo := postgres.NewSkillRepository(db)
	userSkillRepo := postgres.NewUserSkillRepository(db)
	actFieldRepo := postgres.NewActivityFieldRepository(db)
	compRepo := postgres.NewCompanyRepository(db)
	revRepo := postgres.NewReviewRepository(db)

	crypto := base.NewHashCrypto()

	authSvc := auth.NewService(authRepo, crypto, cfg.JwtKey)
	userSvc := user.NewService(userRepo, compRepo, actFieldRepo)
	finSvc := fin_report.NewService(finRepo)
	conSvc := contact.NewService(conRepo)
	skillSvc := skill.NewService(skillRepo)
	userSkillSvc := user_skill.NewService(userSkillRepo, userRepo, skillRepo)
	actFieldSvc := activity_field.NewService(actFieldRepo, compRepo)
	compSvc := company.NewService(compRepo, actFieldRepo)
	revSvc := review.NewService(revRepo)
	interactor := user_activity_field.NewInteractor(userSvc, actFieldSvc, compSvc, finSvc)

	return &App{
		AuthSvc:      authSvc,
		UserSvc:      userSvc,
		FinSvc:       finSvc,
		ConSvc:       conSvc,
		SkillSvc:     skillSvc,
		UserSkillSvc: userSkillSvc,
		ActFieldSvc:  actFieldSvc,
		CompSvc:      compSvc,
		RevSvc:       revSvc,
		Interactor:   interactor,
		Config:       *cfg,
	}
}
