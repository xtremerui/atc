package auth

const BuildContextKey = "build"
const PipelineContextKey = "pipeline"

const TokenTypeBearer = "Bearer"
const userIdClaimKey = "user_id"
const teamsClaimKey = "teams"
const isAdminClaimKey = "is_admin"
const csrfTokenClaimKey = "csrf"

const AuthCookieName = "ATC-Authorization"
const CSRFRequiredKey = "CSRFRequired"
const CSRFHeaderName = "X-Csrf-Token"

const isAuthenticatedKey = "isAuthenticated"
const userIdKey = "userId"
const teamsKey = "teams"
const isAdminKey = "isAdmin"
const isSystemKey = "isSystem"
const csrfTokenKey = "csrfToken"
