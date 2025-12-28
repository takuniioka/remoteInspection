import { expect, test } from '@playwright/test';

test.describe('Admin Features', () => {
    test.beforeEach(async ({ page, context }) => {
        // Login as admin
        await context.addInitScript(() => {
            localStorage.setItem('authToken', 'mock-admin-token');
        });
    });

    test('should access admin templates page', async ({ page }) => {
        await page.goto('/admin/templates');
        await expect(page.locator('text=Annotation Templates')).toBeVisible();
    });

    test('should create new template', async ({ page }) => {
        await page.goto('/admin/templates');
        await page.click('button:has-text("Create Template")');
        await page.fill('input[placeholder="Template Name"]', 'New Template');
        await page.click('button:has-text("Save")');
        // Template should be added to list
        await expect(page.locator('text=New Template')).toBeVisible();
    });

    test('should edit template', async ({ page }) => {
        await page.goto('/admin/templates');
        await page.click('button[title="Edit"]');
        await page.fill('input[placeholder="Template Name"]', 'Updated Template');
        await page.click('button:has-text("Save")');
        // Template should be updated
        await expect(page.locator('text=Updated Template')).toBeVisible();
    });

    test('should delete template', async ({ page }) => {
        await page.goto('/admin/templates');
        await page.click('button[title="Delete"]');
        // Should show confirmation
        await page.click('button:has-text("Confirm")');
        // Template should be removed
    });

    test('should view template usage statistics', async ({ page }) => {
        await page.goto('/admin/templates');
        await page.click('[data-test="template-item"]');
        // Should show statistics
        await expect(page.locator('text=Usage|Statistics')).toBeVisible();
    });
});

test.describe('User Management', () => {
    test.beforeEach(async ({ page, context }) => {
        await context.addInitScript(() => {
            localStorage.setItem('authToken', 'mock-admin-token');
        });
    });

    test('should view user list', async ({ page }) => {
        await page.goto('/admin/users');
        await expect(page.locator('text=Users')).toBeVisible();
    });

    test('should add new user', async ({ page }) => {
        await page.goto('/admin/users');
        await page.click('button:has-text("Add User")');
        await page.fill('input[placeholder="Email"]', 'newuser@example.com');
        await page.selectOption('select[name="role"]', 'editor');
        await page.click('button:has-text("Save")');
        // User should be added
        await expect(page.locator('text=newuser@example.com')).toBeVisible();
    });

    test('should change user role', async ({ page }) => {
        await page.goto('/admin/users');
        const userRow = page.locator('tr', { hasText: 'testuser' });
        await userRow.locator('select').selectOption('admin');
        await page.click('button:has-text("Save")');
        // Role should be updated
        await expect(userRow.locator('text=admin')).toBeVisible();
    });

    test('should deactivate user', async ({ page }) => {
        await page.goto('/admin/users');
        const userRow = page.locator('tr', { hasText: 'testuser' });
        await userRow.locator('button:has-text("Deactivate")').click();
        // User should be marked as inactive
        await expect(userRow.locator('text=Inactive')).toBeVisible();
    });
});

test.describe('Audit Logs', () => {
    test.beforeEach(async ({ page, context }) => {
        await context.addInitScript(() => {
            localStorage.setItem('authToken', 'mock-admin-token');
        });
    });

    test('should view audit logs', async ({ page }) => {
        await page.goto('/admin/audit-logs');
        await expect(page.locator('text=Audit Logs')).toBeVisible();
    });

    test('should filter audit logs by action', async ({ page }) => {
        await page.goto('/admin/audit-logs');
        await page.selectOption('select[name="action"]', 'CREATE_INSPECTION');
        // List should be filtered
        await expect(page.locator('text=CREATE_INSPECTION')).toBeVisible();
    });

    test('should filter audit logs by date', async ({ page }) => {
        await page.goto('/admin/audit-logs');
        await page.fill('input[type="date"]', '2024-01-01');
        // List should be filtered
    });

    test('should export audit logs', async ({ page }) => {
        await page.goto('/admin/audit-logs');
        const downloadPromise = page.waitForEvent('download');
        await page.click('button:has-text("Export")');
        const download = await downloadPromise;
        expect(download.suggestedFilename()).toContain('audit');
    });
});

test.describe('Report Generation', () => {
    test.beforeEach(async ({ page, context }) => {
        await context.addInitScript(() => {
            localStorage.setItem('authToken', 'mock-token');
        });
    });

    test('should generate HTML report', async ({ page }) => {
        await page.goto('/inspections/test-id');
        await page.click('button:has-text("Generate Report")');
        await page.selectOption('select[name="format"]', 'html');
        await page.click('button:has-text("Generate")');
        // Report should be generated
        await expect(page.locator('text=Report generated')).toBeVisible();
    });

    test('should download report', async ({ page }) => {
        await page.goto('/inspections/test-id');
        await page.click('button:has-text("Generate Report")');
        const downloadPromise = page.waitForEvent('download');
        await page.click('button:has-text("Download")');
        const download = await downloadPromise;
        expect(download.suggestedFilename()).toContain('report');
    });

    test('should preview report', async ({ page }) => {
        await page.goto('/inspections/test-id');
        await page.click('button:has-text("Generate Report")');
        await page.click('button:has-text("Preview")');
        // Report preview should open
        const popup = await page.context().pages()[1];
        await expect(popup.locator('text=Report')).toBeVisible();
    });
});

test.describe('Accessibility', () => {
    test('should have proper ARIA labels', async ({ page }) => {
        await page.goto('/login');
        // Check for ARIA labels
        await expect(page.locator('input[aria-label]')).toHaveCount(2);
    });

    test('should navigate with keyboard', async ({ page }) => {
        await page.goto('/inspections');
        await page.keyboard.press('Tab');
        // Focus should move to next element
        await expect(page.locator(':focus')).toBeVisible();
    });

    test('should support screen readers', async ({ page }) => {
        await page.goto('/inspections');
        // Page should be readable by screen readers
        const roleElements = await page.locator('[role]').count();
        expect(roleElements).toBeGreaterThan(0);
    });
});

test.describe('Mobile Responsiveness', () => {
    test.use({ viewport: { width: 375, height: 667 } });

    test('should display correctly on mobile', async ({ page }) => {
        await page.goto('/inspections');
        // Content should be visible and properly formatted
        await expect(page.locator('text=Inspections')).toBeVisible();
    });

    test('should handle touch interactions', async ({ page }) => {
        await page.goto('/inspections');
        // Touch events should work
        await page.locator('button').first().tap();
        await expect(page.locator('text=Inspections')).toBeVisible();
    });
});
