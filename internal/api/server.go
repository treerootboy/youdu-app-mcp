package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"unicode"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yourusername/youdu-app-mcp/internal/adapter"
	"github.com/yourusername/youdu-app-mcp/internal/config"
)

// Server represents the HTTP API server
type Server struct {
	router       chi.Router
	adapter      *adapter.Adapter
	config       *config.Config
	tokenEnabled bool
}

// New creates a new API server
func New(cfg *config.Config) (*Server, error) {
	// 创建 adapter
	adp, err := adapter.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create adapter: %w", err)
	}

	// 创建路由器
	r := chi.NewRouter()

	// 添加中间件
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(corsMiddleware)
	r.Use(jsonContentTypeMiddleware)

	// 检查是否启用 token 认证
	tokenEnabled := cfg.TokenManager != nil && cfg.TokenManager.Count() > 0

	s := &Server{
		router:       r,
		adapter:      adp,
		config:       cfg,
		tokenEnabled: tokenEnabled,
	}

	// 添加 token 认证中间件（如果启用）
	if tokenEnabled {
		s.router.Use(s.tokenAuthMiddleware)
	}

	// 自动注册所有 adapter 方法为 HTTP endpoint
	if err := s.registerRoutes(); err != nil {
		return nil, fmt.Errorf("failed to register routes: %w", err)
	}

	// 添加健康检查和元信息端点
	s.registerMetaRoutes()

	return s, nil
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	fmt.Printf("🚀 YouDu API Server 启动在 %s\n", addr)
	fmt.Println("📖 API 文档: GET /api/v1/endpoints")
	fmt.Println("💚 健康检查: GET /health")
	if s.tokenEnabled {
		fmt.Println("🔒 Token 认证: 已启用")
		fmt.Printf("   当前有效 token 数量: %d\n", s.config.TokenManager.Count())
	} else {
		fmt.Println("⚠️  Token 认证: 未启用")
	}
	return http.ListenAndServe(addr, s.router)
}

// Close closes the server and releases resources
func (s *Server) Close() error {
	return s.adapter.Close()
}

// registerRoutes 使用反射自动注册所有 adapter 方法为 HTTP endpoint
func (s *Server) registerRoutes() error {
	adapterType := reflect.TypeOf(s.adapter)
	adapterValue := reflect.ValueOf(s.adapter)

	fmt.Println("\n📋 正在注册 API Endpoints:")

	// 遍历 adapter 的所有方法
	for i := 0; i < adapterType.NumMethod(); i++ {
		method := adapterType.Method(i)

		// 跳过非导出方法和特殊方法
		if !method.IsExported() || method.Name == "Close" || method.Name == "Context" || method.Name == "GetConfig" {
			continue
		}

		// 转换方法名为 snake_case 作为路径
		path := toSnakeCase(method.Name)

		// 获取方法类型信息
		methodType := method.Type

		// 验证方法签名: func(ctx context.Context, input InputType) (*OutputType, error)
		if methodType.NumIn() != 3 || methodType.NumOut() != 2 {
			continue
		}

		// 检查第一个参数是 context.Context
		contextType := reflect.TypeOf((*context.Context)(nil)).Elem()
		if !methodType.In(1).Implements(contextType) {
			continue
		}

		// 检查最后返回值是 error
		errorType := reflect.TypeOf((*error)(nil)).Elem()
		if !methodType.Out(1).Implements(errorType) {
			continue
		}

		// 获取输入和输出类型
		inputType := methodType.In(2)
		outputType := methodType.Out(0)

		// 注册路由
		if err := s.registerRoute(path, method, adapterValue, inputType, outputType); err != nil {
			return fmt.Errorf("failed to register route %s: %w", path, err)
		}

		fmt.Printf("  ✓ POST /api/v1/%s\n", path)
	}

	return nil
}

// registerRoute 注册单个路由
func (s *Server) registerRoute(path string, method reflect.Method, adapterValue reflect.Value, inputType, outputType reflect.Type) error {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// 创建输入实例
		input := reflect.New(inputType).Interface()

		// 解析 JSON 请求体
		if err := json.NewDecoder(r.Body).Decode(input); err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid JSON: %v", err))
			return
		}

		// 调用 adapter 方法
		results := method.Func.Call([]reflect.Value{
			adapterValue,
			reflect.ValueOf(r.Context()),
			reflect.ValueOf(input).Elem(),
		})

		// 检查错误
		if !results[1].IsNil() {
			err := results[1].Interface().(error)
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// 返回成功响应
		output := results[0].Interface()
		respondJSON(w, http.StatusOK, output)
	}

	// 注册 POST 路由
	s.router.Post(fmt.Sprintf("/api/v1/%s", path), handler)

	return nil
}

// registerMetaRoutes 注册元信息路由
func (s *Server) registerMetaRoutes() {
	// 健康检查
	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"status":  "ok",
			"service": "youdu-api",
			"version": "1.0.0",
		})
	})

	// API 列表
	s.router.Get("/api/v1/endpoints", func(w http.ResponseWriter, r *http.Request) {
		endpoints := s.listEndpoints()
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"endpoints": endpoints,
			"count":     len(endpoints),
		})
	})
}

// listEndpoints 列出所有可用的 API endpoint
func (s *Server) listEndpoints() []map[string]interface{} {
	adapterType := reflect.TypeOf(s.adapter)
	endpoints := []map[string]interface{}{}

	for i := 0; i < adapterType.NumMethod(); i++ {
		method := adapterType.Method(i)

		if !method.IsExported() || method.Name == "Close" || method.Name == "Context" || method.Name == "GetConfig" {
			continue
		}

		methodType := method.Type
		if methodType.NumIn() != 3 || methodType.NumOut() != 2 {
			continue
		}

		path := toSnakeCase(method.Name)
		inputType := methodType.In(2)
		outputType := methodType.Out(0)

		endpoints = append(endpoints, map[string]interface{}{
			"method":      "POST",
			"path":        fmt.Sprintf("/api/v1/%s", path),
			"name":        method.Name,
			"description": generateDescription(method.Name),
			"input_type":  inputType.String(),
			"output_type": outputType.String(),
		})
	}

	return endpoints
}

// toSnakeCase 将 PascalCase 转换为 snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

// generateDescription 从方法名生成可读描述
func generateDescription(methodName string) string {
	var words []string
	var currentWord strings.Builder

	for i, r := range methodName {
		if i > 0 && unicode.IsUpper(r) {
			words = append(words, currentWord.String())
			currentWord.Reset()
		}
		currentWord.WriteRune(unicode.ToLower(r))
	}
	words = append(words, currentWord.String())

	return strings.Join(words, " ")
}

// respondJSON 返回 JSON 响应
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError 返回错误响应
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]interface{}{
		"error":   true,
		"message": message,
	})
}

// corsMiddleware 添加 CORS 支持
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// jsonContentTypeMiddleware 设置 JSON Content-Type
func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

// tokenAuthMiddleware 验证 token
func (s *Server) tokenAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 跳过健康检查和 endpoints 列表
		if r.URL.Path == "/health" || r.URL.Path == "/api/v1/endpoints" {
			next.ServeHTTP(w, r)
			return
		}

		// 从 Authorization header 获取 token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondError(w, http.StatusUnauthorized, "缺少 Authorization header")
			return
		}

		// 支持两种格式: "Bearer <token>" 或 直接 "<token>"
		token := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// 验证 token
		if !s.config.TokenManager.Validate(token) {
			respondError(w, http.StatusUnauthorized, "无效的 token")
			return
		}

		next.ServeHTTP(w, r)
	})
}
