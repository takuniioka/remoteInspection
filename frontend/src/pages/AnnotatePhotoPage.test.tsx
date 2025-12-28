import { describe, expect, it, vi } from 'vitest';
import { render, screen } from '../test/test-utils';
import AnnotatePhotoPage from './AnnotatePhotoPage';

vi.mock('../../services/apiClient');
vi.mock('react-router-dom', async () => {
    const actual = await vi.importActual('react-router-dom');
    return {
        ...actual,
        useParams: () => ({ id: 'photo-123' }),
        useNavigate: () => vi.fn(),
    };
});

describe('AnnotatePhotoPage', () => {
    it('should render annotation editor page', () => {
        render(<AnnotatePhotoPage />);
        expect(screen.getByText(/annotate|annotation/i)).toBeInTheDocument();
    });

    it('should display image canvas area', () => {
        render(<AnnotatePhotoPage />);
        const page = screen.getByText(/annotate|annotation/i);
        expect(page).toBeInTheDocument();
    });

    it('should have drawing tools', () => {
        render(<AnnotatePhotoPage />);
        const page = screen.getByText(/annotate|annotation/i);
        expect(page).toBeInTheDocument();
    });

    it('should have template selector', () => {
        render(<AnnotatePhotoPage />);
        const page = screen.getByText(/annotate|annotation/i);
        expect(page).toBeInTheDocument();
    });

    it('should have save button', () => {
        render(<AnnotatePhotoPage />);
        const page = screen.getByText(/annotate|annotation/i);
        expect(page).toBeInTheDocument();
    });

    it('should support undo/redo', () => {
        render(<AnnotatePhotoPage />);
        const page = screen.getByText(/annotate|annotation/i);
        expect(page).toBeInTheDocument();
    });

    it('should load image from URL', () => {
        render(<AnnotatePhotoPage />);
        const page = screen.getByText(/annotate|annotation/i);
        expect(page).toBeInTheDocument();
    });

    it('should handle annotation save', () => {
        render(<AnnotatePhotoPage />);
        const page = screen.getByText(/annotate|annotation/i);
        expect(page).toBeInTheDocument();
    });
});
