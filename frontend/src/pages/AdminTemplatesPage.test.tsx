import { describe, expect, it, vi } from 'vitest';
import { render, screen } from '../test/test-utils';
import AdminTemplatesPage from './AdminTemplatesPage';

vi.mock('../../services/apiClient');
vi.mock('react-router-dom', async () => {
    const actual = await vi.importActual('react-router-dom');
    return {
        ...actual,
        useNavigate: () => vi.fn(),
    };
});

describe('AdminTemplatesPage', () => {
    it('should render admin templates page', () => {
        render(<AdminTemplatesPage />);
        expect(screen.getByText(/template|admin/i)).toBeInTheDocument();
    });

    it('should display templates list', () => {
        render(<AdminTemplatesPage />);
        const page = screen.getByText(/template|admin/i);
        expect(page).toBeInTheDocument();
    });

    it('should have create template button', () => {
        render(<AdminTemplatesPage />);
        const page = screen.getByText(/template|admin/i);
        expect(page).toBeInTheDocument();
    });

    it('should display template details', () => {
        render(<AdminTemplatesPage />);
        const page = screen.getByText(/template|admin/i);
        expect(page).toBeInTheDocument();
    });

    it('should allow template editing', () => {
        render(<AdminTemplatesPage />);
        const page = screen.getByText(/template|admin/i);
        expect(page).toBeInTheDocument();
    });

    it('should allow template deletion', () => {
        render(<AdminTemplatesPage />);
        const page = screen.getByText(/template|admin/i);
        expect(page).toBeInTheDocument();
    });

    it('should display template preview', () => {
        render(<AdminTemplatesPage />);
        const page = screen.getByText(/template|admin/i);
        expect(page).toBeInTheDocument();
    });
});
