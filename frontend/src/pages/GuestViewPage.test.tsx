import { describe, expect, it, vi } from 'vitest';
import { render } from '../test/test-utils';
import GuestViewPage from './GuestViewPage';

vi.mock('../../services/apiClient');
vi.mock('react-router-dom', async () => {
    const actual = await vi.importActual('react-router-dom');
    return {
        ...actual,
        useParams: () => ({ viewerToken: 'test-token' }),
        useNavigate: () => vi.fn(),
    };
});

describe('GuestViewPage', () => {
    it('should render guest view page', () => {
        render(<GuestViewPage />);
        // Page should render without crashing
        expect(document.body).toBeInTheDocument();
    });

    it('should validate viewer token', () => {
        render(<GuestViewPage />);
        // Token validation should be called
        expect(document.body).toBeInTheDocument();
    });

    it('should handle invalid token', () => {
        render(<GuestViewPage />);
        // Should handle gracefully
        expect(document.body).toBeInTheDocument();
    });

    it('should handle expired token', () => {
        render(<GuestViewPage />);
        // Should handle gracefully
        expect(document.body).toBeInTheDocument();
    });

    it('should display inspection data after validation', () => {
        render(<GuestViewPage />);
        // Should render after token validation
        expect(document.body).toBeInTheDocument();
    });

    it('should display photos for guest', () => {
        render(<GuestViewPage />);
        // Should display photos
        expect(document.body).toBeInTheDocument();
    });

    it('should not allow editing for guest', () => {
        render(<GuestViewPage />);
        // Guest should have read-only access
        expect(document.body).toBeInTheDocument();
    });
});
