# Comply360 Mobile App - Testing Guide

**Comprehensive testing guide for the Comply360 mobile application**

---

## 📋 Table of Contents

1. [Manual Testing](#manual-testing)
2. [Automated Testing](#automated-testing)
3. [Test Scenarios](#test-scenarios)
4. [Platform-Specific Testing](#platform-specific-testing)
5. [Performance Testing](#performance-testing)
6. [Security Testing](#security-testing)

---

## 🧪 Manual Testing

### Pre-Testing Setup

**Required:**
- Physical iOS device (iPhone 11 or later)
- Physical Android device (API 23+)
- Backend API running on `http://localhost:8080` or configured URL
- Test user accounts created in the system

**Test Environment Variables:**
```env
API_URL=http://localhost:8080/api/v1
```

---

## 🔐 Authentication Testing

### Test Case 1: Email/Password Login

**Steps:**
1. Launch app
2. Enter valid email and password
3. Tap "Sign In"

**Expected:**
- ✅ Loading indicator appears
- ✅ Dashboard loads with user data
- ✅ Token stored in Redux persist
- ✅ No console errors

**Test Data:**
```
Email: test@comply360.com
Password: Test@123
```

---

### Test Case 2: Invalid Login

**Steps:**
1. Launch app
2. Enter invalid credentials
3. Tap "Sign In"

**Expected:**
- ✅ Error message displayed
- ✅ Form not cleared
- ✅ User remains on login screen
- ✅ Error is user-friendly

**Test Data:**
```
Email: invalid@test.com
Password: wrongpassword
```

---

### Test Case 3: Biometric Setup

**Steps:**
1. Login successfully
2. When prompted, tap "Enable [Biometric Type]"
3. Authenticate with biometric
4. Complete setup

**Expected:**
- ✅ Device biometric prompt appears
- ✅ Success message shown
- ✅ Credentials stored in Keychain
- ✅ Biometric flag set in Redux

**Platforms:**
- iOS: Face ID / Touch ID
- Android: Fingerprint

---

### Test Case 4: Biometric Login

**Steps:**
1. Logout from app
2. Tap biometric login button
3. Authenticate with biometric

**Expected:**
- ✅ Biometric prompt appears
- ✅ Credentials retrieved from Keychain
- ✅ Auto-login successful
- ✅ Dashboard loads

---

### Test Case 5: Forgot Password

**Steps:**
1. On login screen, tap "Forgot Password?"
2. Enter email address
3. Tap "Send Reset Instructions"

**Expected:**
- ✅ Success screen shown
- ✅ Email sent to user
- ✅ Instructions displayed
- ✅ Can return to login

**Test Data:**
```
Email: test@comply360.com
```

---

### Test Case 6: Remember Me

**Steps:**
1. Check "Remember me" checkbox
2. Login successfully
3. Close app completely
4. Reopen app

**Expected:**
- ✅ User auto-logged in
- ✅ Dashboard shown immediately
- ✅ Token valid

---

### Test Case 7: Logout

**Steps:**
1. Navigate to Profile tab
2. Tap "Logout"
3. Confirm logout

**Expected:**
- ✅ Confirmation dialog shown
- ✅ User logged out
- ✅ Login screen shown
- ✅ Redux state cleared
- ✅ Biometric credentials cleared (if enabled)

---

## 📊 Dashboard Testing

### Test Case 8: Dashboard Load

**Steps:**
1. Login successfully
2. Observe dashboard

**Expected:**
- ✅ Loading spinner shown initially
- ✅ 3 API calls made (registrations, commissions, documents)
- ✅ Stats displayed correctly
- ✅ Quick actions visible
- ✅ Summary sections shown
- ✅ No errors in console

---

### Test Case 9: Pull to Refresh

**Steps:**
1. On dashboard
2. Pull down to refresh

**Expected:**
- ✅ Refresh indicator shown
- ✅ All 3 APIs re-fetched
- ✅ Data updates
- ✅ Loading states correct

---

### Test Case 10: Quick Actions

**Steps:**
1. Tap each quick action button
   - New Registration
   - Upload Document
   - View Commissions
   - Profile

**Expected:**
- ✅ Navigation works (or console log for now)
- ✅ No crashes
- ✅ Smooth transitions

---

### Test Case 11: Stats Display

**Steps:**
1. View dashboard stats
2. Verify numbers match backend

**Expected:**
- ✅ Registration count correct
- ✅ Commission amount formatted (R X,XXX)
- ✅ Document counts correct
- ✅ Notification badge shows pending count

---

## 🔄 State Persistence Testing

### Test Case 12: Offline Persistence

**Steps:**
1. Login successfully
2. Force close app (swipe up / kill process)
3. Reopen app

**Expected:**
- ✅ User still logged in
- ✅ Dashboard loads immediately
- ✅ Token persisted
- ✅ User data persisted

---

### Test Case 13: Network Offline

**Steps:**
1. Login successfully
2. Turn off WiFi and cellular
3. Navigate around app

**Expected:**
- ✅ Cached data still displayed
- ✅ Error messages for new requests
- ✅ No crashes
- ✅ Graceful error handling

---

### Test Case 14: Network Recovery

**Steps:**
1. While offline, attempt API call
2. Turn network back on
3. Pull to refresh

**Expected:**
- ✅ API calls succeed
- ✅ Data refreshed
- ✅ Cache updated

---

## 🧩 Component Testing

### Test Case 15: Button Component

**Test all variants:**
```typescript
- Primary (purple)
- Secondary (gray)
- Outline (border)
- Ghost (transparent)
- Danger (red)
```

**Test all states:**
- Normal
- Loading
- Disabled
- With icon (left/right)

**Expected:**
- ✅ Correct styling for each variant
- ✅ Loading spinner shown when loading
- ✅ Disabled state prevents clicks
- ✅ Icons positioned correctly

---

### Test Case 16: Modal Component

**Steps:**
1. Trigger modal open
2. Interact with modal
3. Close via backdrop / close button / action button

**Expected:**
- ✅ Backdrop visible
- ✅ Modal centered
- ✅ Content scrollable if long
- ✅ Closes on backdrop tap (if enabled)
- ✅ Buttons work correctly

---

### Test Case 17: Toast Notifications

**Steps:**
1. Trigger success toast
2. Trigger error toast
3. Trigger warning toast
4. Trigger info toast

**Expected:**
- ✅ Toast slides from top
- ✅ Correct color for each type
- ✅ Auto-dismisses after 4s
- ✅ Can manually dismiss
- ✅ Multiple toasts queue properly

---

### Test Case 18: Error Boundary

**Steps:**
1. Trigger a React error (simulate)
2. Observe error screen

**Expected:**
- ✅ Error boundary catches error
- ✅ Fallback UI shown
- ✅ Error details shown (dev mode)
- ✅ "Try Again" button works
- ✅ No app crash

---

## 📱 Platform-Specific Testing

### iOS Testing

**Test On:**
- iPhone SE (small screen)
- iPhone 12/13/14 (notch)
- iPhone 14 Pro Max (large)
- iPad (tablet)

**iOS-Specific Checks:**
- ✅ Safe area respected (notch, home indicator)
- ✅ Face ID works correctly
- ✅ Touch ID works correctly (if device supports)
- ✅ Keyboard avoidance works
- ✅ Status bar color correct
- ✅ Navigation gestures work

---

### Android Testing

**Test On:**
- Small phone (5-6")
- Large phone (6.5"+)
- Tablet (10"+)
- Different Android versions (API 23-33)

**Android-Specific Checks:**
- ✅ Fingerprint works correctly
- ✅ Back button behavior correct
- ✅ Material design consistent
- ✅ Keyboard behavior correct
- ✅ Status bar color correct
- ✅ Permission requests work

---

## ⚡ Performance Testing

### Test Case 19: App Launch Time

**Steps:**
1. Force close app
2. Time app launch to dashboard

**Expected:**
- ✅ Cold start: < 3 seconds
- ✅ Warm start: < 1 second
- ✅ No white screen flash
- ✅ Smooth splash → dashboard transition

---

### Test Case 20: Memory Usage

**Steps:**
1. Use Android Studio Profiler / Xcode Instruments
2. Navigate through all screens
3. Observe memory usage

**Expected:**
- ✅ Memory usage < 150MB
- ✅ No memory leaks
- ✅ Garbage collection working

---

### Test Case 21: API Response Time

**Steps:**
1. Monitor network tab
2. Trigger API calls
3. Measure response times

**Expected:**
- ✅ Dashboard stats: < 1s
- ✅ Login: < 2s
- ✅ Document upload: < 5s
- ✅ Loading states shown

---

### Test Case 22: Smooth Scrolling

**Steps:**
1. Scroll long lists
2. Navigate between tabs
3. Open/close modals

**Expected:**
- ✅ 60 FPS scrolling
- ✅ No janky animations
- ✅ Smooth transitions

---

## 🔒 Security Testing

### Test Case 23: Secure Storage

**Steps:**
1. Login with biometric
2. Use device file explorer (root required)
3. Check Keychain/Keystore

**Expected:**
- ✅ Credentials encrypted
- ✅ Token not in plain text
- ✅ Keychain protected by biometric

---

### Test Case 24: Token Expiry

**Steps:**
1. Login successfully
2. Wait for token expiry (or manually expire)
3. Attempt API call

**Expected:**
- ✅ 401 response received
- ✅ Refresh token used automatically
- ✅ New token stored
- ✅ API call retried successfully

**OR if refresh fails:**
- ✅ User logged out
- ✅ Redirected to login

---

### Test Case 25: Input Validation

**Steps:**
1. Enter invalid data in forms
2. Submit forms

**Expected:**
- ✅ Validation errors shown
- ✅ No API call made with invalid data
- ✅ User-friendly error messages
- ✅ XSS prevented

**Test Cases:**
```
- Invalid email format
- Weak password
- Invalid phone number
- Invalid ID number
- SQL injection attempts
```

---

### Test Case 26: HTTPS Enforcement

**Steps:**
1. Check API calls in network inspector
2. Attempt HTTP calls (if possible)

**Expected:**
- ✅ All API calls use HTTPS (production)
- ✅ Certificate validation enabled
- ✅ No mixed content

---

## 🧪 Automated Testing

### Unit Tests

**Run Tests:**
```bash
npm test
```

**Test Coverage:**
```bash
npm run test:coverage
```

**Expected Coverage:**
- ✅ Utilities: > 90%
- ✅ Services: > 80%
- ✅ Redux slices: > 80%
- ✅ Components: > 70%

---

### Example Test Cases

**Validation Tests:**
```typescript
describe('validateEmail', () => {
  it('should validate correct email', () => {
    expect(validateEmail('test@example.com')).toBe(true);
  });

  it('should reject invalid email', () => {
    expect(validateEmail('invalid')).toBe(false);
  });
});
```

**Formatting Tests:**
```typescript
describe('formatCurrency', () => {
  it('should format currency correctly', () => {
    expect(formatCurrency(1234.56)).toBe('R 1,234.56');
  });
});
```

**Component Tests:**
```typescript
describe('Button', () => {
  it('should render correctly', () => {
    const { getByText } = render(
      <Button title="Test" onPress={() => {}} />
    );
    expect(getByText('Test')).toBeTruthy();
  });

  it('should show loading state', () => {
    const { getByTestId } = render(
      <Button title="Test" onPress={() => {}} loading />
    );
    expect(getByTestId('loading-spinner')).toBeTruthy();
  });
});
```

---

## 📊 Test Results Template

### Test Session Report

**Date:** ___________
**Tester:** ___________
**Platform:** iOS / Android
**Device:** ___________
**Build:** ___________

**Results:**

| Test Case | Status | Notes |
|-----------|--------|-------|
| Login     | ✅ Pass | |
| Biometric | ✅ Pass | |
| Dashboard | ❌ Fail | Stats not loading |
| ...       | ...    | ... |

**Issues Found:**
1. _______________________________
2. _______________________________
3. _______________________________

**Overall Status:** Pass / Fail
**Ready for Release:** Yes / No

---

## 🔍 Regression Testing

**After Each Update:**

✅ Authentication flow works
✅ Dashboard loads correctly
✅ All navigation works
✅ API calls succeed
✅ State persists
✅ No console errors
✅ No crashes
✅ Performance acceptable

---

## 📞 Support

**Found a bug?**
1. Check if it's reproducible
2. Note device, OS version, and steps
3. Take screenshots
4. Report to: dev@comply360.com

---

**Version:** 1.0.0
**Last Updated:** December 28, 2025
**Test Coverage:** Comprehensive
