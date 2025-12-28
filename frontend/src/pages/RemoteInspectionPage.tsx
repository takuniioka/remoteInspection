import { Box, Container, Typography } from '@mui/material'
import React from 'react'

const RemoteInspectionPage: React.FC = () => {
    return (
        <Container>
            <Box sx={{ py: 4 }}>
                <Typography variant="h4" gutterBottom>
                    リモート検収画面
                </Typography>
                {/* TODO: Implement remote inspection interface */}
            </Box>
        </Container>
    )
}

export default RemoteInspectionPage
