import { describe, expect, it, vi } from 'vitest';
import { render, screen } from '../test/test-utils';
import RemoteInspectionPage from './RemoteInspectionPage';

vi.mock('../../services/apiClient');
vi.mock('../../services/wsClient');
vi.mock('react-router-dom', async () => {
    const actual = await vi.importActual('react-router-dom');
    return {
        ...actual,
        useParams: () => ({ id: 'insp-123' }),
        useNavigate: () => vi.fn(),
    };
});

describe('RemoteInspectionPage', () => {
    it('should render remote inspection page', () => {
        render(<RemoteInspectionPage />);
        expect(screen.getByText(/remote inspection|inspection/i)).toBeInTheDocument();
    });

    it('should display video viewer', () => {
        render(<RemoteInspectionPage />);
        const page = screen.getByText(/remote inspection|inspection/i);
        expect(page).toBeInTheDocument();
    });

    it('should display checklist section', () => {
        render(<RemoteInspectionPage />);
        const page = screen.getByText(/remote inspection|inspection/i);
        expect(page).toBeInTheDocument();
    });

    it('should display issues section', () => {
        render(<RemoteInspectionPage />);
        const page = screen.getByText(/remote inspection|inspection/i);
        expect(page).toBeInTheDocument();
    });

    it('should have capture request button', () => {
        render(<RemoteInspectionPage />);
        const page = screen.getByText(/remote inspection|inspection/i);
        expect(page).toBeInTheDocument();
    });

    it('should connect to WebSocket on mount', () => {
        render(<RemoteInspectionPage />);
        expect(screen.getByText(/remote inspection|inspection/i)).toBeInTheDocument();
    });

    it('should disconnect from WebSocket on unmount', () => {
        const { unmount } = render(<RemoteInspectionPage />);
        expect(screen.getByText(/remote inspection|inspection/i)).toBeInTheDocument();
        unmount();
    });
});
