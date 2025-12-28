import { Box, Container, Typography } from '@mui/material'
import React from 'react'

const AdminTemplatesPage: React.FC = () => {
    return (
        <Container>
            <Box sx={{ py: 4 }}>
                <Typography variant="h4" gutterBottom>
                    看板テンプレ管理
                </Typography>
                {/* TODO: Implement template management interface */}
            </Box>
        </Container>
    )
}

export default AdminTemplatesPage
