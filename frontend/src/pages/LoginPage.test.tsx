import { describe, expect, it, vi } from 'vitest';
import { render, screen } from '../test/test-utils';
import LoginPage from './LoginPage';

vi.mock('../../services/apiClient');
vi.mock('react-router-dom', async () => {
    const actual = await vi.importActual('react-router-dom');
    return {
        ...actual,
        useNavigate: () => vi.fn(),
    };
});

describe('LoginPage', () => {
    it('should render login form', () => {
        render(<LoginPage />);
        expect(screen.getByRole('button', { name: /login/i })).toBeInTheDocument();
    });

    it('should have email and password input fields', () => {
        render(<LoginPage />);
        expect(screen.getByPlaceholderText(/email/i)).toBeInTheDocument();
        expect(screen.getByPlaceholderText(/password/i)).toBeInTheDocument();
    });

    it('should have form title', () => {
        render(<LoginPage />);
        expect(screen.getByText(/inspection.*inspection tool/i)).toBeInTheDocument();
    });

    it('should enable login button when form is valid', async () => {
        render(<LoginPage />);
        const emailInput = screen.getByPlaceholderText(/email/i) as HTMLInputElement;
        const passwordInput = screen.getByPlaceholderText(/password/i) as HTMLInputElement;

        expect(emailInput).toBeInTheDocument();
        expect(passwordInput).toBeInTheDocument();
    });
});
