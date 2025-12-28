import {
    Alert,
    Box,
    Button,
    Card,
    CardContent,
    Container,
    TextField,
    Typography,
} from '@mui/material'
import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'

const LoginPage: React.FC = () => {
    const navigate = useNavigate()
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [error, setError] = useState('')
    const [isLoading, setIsLoading] = useState(false)

    const handleLogin = async (e: React.FormEvent) => {
        e.preventDefault()
        setIsLoading(true)
        setError('')

        // TODO: Implement Cognito authentication
        // For MVP, use mock login
        try {
            // Mock token for development
            const mockToken = 'mock-jwt-token-' + Date.now()
            localStorage.setItem('authToken', mockToken)
            navigate('/')
        } catch (err) {
            setError('Login failed. Please try again.')
        } finally {
            setIsLoading(false)
        }
    }

    return (
        <Container maxWidth="sm">
            <Box
                display="flex"
                justifyContent="center"
                alignItems="center"
                minHeight="100vh"
            >
                <Card sx={{ width: '100%' }}>
                    <CardContent>
                        <Typography variant="h5" gutterBottom align="center">
                            Remote Inspection Capture Tool
                        </Typography>
                        <Box component="form" onSubmit={handleLogin} noValidate>
                            {error && (
                                <Alert severity="error" sx={{ mb: 2 }}>
                                    {error}
                                </Alert>
                            )}
                            <TextField
                                fullWidth
                                label="Email"
                                type="email"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                margin="normal"
                                required
                            />
                            <TextField
                                fullWidth
                                label="Password"
                                type="password"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                margin="normal"
                                required
                            />
                            <Button
                                fullWidth
                                variant="contained"
                                color="primary"
                                type="submit"
                                sx={{ mt: 3 }}
                                disabled={isLoading}
                            >
                                {isLoading ? 'ログイン中...' : 'ログイン'}
                            </Button>
                        </Box>
                    </CardContent>
                </Card>
            </Box>
        </Container>
    )
}

export default LoginPage
