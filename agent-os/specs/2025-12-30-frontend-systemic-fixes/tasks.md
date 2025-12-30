# Frontend Systemic Fixes - Task List

**Spec:** `frontend-systemic-fixes-2025-12-30`
**Created:** 2025-12-30
**Status:** IN_PROGRESS

---

## Phase 1: Unified Auth State (P0) - CRITICAL

**Agent:** auth-system-agent
**Status:** COMPLETE
**Start Time:** 2025-12-30 23:30:00
**End Time:** 2025-12-30 23:45:00
**Duration:** 15 minutes
**Duration Estimate:** 2 hours

### Tasks

- [x] **TASK-001:** Create AuthManager class
  - File: `/frontend/src/lib/auth/AuthManager.ts`
  - Description: Singleton class managing all auth state
  - Features: User, tokens, permissions, tenant context
  - Status: COMPLETE
  - Notes: Created with full singleton pattern, in-memory token storage, session management

- [x] **TASK-002:** Refactor auth store to use AuthManager
  - File: `/frontend/src/lib/stores/auth.ts`
  - Description: Simplify store to wrap AuthManager
  - Remove: Direct localStorage access
  - Status: COMPLETE
  - Notes: Removed persisted-store, wrapped AuthManager with reactive interface

- [x] **TASK-003:** Implement token management
  - Description: Access token in-memory, refresh in httpOnly
  - Add: Token expiry validation
  - Add: Automatic token refresh
  - Status: COMPLETE
  - Notes: Token expiry checked before use, automatic refresh on 401

- [x] **TASK-004:** Add unified logout cleanup
  - Description: Clear ALL state on logout
  - Clear: AuthManager, all stores, WebSocket, localStorage
  - Status: COMPLETE
  - Notes: Comprehensive cleanup preserving only user preferences

- [x] **TASK-005:** Add permission loading on login
  - Description: Load permissions immediately after login
  - Integration: permissionStore.loadPermissions()
  - Status: COMPLETE
  - Notes: Permissions loaded atomically during login

- [x] **TASK-006:** Implement session timeout
  - Description: Auto-logout after inactivity
  - Duration: Use VITE_SESSION_TIMEOUT env var
  - Status: COMPLETE
  - Notes: Session timeout from env (default 60 min), auto-logout on expiry

- [x] **TASK-007:** Update API client to use AuthManager
  - File: `/frontend/src/lib/api/client.ts`
  - Description: Replace localStorage with AuthManager
  - Remove: Direct token access
  - Status: COMPLETE
  - Notes: All token/tenant operations now use AuthManager, 401 refresh integrated

- [x] **TASK-008:** Write tests for auth flow
  - Tests: Login, logout, token refresh, session timeout
  - Coverage: > 90%
  - Status: COMPLETE
  - Notes: 16 comprehensive tests written, all passing (100% pass rate)

**Exit Criteria:**
- ✅ All tests pass (16/16 passing)
- ✅ Login stores everything atomically (verified in tests)
- ✅ Logout clears everything (verified in tests)
- ✅ No auth state desynchronization possible (AuthManager is single source of truth)
- ✅ Token expiry validated before use
- ✅ Permissions load on login
- ✅ Session timeout implemented

---

## Phase 2: Single Route Guard (P0) - CRITICAL

**Agent:** route-guard-agent
**Status:** PENDING
**Dependencies:** Phase 1 complete
**Duration Estimate:** 1 hour

### Tasks

- [ ] **TASK-101:** Remove auth logic from root layout
  - File: `/frontend/src/routes/+layout.svelte`
  - Remove: onMount auth check
  - Remove: loadUser call
  - Keep: Basic HTML structure only
  - Status: PENDING

- [ ] **TASK-102:** Implement single auth guard in app layout
  - File: `/frontend/src/routes/app/+layout.svelte`
  - Add: Single onMount guard
  - Add: Loading state
  - Add: Proper redirect logic
  - Status: PENDING

- [ ] **TASK-103:** Remove redirect from login page
  - File: `/frontend/src/routes/auth/login/+page.svelte`
  - Remove: onMount redirect check
  - Reason: Should only redirect after successful login action
  - Status: PENDING

- [ ] **TASK-104:** Add loading component
  - File: `/frontend/src/lib/components/LoadingScreen.svelte`
  - Description: Full-screen loading indicator
  - Status: PENDING

- [ ] **TASK-105:** Write tests for route guards
  - Test: Unauthenticated redirect
  - Test: Authenticated access
  - Test: Loading state display
  - Status: PENDING

**Exit Criteria:**
- ✅ Only one auth check per navigation
- ✅ No race conditions
- ✅ No flash of wrong content

---

## Phase 3: AdminSidebar Simplification (P1)

**Agent:** admin-sidebar-agent
**Status:** PENDING
**Dependencies:** Phase 1 complete
**Duration Estimate:** 1 hour

### Tasks

- [ ] **TASK-201:** Remove reactive accordion logic
  - File: `/frontend/src/lib/components/admin/AdminSidebar.svelte`
  - Remove: Complex $: reactive statements
  - Remove: manualOverride flag
  - Remove: previousPath tracking
  - Status: PENDING

- [ ] **TASK-202:** Implement explicit accordion functions
  - Add: expandSection(title: string)
  - Add: collapseSection(title: string)
  - Add: toggleSection(title: string)
  - Logic: Simple, predictable, testable
  - Status: PENDING

- [ ] **TASK-203:** Add sidebar state to UI store
  - File: `/frontend/src/lib/stores/ui.ts`
  - Add: expandedSections: string[]
  - Add: Make store persisted
  - Status: PENDING

- [ ] **TASK-204:** Auto-expand active section on mount
  - Trigger: onMount only (not reactive)
  - Logic: Check current route, expand matching section
  - Status: PENDING

- [ ] **TASK-205:** Write tests for accordion
  - Test: Expand/collapse works
  - Test: Active section auto-expands
  - Test: Preferences persist
  - Status: PENDING

**Exit Criteria:**
- ✅ No reactive loop warnings
- ✅ Accordion behaves predictably
- ✅ User preferences respected

---

## Phase 4: Tenant ID & API Client (P1)

**Agent:** api-client-agent
**Status:** PENDING
**Dependencies:** Phase 1 complete
**Duration Estimate:** 1 hour

### Tasks

- [ ] **TASK-301:** Remove tenant ID fallback layers
  - File: `/frontend/src/lib/api/client.ts`
  - Keep: Single source from AuthManager
  - Remove: All fallbacks, hardcoded UUIDs
  - Fail fast: Throw error if tenant ID missing
  - Status: PENDING

- [ ] **TASK-302:** Merge admin API client
  - Delete: `/frontend/src/lib/api/admin.ts`
  - Merge: All admin functions into client.ts
  - Reason: Eliminate duplication
  - Status: PENDING

- [ ] **TASK-303:** Add token expiry validation
  - Add: JWT decode and expiry check
  - Check: Before every request
  - Trigger: Auto-refresh if expired
  - Status: PENDING

- [ ] **TASK-304:** Fix CORS proxy configuration
  - File: `/frontend/vite.config.ts`
  - Add: Proxy for analytics service (8088)
  - Add: Proxy for ML service (8004)
  - Status: PENDING

- [ ] **TASK-305:** Set proper .env.example defaults
  - File: `/frontend/.env.example`
  - Set: VITE_DEFAULT_TENANT_ID with valid UUID
  - Document: How to get tenant ID
  - Status: PENDING

- [ ] **TASK-306:** Write API client tests
  - Test: Tenant ID required
  - Test: Token expiry checked
  - Test: CORS works
  - Status: PENDING

**Exit Criteria:**
- ✅ No silent fallbacks
- ✅ Single API client
- ✅ Token validation works

---

## Phase 5: Error Handling (P1)

**Agent:** error-handling-agent
**Status:** PENDING
**Dependencies:** Phase 1 complete
**Duration Estimate:** 1 hour

### Tasks

- [ ] **TASK-401:** Create global error boundary
  - File: `/frontend/src/lib/components/ErrorBoundary.svelte`
  - Catch: Unhandled errors
  - Display: User-friendly message
  - Option: Reload or contact support
  - Status: PENDING

- [ ] **TASK-402:** Enhance app error page
  - File: `/frontend/src/routes/app/+error.svelte`
  - Add: Better error messages
  - Add: Error reporting
  - Status: PENDING

- [ ] **TASK-403:** Fix API error interceptor
  - File: `/frontend/src/lib/api/client.ts`
  - Add: Automatic redirect to login on 401
  - Add: Retry logic with max attempts (3)
  - Add: User notification for 500 errors
  - Status: PENDING

- [ ] **TASK-404:** Create error utilities
  - File: `/frontend/src/lib/utils/errors.ts`
  - Add: Error type detection
  - Add: User-friendly message mapping
  - Add: Error logging helpers
  - Status: PENDING

- [ ] **TASK-405:** Standardize error responses
  - Description: Document expected error format
  - Format: `{ error: string, message: string, details?: any }`
  - Update: All API error handlers
  - Status: PENDING

- [ ] **TASK-406:** Write error handling tests
  - Test: 401 redirect
  - Test: 500 error UI
  - Test: Retry mechanism
  - Status: PENDING

**Exit Criteria:**
- ✅ Errors don't crash app
- ✅ 401 redirects work
- ✅ User sees helpful messages

---

## Phase 6: Integration & Testing (P1)

**Agent:** integration-agent
**Status:** PENDING
**Dependencies:** All previous phases complete
**Duration Estimate:** 1 hour

### Tasks

- [ ] **TASK-501:** End-to-end testing
  - Test: Login → Navigate → Logout flow
  - Test: Multi-tab behavior
  - Test: Token refresh during navigation
  - Test: Permission checks
  - Test: All original 23 issues resolved
  - Status: PENDING

- [ ] **TASK-502:** Clean up duplicate code
  - Search: Duplicate auth checks
  - Search: Duplicate API logic
  - Remove: All found duplications
  - Status: PENDING

- [ ] **TASK-503:** Add missing cleanup handlers
  - Add: WebSocket cleanup on logout
  - Add: Store reset on tenant switch
  - Add: Timer cleanup on unmount
  - Status: PENDING

- [ ] **TASK-504:** Update documentation
  - Update: README with new auth flow
  - Update: Contributing guide
  - Add: Migration guide for breaking changes
  - Status: PENDING

- [ ] **TASK-505:** Verify all issues resolved
  - Check: All 23 original issues
  - Document: Resolution for each
  - Status: PENDING

- [ ] **TASK-506:** Performance testing
  - Test: Page load times
  - Test: API response times
  - Test: Memory leaks
  - Status: PENDING

- [ ] **TASK-507:** Create rollback plan
  - Document: Rollback procedure
  - Test: Rollback works
  - Status: PENDING

**Exit Criteria:**
- ✅ All acceptance criteria met
- ✅ All tests passing
- ✅ Zero regression
- ✅ Documentation complete

---

## Summary

**Total Tasks:** 37
**Completed:** 0
**In Progress:** 0
**Pending:** 37
**Blocked:** 0

**Overall Progress:** 0%

**Critical Path:**
Phase 1 → (Phase 2 || Phase 3) → (Phase 4 || Phase 5) → Phase 6

**Estimated Total Duration:** 6-7 hours
**Actual Duration:** TBD

---

## Status Legend

- **PENDING:** Not started
- **IN_PROGRESS:** Currently being worked on
- **BLOCKED:** Waiting for dependency
- **COMPLETED:** Finished and tested
- **FAILED:** Encountered issues, needs retry

---

**Last Updated:** 2025-12-30
**Next Review:** After Phase 1 completion
