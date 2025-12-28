import {
    WSMessage
} from '@/types/api'

export type WSMessageHandler = (message: WSMessage) => void
export type WSConnectionHandler = () => void
export type WSErrorHandler = (error: Event) => void
export type WSCloseHandler = () => void

class WebSocketClient {
    private ws: WebSocket | null = null
    private url: string
    private token: string | null = null
    private reconnectAttempts = 0
    private maxReconnectAttempts = 5
    private reconnectDelay = 3000

    private messageHandlers: Map<string, WSMessageHandler[]> = new Map()
    private onConnect: WSConnectionHandler | null = null
    private onError: WSErrorHandler | null = null
    private onClose: WSCloseHandler | null = null

    constructor(url: string) {
        this.url = url
    }

    connect(token: string, inspectionId: string): Promise<void> {
        return new Promise((resolve, reject) => {
            this.token = token

            const wsURL = `${this.url}?inspectionId=${inspectionId}&token=${token}`

            this.ws = new WebSocket(wsURL)

            this.ws.onopen = () => {
                console.log('WebSocket connected')
                this.reconnectAttempts = 0
                if (this.onConnect) {
                    this.onConnect()
                }
                resolve()
            }

            this.ws.onmessage = (event) => {
                try {
                    const message: WSMessage = JSON.parse(event.data)
                    this.handleMessage(message)
                } catch (error) {
                    console.error('Failed to parse WebSocket message:', error)
                }
            }

            this.ws.onerror = (error) => {
                console.error('WebSocket error:', error)
                if (this.onError) {
                    this.onError(error)
                }
                reject(error)
            }

            this.ws.onclose = () => {
                console.log('WebSocket closed')
                if (this.onClose) {
                    this.onClose()
                }
                this.attemptReconnect(token, inspectionId)
            }
        })
    }

    private attemptReconnect(token: string, inspectionId: string) {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++
            setTimeout(() => {
                console.log(
                    `Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})`
                )
                this.connect(token, inspectionId).catch(() => {
                    // Silently fail and let the next attempt happen
                })
            }, this.reconnectDelay)
        }
    }

    disconnect() {
        if (this.ws) {
            this.ws.close()
            this.ws = null
        }
    }

    send(message: WSMessage) {
        if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
            console.error('WebSocket is not connected')
            return
        }
        this.ws.send(JSON.stringify(message))
    }

    on(messageType: string, handler: WSMessageHandler) {
        if (!this.messageHandlers.has(messageType)) {
            this.messageHandlers.set(messageType, [])
        }
        this.messageHandlers.get(messageType)!.push(handler)
    }

    off(messageType: string, handler: WSMessageHandler) {
        const handlers = this.messageHandlers.get(messageType)
        if (handlers) {
            const index = handlers.indexOf(handler)
            if (index !== -1) {
                handlers.splice(index, 1)
            }
        }
    }

    onConnected(handler: WSConnectionHandler) {
        this.onConnect = handler
    }

    onConnectionError(handler: WSErrorHandler) {
        this.onError = handler
    }

    onConnectionClosed(handler: WSCloseHandler) {
        this.onClose = handler
    }

    private handleMessage(message: WSMessage) {
        // Handle message type routing
        const handlers = this.messageHandlers.get(message.type) || []
        handlers.forEach((handler) => {
            try {
                handler(message)
            } catch (error) {
                console.error(`Error in message handler for ${message.type}:`, error)
            }
        })

        // Also call generic message handlers if they exist
        const genericHandlers = this.messageHandlers.get('*') || []
        genericHandlers.forEach((handler) => {
            try {
                handler(message)
            } catch (error) {
                console.error('Error in generic message handler:', error)
            }
        })
    }

    isConnected(): boolean {
        return this.ws !== null && this.ws.readyState === WebSocket.OPEN
    }
}

const wsClient = new WebSocketClient(
    import.meta.env.VITE_WS_URL || 'ws://localhost:8081'
)

export default wsClient
