import { Box, Container, Typography } from '@mui/material'
import React from 'react'

const FieldCameraPage: React.FC = () => {
    return (
        <Container>
            <Box sx={{ py: 4 }}>
                <Typography variant="h4" gutterBottom>
                    現場配信・撮影画面
                </Typography>
                {/* TODO: Implement field camera interface */}
            </Box>
        </Container>
    )
}

export default FieldCameraPage
