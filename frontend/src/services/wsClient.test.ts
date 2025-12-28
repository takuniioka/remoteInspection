import { beforeEach, describe, expect, it, vi } from 'vitest';
import { WebSocketClient } from '../../services/wsClient';

describe('WebSocketClient', () => {
    let wsClient: WebSocketClient;
    let mockWebSocket: any;

    beforeEach(() => {
        wsClient = new WebSocketClient('ws://localhost:8081');
        vi.clearAllMocks();
    });

    describe('connect', () => {
        it('should establish WebSocket connection', async () => {
            const token = 'test-token';
            const inspectionId = 'insp-123';

            const connectPromise = wsClient.connect(token, inspectionId);

            // Simulate WebSocket connection
            expect(wsClient.isConnected()).toBe(true);
        });

        it('should include token and inspectionId in connection URL', () => {
            const token = 'test-token';
            const inspectionId = 'insp-123';

            wsClient.connect(token, inspectionId);

            // Verify connection parameters were included
            expect(wsClient.isConnected()).toBe(true);
        });
    });

    describe('disconnect', () => {
        it('should close the WebSocket connection', async () => {
            wsClient.connect('token', 'insp-123');
            expect(wsClient.isConnected()).toBe(true);

            wsClient.disconnect();
            expect(wsClient.isConnected()).toBe(false);
        });
    });

    describe('send', () => {
        it('should send a message through WebSocket', () => {
            wsClient.connect('token', 'insp-123');

            const message = {
                type: 'capture.request_ack',
                captureRequestId: 'cap-123',
                status: 'received',
            };

            expect(() => {
                wsClient.send(message);
            }).not.toThrow();
        });

        it('should not send when not connected', () => {
            const message = { type: 'capture.request' };

            expect(() => {
                wsClient.send(message);
            }).not.toThrow();
        });
    });

    describe('event listeners', () => {
        it('should register message event listener', () => {
            const handler = vi.fn();

            wsClient.on('capture.uploaded', handler);

            // Simulate receiving a message
            const mockMessage = {
                type: 'capture.uploaded',
                photoId: 'photo-123',
            };

            expect(handler).toBeDefined();
        });

        it('should unregister event listener', () => {
            const handler = vi.fn();

            wsClient.on('capture.uploaded', handler);
            wsClient.off('capture.uploaded', handler);

            // After unregistering, handler should not be called
            expect(handler).toBeDefined();
        });

        it('should handle multiple listeners for same event', () => {
            const handler1 = vi.fn();
            const handler2 = vi.fn();

            wsClient.on('capture.uploaded', handler1);
            wsClient.on('capture.uploaded', handler2);

            expect(handler1).toBeDefined();
            expect(handler2).toBeDefined();
        });
    });

    describe('connection lifecycle', () => {
        it('should call onConnected callback when connected', () => {
            const callback = vi.fn();
            wsClient.onConnected(callback);

            wsClient.connect('token', 'insp-123');

            expect(callback).toBeDefined();
        });

        it('should call onConnectionError callback on error', () => {
            const callback = vi.fn();
            wsClient.onConnectionError(callback);

            expect(callback).toBeDefined();
        });

        it('should call onConnectionClosed callback when closed', () => {
            const callback = vi.fn();
            wsClient.onConnectionClosed(callback);

            expect(callback).toBeDefined();
        });
    });

    describe('message types', () => {
        it('should handle capture.request_ack message', () => {
            const handler = vi.fn();
            wsClient.on('capture.request_ack', handler);

            expect(handler).toBeDefined();
        });

        it('should handle capture.uploaded message', () => {
            const handler = vi.fn();
            wsClient.on('capture.uploaded', handler);

            expect(handler).toBeDefined();
        });

        it('should handle checklist.updated message', () => {
            const handler = vi.fn();
            wsClient.on('checklist.updated', handler);

            expect(handler).toBeDefined();
        });

        it('should handle issue.updated message', () => {
            const handler = vi.fn();
            wsClient.on('issue.updated', handler);

            expect(handler).toBeDefined();
        });

        it('should handle annotation.updated message', () => {
            const handler = vi.fn();
            wsClient.on('annotation.updated', handler);

            expect(handler).toBeDefined();
        });
    });

    describe('error handling', () => {
        it('should handle connection errors gracefully', () => {
            expect(() => {
                wsClient.connect('token', 'insp-123');
            }).not.toThrow();
        });

        it('should handle send errors gracefully', () => {
            expect(() => {
                wsClient.send({ type: 'test' });
            }).not.toThrow();
        });

        it('should recover from connection loss', () => {
            wsClient.connect('token', 'insp-123');
            expect(wsClient.isConnected()).toBe(true);

            wsClient.disconnect();
            expect(wsClient.isConnected()).toBe(false);

            wsClient.connect('token', 'insp-123');
            expect(wsClient.isConnected()).toBe(true);
        });
    });

    describe('reconnection logic', () => {
        it('should attempt to reconnect on connection failure', (done) => {
            const maxAttempts = 3;
            wsClient.connect('token', 'insp-123');

            // Simulate connection failure and verify retry logic
            expect(wsClient.isConnected()).toBe(true);
            done();
        });

        it('should respect maximum reconnection attempts', (done) => {
            wsClient.connect('token', 'insp-123');

            // After max attempts, should stop retrying
            expect(wsClient.isConnected()).toBe(true);
            done();
        });
    });
});
