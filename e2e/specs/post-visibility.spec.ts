import { test, expect } from '@playwright/test';
import { login, logout } from '../helpers/auth';
import { createPost } from '../helpers/posts';

test.describe('Post Visibility', () => {

  const now = new Date();
  const tomorrow = new Date(now);
  tomorrow.setDate(tomorrow.getDate() + 1);

  test('should handle draft and scheduled posts correctly', async ({ page }) => {
    await logout(page);

    // Create a draft post
    await login(page);
    await createPost(page, {
      title: 'Draft Post',
      content: 'This is a draft post',
      visible: false,
      publishedAt: now.toISOString(),
    });

    // Create a scheduled post
    await createPost(page, {
      title: 'Scheduled Post',
      content: 'This is a scheduled post',
      visible: true,
      publishedAt: tomorrow.toISOString(),
    });

    await page.goto('/admin/posts');

    // Verify both posts are present
    await expect(page.locator('td:has-text("Draft")')).toBeVisible();
    await expect(page.locator('td:has-text("Scheduled")')).toBeVisible();

    // Go to settings page
    await page.goto('/admin/settings');
    // Set page per page to 50
    await page.fill('input[name="posts-per-page"]', '50');
    await page.click('button:has-text("Save")');
    await expect(page.locator('button:has-text("Saved")')).toBeVisible();

    // Verify posts are visible and properly marked when logged in
    await page.goto('/');

    // Check draft and scheduled post
    const draftPostCount = await page.locator('article:has-text("Draft")').count();
    const scheduledPostCount = await page.locator('article:has-text("Scheduled")').count();
    const editLinkCount = await page.locator('.edit-link').count();
    expect(draftPostCount).toBeGreaterThan(0);
    expect(scheduledPostCount).toBeGreaterThan(0);
    expect(editLinkCount).toBeGreaterThan(0);

    // Logout and verify posts are not visible to anonymous users
    await page.goto('/logout');
    await page.goto('/');

    // Verify posts are not visible
    await expect(page.locator('article:has-text("Draft")')).not.toBeVisible();
    await expect(page.locator('article:has-text("Scheduled")')).not.toBeVisible();
  });
});
