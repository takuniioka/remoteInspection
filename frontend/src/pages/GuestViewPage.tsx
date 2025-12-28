import apiClient from '@/services/apiClient'
import { Box, CircularProgress, Container, Typography } from '@mui/material'
import React, { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'

const GuestViewPage: React.FC = () => {
    const { viewerToken } = useParams<{ viewerToken: string }>()
    const navigate = useNavigate()
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState('')

    useEffect(() => {
        const validateAndCreateSession = async () => {
            if (!viewerToken) {
                setError('Invalid viewer token')
                return
            }

            try {
                // Exchange viewer access token for session token
                const session = await apiClient.createGuestSession(viewerToken)
                localStorage.setItem('guestSessionToken', session.guestSessionToken)
                apiClient.setToken(session.guestSessionToken)

                // Redirect to inspection view (need inspectionId)
                // TODO: Get inspectionId from ViewerLink
                navigate('/')
            } catch (err) {
                setError('Failed to create guest session. Link may be invalid or expired.')
            } finally {
                setLoading(false)
            }
        }

        validateAndCreateSession()
    }, [viewerToken, navigate])

    if (loading) {
        return (
            <Container>
                <Box
                    display="flex"
                    justifyContent="center"
                    alignItems="center"
                    minHeight="100vh"
                >
                    <CircularProgress />
                </Box>
            </Container>
        )
    }

    if (error) {
        return (
            <Container>
                <Box sx={{ py: 4 }}>
                    <Typography color="error">{error}</Typography>
                </Box>
            </Container>
        )
    }

    return (
        <Container>
            <Box sx={{ py: 4 }}>
                <Typography variant="h4" gutterBottom>
                    ゲスト閲覧
                </Typography>
            </Box>
        </Container>
    )
}

export default GuestViewPage
