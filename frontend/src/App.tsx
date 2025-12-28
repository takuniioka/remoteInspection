import { Box, CircularProgress } from '@mui/material'
import React, { useEffect, useState } from 'react'
import {
    Navigate,
    Route,
    BrowserRouter as Router,
    Routes,
} from 'react-router-dom'
import apiClient from './services/apiClient'
import { User } from './types/api'

// Pages
import AdminTemplatesPage from './pages/AdminTemplatesPage'
import AnnotatePhotoPage from './pages/AnnotatePhotoPage'
import FieldCameraPage from './pages/FieldCameraPage'
import GuestViewPage from './pages/GuestViewPage'
import InspectionListPage from './pages/InspectionListPage'
import LoginPage from './pages/LoginPage'
import RemoteInspectionPage from './pages/RemoteInspectionPage'

/**
 * Top-level application component. Responsible for reading any stored
 * authentication (or guest session) tokens, initializing the `apiClient`
 * with a token, and rendering the appropriate routes based on user
 * authentication state.
 */
const App: React.FC = () => {
    const [user, setUser] = useState<User | null>(null)
    const [loading, setLoading] = useState(true)
    const [guestMode, setGuestMode] = useState(false)

    useEffect(() => {
        // On mount, check localStorage for an auth token and attempt to
        // load the user profile. If the token is invalid, remove it.
        const token = localStorage.getItem('authToken')
        if (token) {
            apiClient.setToken(token)
            apiClient
                .getMe()
                .then(setUser)
                .catch(() => {
                    localStorage.removeItem('authToken')
                    setUser(null)
                })
                .finally(() => setLoading(false))
        } else {
            // If no auth token, check for a guest session token and set
            // a minimal guest user context.
            const guestToken = localStorage.getItem('guestSessionToken')
            if (guestToken) {
                apiClient.setToken(guestToken)
                setGuestMode(true)
                // Set a guest user
                setUser({
                    userId: 'guest',
                    username: 'guest',
                    email: '',
                    role: 'guestViewer',
                    createdAt: new Date().toISOString(),
                })
            }
            setLoading(false)
        }
    }, [])

    if (loading) {
        return (
            <Box
                display="flex"
                justifyContent="center"
                alignItems="center"
                minHeight="100vh"
            >
                <CircularProgress />
            </Box>
        )
    }

    return (
        <Router>
            <Routes>
                {/* Public routes */}
                <Route path="/login" element={<LoginPage />} />
                <Route path="/login/callback" element={<LoginPage />} />
                <Route path="/guest/:viewerToken" element={<GuestViewPage />} />

                {/* Protected routes */}
                {user ? (
                    <>
                        <Route path="/" element={<InspectionListPage />} />
                        <Route
                            path="/inspections/:inspectionId/remote"
                            element={<RemoteInspectionPage />}
                        />
                        <Route
                            path="/inspections/:inspectionId/field"
                            element={<FieldCameraPage />}
                        />
                        <Route
                            path="/inspections/:inspectionId/view"
                            element={<RemoteInspectionPage />}
                        />
                        <Route
                            path="/photos/:photoId/annotate"
                            element={<AnnotatePhotoPage />}
                        />
                        {user.role === 'admin' && (
                            <Route path="/admin/templates" element={<AdminTemplatesPage />} />
                        )}
                        <Route path="*" element={<Navigate to="/" replace />} />
                    </>
                ) : (
                    <Route path="*" element={<Navigate to="/login" replace />} />
                )}
            </Routes>
        </Router>
    )
}

export default App
