package server

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"antd-gin-admin-backend/internal/api/v1/scan"
	"antd-gin-admin-backend/internal/api/v1/system"
	"antd-gin-admin-backend/internal/bootstrap"
	"antd-gin-admin-backend/internal/config"
	"antd-gin-admin-backend/internal/middleware"
	"antd-gin-admin-backend/internal/model/entity"
	repoDB "antd-gin-admin-backend/internal/repository/db"
	repoInterfaces "antd-gin-admin-backend/internal/repository/interfaces"
	repoMock "antd-gin-admin-backend/internal/repository/mock"
	assetsvc "antd-gin-admin-backend/internal/service/asset"
	authsvc "antd-gin-admin-backend/internal/service/auth"
	detectionrulesvc "antd-gin-admin-backend/internal/service/detectionrule"
	scanjobsvc "antd-gin-admin-backend/internal/service/scanjob"
	cachesvc "antd-gin-admin-backend/internal/service/cache"
	dataManageSvc "antd-gin-admin-backend/internal/service/data_manage"
	dataScopeSvc "antd-gin-admin-backend/internal/service/data_scope"
	deptsvc "antd-gin-admin-backend/internal/service/dept"
	serviceInterfaces "antd-gin-admin-backend/internal/service/interfaces"
	menusvc "antd-gin-admin-backend/internal/service/menu"
	monitorsvc "antd-gin-admin-backend/internal/service/monitor"
	operationlogsvc "antd-gin-admin-backend/internal/service/operation_log"
	permissionsvc "antd-gin-admin-backend/internal/service/permission"
	profilesvc "antd-gin-admin-backend/internal/service/profile"
	rolesvc "antd-gin-admin-backend/internal/service/role"
	roleDeptSvc "antd-gin-admin-backend/internal/service/role_dept"
	roleMenuSvc "antd-gin-admin-backend/internal/service/role_menu"
	usersvc "antd-gin-admin-backend/internal/service/user"
	userRoleSvc "antd-gin-admin-backend/internal/service/user_role"
	"antd-gin-admin-backend/pkg/database"
	"antd-gin-admin-backend/pkg/logger"
	redisclient "antd-gin-admin-backend/pkg/redis"
)

// Run starts the HTTP server with baseline middleware and routes.
func Run(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("config is required")
	}

	logger.Init(cfg.Log.Level, cfg.Log.Format, cfg.Log.RetentionDays, cfg.Log.FilePath, cfg.Log.MaxSizeMB)

	if cfg.Server.Mode != "" {
		gin.SetMode(cfg.Server.Mode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.Trace())
	engine.Use(middleware.RequestLogger())
	engine.Use(middleware.CORS())

	var db *gorm.DB
	var err error
	var redisClient *redisclient.Client
	if !cfg.Auth.UseMock {
		db, err = database.Connect(cfg)
		if err != nil {
			return err
		}
		if err := bootstrap.AutoMigrate(db); err != nil {
			return err
		}
		if err := bootstrap.ApplySchemaComments(db); err != nil {
			return err
		}
		if err := bootstrap.EnsureAdminUser(db, cfg); err != nil {
			return err
		}
		if err := bootstrap.EnsureRBACSeed(db, cfg); err != nil {
			return err
		}
		// Initialize Redis client for monitoring and caching (optional).
		redisClient, err = redisclient.NewClient(cfg)
		if err != nil {
			logger.Warn(nil, "redis init failed, monitor features disabled", "error", err.Error())
		}
	}

	var authRepo repoInterfaces.AuthRepository
	var userRepo repoInterfaces.UserRepository
	var roleRepo repoInterfaces.RoleRepository
	var menuRepo repoInterfaces.MenuRepository
	var deptRepo repoInterfaces.DeptRepository
	var userRoleRepo repoInterfaces.UserRoleRepository
	var roleMenuRepo repoInterfaces.RoleMenuRepository
	var roleDeptRepo repoInterfaces.RoleDeptRepository
	var operationLogRepo repoInterfaces.OperationLogRepository
	var dbMetaRepo repoInterfaces.DBMetaRepository
	var assetRepo repoInterfaces.AssetRepository
	var detectionRuleRepo repoInterfaces.DetectionRuleRepository
	var scanJobRepo repoInterfaces.ScanJobRepository
	if cfg.Auth.UseMock {
		authRepo = repoMock.NewAuthMockRepository()
	} else {
		authRepo = repoDB.NewAuthRepository(db)
		userRepo = repoDB.NewUserRepository(db)
		roleRepo = repoDB.NewRoleRepository(db)
		menuRepo = repoDB.NewMenuRepository(db)
		deptRepo = repoDB.NewDeptRepository(db)
		userRoleRepo = repoDB.NewUserRoleRepository(db)
		roleMenuRepo = repoDB.NewRoleMenuRepository(db)
		roleDeptRepo = repoDB.NewRoleDeptRepository(db)
		operationLogRepo = repoDB.NewOperationLogRepository(db)
		dbMetaRepo = repoDB.NewDBMetaRepository(db)
		assetRepo = repoDB.NewAssetRepository(db)
		detectionRuleRepo = repoDB.NewDetectionRuleRepository(db)
		scanJobRepo = repoDB.NewScanJobRepository(db)
	}
	authSvc := authsvc.New(authRepo, cfg.JWT.Secret, cfg.JWT.Issuer, cfg.JWT.ExpireTime)
	var userSvc serviceInterfaces.UserService
	if userRepo != nil {
		userSvc = usersvc.New(userRepo)
	}
	var roleSvc serviceInterfaces.RoleService
	if roleRepo != nil {
		roleSvc = rolesvc.New(roleRepo)
	}
	var menuSvc serviceInterfaces.MenuService
	if menuRepo != nil {
		menuSvc = menusvc.New(menuRepo)
	}
	var deptSvc serviceInterfaces.DeptService
	if deptRepo != nil {
		deptSvc = deptsvc.New(deptRepo)
	}
	var userRoleSvcVar serviceInterfaces.UserRoleService
	if userRoleRepo != nil {
		userRoleSvcVar = userRoleSvc.New(userRoleRepo)
	}
	var permissionSvcVar serviceInterfaces.PermissionService
	if userRoleRepo != nil && roleMenuRepo != nil {
		permissionSvcVar = permissionsvc.New(userRoleRepo, roleMenuRepo)
	}
	var dataScopeSvcVar serviceInterfaces.DataScopeService
	if userRepo != nil && userRoleRepo != nil {
		dataScopeSvcVar = dataScopeSvc.New(userRepo, userRoleRepo, roleDeptRepo, deptRepo)
	}
	var profileSvc serviceInterfaces.ProfileService
	if userSvc != nil {
		profileSvc = profilesvc.New(userSvc, deptSvc, userRoleSvcVar, permissionSvcVar)
	}
	var roleMenuSvcVar serviceInterfaces.RoleMenuService
	if roleMenuRepo != nil {
		roleMenuSvcVar = roleMenuSvc.New(roleMenuRepo)
	}
	var roleDeptSvcVar serviceInterfaces.RoleDeptService
	if roleDeptRepo != nil {
		roleDeptSvcVar = roleDeptSvc.New(roleDeptRepo)
	}
	var operationLogSvcVar serviceInterfaces.OperationLogService
	if operationLogRepo != nil {
		operationLogSvcVar = operationlogsvc.New(operationLogRepo)
	}
	var dataManageSvcVar serviceInterfaces.DataManageService
	if dbMetaRepo != nil {
		dataManageSvcVar = dataManageSvc.New(dbMetaRepo, cfg.Database.Username, cfg.Database.Password)
	}
	var monitorSvcVar serviceInterfaces.MonitorService
	if redisClient != nil || db != nil {
		// Online user considered online if active within last 10 minutes.
		monitorSvcVar = monitorsvc.New(redisClient, 10*time.Minute, db, time.Duration(cfg.JWT.ExpireTime)*time.Second)
	}
	cacheSvcVar := cachesvc.New(redisClient)
	authMiddleware := middleware.Auth(cfg.JWT.Secret, monitorSvcVar)

	var assetSvcVar serviceInterfaces.AssetService
	if assetRepo != nil {
		assetSvcVar = assetsvc.New(assetRepo)
	}
	var detectionRuleSvcVar serviceInterfaces.DetectionRuleService
	if detectionRuleRepo != nil {
		// Default embedded fixture source until Compose ships nuclei-templates mount.
		ruleSource := &repoMock.StaticDetectionRuleSource{
			Rules: []entity.DetectionRule{
				{RuleCode: "exposed-panels", DisplayName: "Exposed Admin Panels"},
				{RuleCode: "http-missing-security-headers", DisplayName: "Missing Security Headers"},
				{RuleCode: "tech-detect", DisplayName: "Technology Detection"},
			},
		}
		detectionRuleSvcVar = detectionrulesvc.New(detectionRuleRepo, ruleSource)
	}
	var scanJobSvcVar serviceInterfaces.ScanJobService
	if scanJobRepo != nil && assetRepo != nil {
		var queue repoInterfaces.ScanJobQueue = repoMock.NewMemoryScanJobQueue()
		if redisClient != nil {
			queue = scanjobsvc.NewRedisQueue(redisClient)
		}
		scanJobSvcVar = scanjobsvc.New(scanJobRepo, assetRepo, queue)
	}

	api := engine.Group("/api/v1")
	{
		system.RegisterRoutes(api, authSvc, authMiddleware, permissionSvcVar, dataScopeSvcVar, profileSvc, userSvc, roleSvc, menuSvc, deptSvc, userRoleSvcVar, roleMenuSvcVar, roleDeptSvcVar, operationLogSvcVar, monitorSvcVar, dataManageSvcVar, cacheSvcVar)
		scan.RegisterRoutes(api, authMiddleware, permissionSvcVar, dataScopeSvcVar, assetSvcVar, detectionRuleSvcVar, scanJobSvcVar)
	}

	port := cfg.Server.Port
	if envPort := os.Getenv("ANTD_GIN_HTTP_PORT"); envPort != "" {
		port = envPort
	}
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	return engine.Run(addr)
}
