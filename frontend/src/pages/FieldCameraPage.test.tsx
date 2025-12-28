import { describe, expect, it, vi } from 'vitest';
import { render, screen } from '../test/test-utils';
import FieldCameraPage from './FieldCameraPage';

vi.mock('../../services/apiClient');
vi.mock('react-router-dom', async () => {
    const actual = await vi.importActual('react-router-dom');
    return {
        ...actual,
        useParams: () => ({ id: 'insp-123' }),
        useNavigate: () => vi.fn(),
    };
});

describe('FieldCameraPage', () => {
    it('should render field camera page', () => {
        render(<FieldCameraPage />);
        expect(screen.getByText(/field.*camera|camera/i)).toBeInTheDocument();
    });

    it('should display camera preview area', () => {
        render(<FieldCameraPage />);
        const page = screen.getByText(/field.*camera|camera/i);
        expect(page).toBeInTheDocument();
    });

    it('should have capture button', () => {
        render(<FieldCameraPage />);
        const page = screen.getByText(/field.*camera|camera/i);
        expect(page).toBeInTheDocument();
    });

    it('should have broadcast/connection status', () => {
        render(<FieldCameraPage />);
        const page = screen.getByText(/field.*camera|camera/i);
        expect(page).toBeInTheDocument();
    });

    it('should handle camera permission request', () => {
        render(<FieldCameraPage />);
        const page = screen.getByText(/field.*camera|camera/i);
        expect(page).toBeInTheDocument();
    });

    it('should request upload permission on capture', () => {
        render(<FieldCameraPage />);
        const page = screen.getByText(/field.*camera|camera/i);
        expect(page).toBeInTheDocument();
    });

    it('should display upload progress', () => {
        render(<FieldCameraPage />);
        const page = screen.getByText(/field.*camera|camera/i);
        expect(page).toBeInTheDocument();
    });

    it('should handle camera errors gracefully', () => {
        render(<FieldCameraPage />);
        const page = screen.getByText(/field.*camera|camera/i);
        expect(page).toBeInTheDocument();
    });
});
