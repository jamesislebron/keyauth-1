package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/infraboard/mcube/http/middleware/accesslog"
	"github.com/infraboard/mcube/http/middleware/cors"
	"github.com/infraboard/mcube/http/middleware/recovery"
	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/http/router/httprouter"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"github.com/infraboard/keyauth/conf"
	"github.com/infraboard/keyauth/pkg"
	"github.com/infraboard/keyauth/pkg/endpoint"
	"github.com/infraboard/keyauth/pkg/micro"
	"github.com/infraboard/keyauth/pkg/token"
	"github.com/infraboard/keyauth/pkg/user/types"
	"github.com/infraboard/keyauth/version"
)

// NewHTTPService 构建函数
func NewHTTPService() *HTTPService {
	r := httprouter.New()
	r.Use(recovery.NewWithLogger(zap.L().Named("Recovery")))
	r.Use(accesslog.NewWithLogger(zap.L().Named("AccessLog")))
	r.Use(cors.AllowAll())
	r.EnableAPIRoot()
	r.SetAuther(pkg.NewInternalAuther())
	r.Auth(true)
	server := &http.Server{
		ReadHeaderTimeout: 20 * time.Second,
		// ReadTimeout:       20 * time.Second,
		// WriteTimeout:      25 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Addr:           conf.C().App.Addr(),
		Handler:        r,
	}

	return &HTTPService{
		r:      r,
		server: server,
		l:      zap.L().Named("API"),
		c:      conf.C(),
	}
}

// HTTPService http服务
type HTTPService struct {
	r      router.Router
	l      logger.Logger
	c      *conf.Config
	server *http.Server
}

// Start 启动服务
func (s *HTTPService) Start() error {
	// 装置子服务路由
	if err := pkg.InitV1HTTPAPI(s.c.App.Name, s.r); err != nil {
		return err
	}

	// 注册服务
	s.l.Info("start registry endpoints ...")
	if err := s.RegistryEndpoints(); err != nil {
		s.l.Warnf("registry endpoints error, %s", err)
	}
	s.l.Infof("service endpoints registry success: \n%s", s.r.GetEndpoints())

	// 启动HTTP服务
	s.l.Infof("服务启动成功, 监听地址: %s", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			s.l.Info("service is stopped")
		}
		return fmt.Errorf("start service error, %s", err.Error())
	}
	return nil
}

// Stop 停止server
func (s *HTTPService) Stop() error {
	s.l.Info("start graceful shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 优雅关闭HTTP服务
	if err := s.server.Shutdown(ctx); err != nil {
		s.l.Errorf("graceful shutdown timeout, force exit")
	}

	return nil
}

// RegistryEndpoints 注册条目
func (s *HTTPService) RegistryEndpoints() error {
	if pkg.Micro == nil {
		return fmt.Errorf("dependence micro service is nil")
	}

	desc := micro.NewDescriptServiceRequest()
	desc.Name = version.ServiceName
	svr, err := pkg.Micro.DescribeService(desc)
	if err != nil {
		return err
	}

	if pkg.Endpoint == nil {
		return fmt.Errorf("dependence endpoint service is nil")
	}

	tk := token.NewDefaultToken()
	tk.AccessToken = svr.AccessToken
	tk.RefreshToken = svr.RefreshToken
	tk.UserType = types.ServiceAccount
	tk.Account = svr.Name
	req := endpoint.NewRegistryRequest(version.Short(), s.r.GetEndpoints().Items)
	req.WithToken(tk)
	return pkg.Endpoint.Registry(req)
}
