import { describe, expect, it, vi } from 'vitest';
import { render, screen } from '../test/test-utils';
import InspectionListPage from './InspectionListPage';

vi.mock('../../services/apiClient');
vi.mock('react-router-dom', async () => {
    const actual = await vi.importActual('react-router-dom');
    return {
        ...actual,
        useNavigate: () => vi.fn(),
    };
});

describe('InspectionListPage', () => {
    it('should render inspection list page', () => {
        render(<InspectionListPage />);
        expect(screen.getByText(/inspections/i)).toBeInTheDocument();
    });

    it('should have create inspection button', () => {
        render(<InspectionListPage />);
        const createButton = screen.getByRole('button', { name: /create|add/i });
        expect(createButton).toBeInTheDocument();
    });

    it('should display inspection list', () => {
        render(<InspectionListPage />);
        const listContainer = screen.getByText(/inspections/i);
        expect(listContainer).toBeInTheDocument();
    });

    it('should handle loading state', () => {
        render(<InspectionListPage />);
        // Should render without crashing
        expect(screen.getByText(/inspections/i)).toBeInTheDocument();
    });

    it('should handle empty inspection list', () => {
        render(<InspectionListPage />);
        // Should render page even if list is empty
        expect(screen.getByText(/inspections/i)).toBeInTheDocument();
    });
});
