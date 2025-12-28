import { Box, Container, Typography } from '@mui/material'
import React from 'react'

const AnnotatePhotoPage: React.FC = () => {
    return (
        <Container>
            <Box sx={{ py: 4 }}>
                <Typography variant="h4" gutterBottom>
                    写真注釈編集
                </Typography>
                {/* TODO: Implement annotation canvas with Konva */}
            </Box>
        </Container>
    )
}

export default AnnotatePhotoPage
