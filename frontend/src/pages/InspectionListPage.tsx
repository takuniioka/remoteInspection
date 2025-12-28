import { Box, Button, Container, Typography } from '@mui/material'
import React from 'react'
import { useNavigate } from 'react-router-dom'

const InspectionListPage: React.FC = () => {
    const navigate = useNavigate()

    return (
        <Container>
            <Box sx={{ py: 4 }}>
                <Typography variant="h4" gutterBottom>
                    検収案件一覧
                </Typography>
                <Button
                    variant="contained"
                    color="primary"
                    onClick={() => navigate('/inspections/new')}
                    sx={{ mt: 2 }}
                >
                    新規案件作成
                </Button>
                {/* TODO: Display list of inspections */}
            </Box>
        </Container>
    )
}

export default InspectionListPage
