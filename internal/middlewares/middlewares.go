package middlewares

// IMiddleware middleware interface
type IMiddleware interface {
	Setup()
}

// Middlewares contains multiple middleware
type Middlewares []IMiddleware

// NewMiddlewares creates new middlewares
// Register the middleware that should be applied directly (globally)
func NewMiddlewares(
	coreMiddleware CoreMiddleware,
	corsMiddleware CorsMiddleware,
	zapMiddleware ZapMiddleware,
	authMiddleware AuthMiddleware,
	casbinMiddleware CasbinMiddleware,
	prometheusMiddleware PrometheusMiddleware,
) Middlewares {
	return Middlewares{
		coreMiddleware,
		zapMiddleware,
		corsMiddleware,
		authMiddleware,
		casbinMiddleware,
		prometheusMiddleware,
	}
}

// Setup sets up middlewares
func (a Middlewares) Setup() {
	for _, middleware := range a {
		middleware.Setup()
	}
}

func isIgnorePath(path string, prefixes ...string) bool {
	pathLen := len(path)

	for _, p := range prefixes {
		if pl := len(p); pathLen >= pl && path[:pl] == p {
			return true
		}
	}

	return false
}
