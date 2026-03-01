# Phase 1 Analysis - Critical Issues Found

## Executive Summary
The migration from Wails to web server is **NOT COMPLETE**. While builds succeed, the core download flow is BROKEN. Real-time progress will NOT work in production.

## Critical Issues (Must Fix)

### 🔴 ISSUE #1: NO SSE PROGRESS DURING DOWNLOADS
**Severity**: BLOCKER
**Location**: `backend/progress.go:149-174`, `server/handlers.go:203-342`

**Problem**:
- ProgressWriter updates global state but does NOT broadcast to SSE broker
- Download handler broadcasts only on start/complete/error
- NO progress events during active download
- Frontend will show 0% → 100% jump, no real-time updates

**Evidence**:
```go
// progress.go:149-174 - ProgressWriter.Write()
SetDownloadProgress(mbDownloaded)  // ← Updates global only
UpdateItemProgress(pw.itemID, mbDownloaded, speedMBps)  // ← Local only
// NO: broker.Publish() call!
```

**Root Cause**: ProgressWriter has no reference to SSE broker. Backend downloaders (tidal.go:251, qobuz.go:321, amazon.go) use ProgressWriter without callback mechanism.

**Fix Required**:
1. Add `progressCallback func(itemID string, percent float64, speed float64)` parameter to downloader functions
2. Pass callback from handler that calls `s.sseBroker.BroadcastJSON()`
3. Modify ProgressWriter to accept and call callback on progress updates

---

### 🔴 ISSUE #2: NO EVENTSOURCE FOR DOWNLOAD PROGRESS IN FRONTEND
**Severity**: BLOCKER
**Location**: `frontend/src/App.tsx:180-234`, `frontend/src/hooks/useDownload.ts`

**Problem**:
- EventSource is created ONLY for FFmpeg install (line 184)
- NO EventSource connection for track downloads
- Frontend cannot receive download progress events
- useDownload hook polls `/api/download-progress` instead of listening to SSE

**Evidence**:
```typescript
// App.tsx:180 - EventSource ONLY in handleInstallFFmpeg()
eventSource = new EventSource('/api/events');
eventSource.addEventListener('ffmpeg:progress', ...);  // ← FFmpeg only!
// NO: 'download:progress' listener
```

**Fix Required**:
1. Create persistent EventSource connection on app mount
2. Add event listeners for download progress events
3. Update useDownload hook to use SSE instead of polling
4. Ensure cleanup on unmount

---

### 🟡 ISSUE #3: DOWNLOAD PATH NOT ENFORCED IN FOLDER TEMPLATE
**Severity**: HIGH
**Location**: `server/handlers.go:238-239`

**Problem**:
```go
// Apply folder template if needed from settings
// The actual folder creation logic should be in backend
```
This is a STUB comment! Folder template is NOT applied. Files go directly to downloadPath root, ignoring user's folder template setting.

**Fix Required**:
1. Load settings in handler
2. Call `backend.BuildFolderPath(settings.FolderTemplate, trackMetadata)`
3. Join with base downloadPath: `filepath.Join(s.downloadPath, subfolder)`
4. Create directory before download

---

### 🟡 ISSUE #4: SETTINGS LOAD/SAVE NOT VERIFIED
**Severity**: MEDIUM
**Location**: `server/handlers.go:667-719`

**Problem**:
- HandleLoadSettings reads from file, not bbolt
- HandleSaveSettings writes to file, not bbolt
- Settings may not persist across restarts
- Need to verify bbolt DB functions exist and work

**Fix Required**:
1. Check if `backend.LoadSettings()` and `backend.SaveSettings()` exist
2. Verify they use bbolt, not file system
3. Ensure defaults are merged on load
4. Test persistence across server restart

---

### 🟡 ISSUE #5: NO SPA FALLBACK ROUTING
**Severity**: MEDIUM
**Location**: `main.go:168`

**Problem**:
```go
e.GET("/*", echo.WrapHandler(http.FileServer(http.FS(frontendFS))))
```
This serves files directly. Refreshing on `/settings` returns 404 because `settings` file doesn't exist. Need SPA fallback to serve `index.html` for non-file paths.

**Fix Required**:
Implement proper SPA routing as per problem statement Fix I

---

### 🟡 ISSUE #6: DOWNLOAD HANDLER USES WRONG FILE PATH CHECK
**Severity**: MEDIUM
**Location**: `backend/tidal.go:485-488`

**Problem**:
```go
if fileInfo, err := os.Stat(outputFilename); err == nil && fileInfo.Size() > 0 {
    return "EXISTS:" + outputFilename, nil
}
```
Handler checks for file existence in backend, but doesn't communicate "already exists" status properly via SSE. Frontend will show error instead of "skipped".

**Fix Required**:
1. Return special status when file exists
2. Broadcast SSE event with status="exists"
3. Frontend should show "Already Downloaded" not error

---

### 🟢 ISSUE #7: ERROR HANDLING IN API.TS
**Severity**: LOW
**Location**: `frontend/src/lib/api.ts`

**Problem**: Some fetch calls don't handle non-JSON error responses. If backend returns plain text error, `.json()` will throw.

**Fix Required**: Implement `safeFetch` helper per problem statement Fix F

---

### 🟢 ISSUE #8: MISSING LOADING STATES IN UI
**Severity**: LOW
**Location**: Various components

**Problem**: Some async actions don't show loading indicators

**Fix Required**: Add spinners/loading states per problem statement Fix K

---

## Analysis Checklist Status

### 1.1 Download Flow
- ✅ React component: SearchBar/TrackInfo trigger download
- ✅ Hook: useDownload.ts
- ✅ API: api.ts downloadTrack()
- ✅ HTTP: POST /api/download
- ✅ Handler: server/handlers.go HandleDownloadTrack
- ✅ Backend: TidalDownloader.Download() exists
- ❌ **Progress callback**: NOT wired to SSE broker
- ✅ OutputDir: Correctly set to server path
- ❌ **Folder template**: NOT applied (stub comment)
- ✅ File write: ProgressWriter writes to disk
- ❌ **SSE events**: Only start/complete/error, NOT progress

### 1.2 SSE Flow
- ✅ Backend has SSEBroker in server/sse.go
- ✅ GET /api/events registered
- ✅ SSEHandler sets correct headers
- ✅ Handler flushes after write
- ❌ **Frontend EventSource**: ONLY for FFmpeg, NOT downloads
- ❌ **Progress callback**: Backend doesn't call broker.Publish()
- ❌ **Event shape**: No percent field in progress events
- ❌ **Frontend parsing**: No listener for download events

### 1.3 Settings Flow
- ✅ Frontend calls GET /api/settings
- ⚠️ **Backend function**: Uses file, not bbolt (needs verification)
- ⚠️ **Defaults**: Need to verify merge logic
- ⚠️ **Persistence**: Need to test across restart
- ✅ Download path: Shown as read-only in UI

### 1.4 Metadata Flow
- ✅ POST /api/metadata endpoint exists
- ✅ Handler calls backend.GetFilteredSpotifyData
- ✅ Response parsed as JSON
- ✅ Frontend displays track/album/playlist info
- ✅ Individual track selection works

### 1.5 UI States
- ✅ Idle state
- ⚠️ Loading metadata - needs verification
- ✅ Metadata loaded states
- ❌ **Downloading** - progress will jump 0→100%
- ✅ Download complete
- ✅ Download error
- ✅ Settings panel
- ⚠️ History panel - needs verification
- ✅ Theme switching

### 1.6 Runtime Errors
- ⚠️ Some fetch calls missing .catch()
- ⚠️ Missing null checks in components
- ❌ **EventSource not closed** on unmount (FFmpeg only, but wrong pattern)
- ✅ CORS: Same origin, should work
- ✅ Content-Type headers set
- ⚠️ JSON parse errors: Some calls vulnerable
- ⚠️ DOWNLOAD_PATH: Created on startup, needs verification
- ⚠️ ffmpeg: Error handling exists but needs testing

### 1.7 Build Pipeline
- ✅ `pnpm run build` succeeds (verified in previous session)
- ✅ `go build .` succeeds (verified in previous session)
- ✅ frontend/dist/ populated
- ✅ embed directive correct
- ✅ Binary runs
- ⚠️ Docker build: Needs verification

---

## Priority Fix Order

1. **BLOCKER** Issue #1: Wire ProgressWriter to SSE broker
2. **BLOCKER** Issue #2: Add EventSource for downloads in frontend
3. **HIGH** Issue #3: Apply folder template in download handler
4. **MEDIUM** Issue #4: Verify settings persistence
5. **MEDIUM** Issue #5: Add SPA fallback routing
6. **MEDIUM** Issue #6: Handle "file exists" status properly
7. **LOW** Issue #7: Add error handling in api.ts
8. **LOW** Issue #8: Add loading states

---

## Next Steps
Proceed to Phase 2 - Fix all issues in priority order.
