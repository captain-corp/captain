import { test, expect } from '@playwright/test';
import { faker } from '@faker-js/faker';
import { login, logout } from '../helpers/auth';
import { fillDateTimeField, setVisibility } from '../helpers/posts';

const randomTitle = faker.lorem.words(3);

test.describe('Admin Panel E2E Tests', () => {
    test('complete admin workflow', async ({ page }) => {
        await logout(page);

        // Login
        await test.step('Login to admin panel', async () => {
            await login(page);

            // Verify successful login
            await expect(page).toHaveURL(/.*\/admin/);
        });

        // Posts Management
        await test.step('Create and manage posts', async () => {
            const postTitle = 'Test Post';
            const titleFieldSelector = 'input[name="title"]';
            const slugFieldSelector = 'input[name="slug"]';
            const expectedSlug = 'test-post';

            // Create first post with specific data
            await page.click('text=Create Your First Post');

            // Set title
            await page.locator(titleFieldSelector).pressSequentially(postTitle, {delay: 100});

            // Validate the slug
            const slugValue = await page.locator(slugFieldSelector).inputValue();
            expect(slugValue).toBe(expectedSlug);

            // Set excerpt
            await page.fill('textarea[name="excerpt"]', 'Test excerpt');

            // Set content
            await page.fill('textarea[name="content"]', 'Test content');

            // Set tags
            await page.fill('input[name="tags"]', 'test');
            await page.press('input[name="tags"]', 'Enter');
            await page.fill('input[name="tags"]', 'e2e');
            await page.press('input[name="tags"]', 'Enter');

            // Set publish type to "scheduled" and set publish date to "1985-10-26T10:00"
            await page.selectOption('select[name="publish"]', 'scheduled');
                                           //2006-01-02T15:04:05Z07:00
            await fillDateTimeField(page, '1985-10-26T10:00:00Z');

            await expect(page.locator('#selected-tags')).toContainText('test');
            await expect(page.locator('#selected-tags')).toContainText('e2e');

            // Set visibility
            await setVisibility(page, true);

            await page.click('button:has-text("Save")');

            // Verify flash message
            await expect(page.locator('.flash-message.flash-success')).toBeVisible();
            await expect(page.locator('.flash-message.flash-success')).toContainText('Post created successfully');

            // Verify post creation
            await expect(page.locator(`text=${postTitle}`)).toBeVisible();

            // Edit post
            await page.click(`a:has-text("Edit")`);

            // Edit title
            await page.locator(titleFieldSelector).clear();
            await page.locator(titleFieldSelector).pressSequentially("New title", {delay: 100});

            // Verify slug remained same
            const slugValueEdited = await page.locator(slugFieldSelector).inputValue();
            expect(slugValueEdited).toBe(expectedSlug);

            await page.click('button:has-text("Save")');

            // Verify post update
            await expect(page.locator(`text="New title"`)).toBeVisible();

            /**
             * Second Post
             */

            // Create second post with random data
            await page.click('text=Create New Post');
            await page.locator(titleFieldSelector).pressSequentially(randomTitle);
            await page.fill('textarea[name="excerpt"]', 'Test excerpt');
            await page.fill('textarea[name="content"]', 'Test content');

            await page.fill('input[id="tags"]', 'another');
            await page.press('input[id="tags"]', 'Enter');
            await page.fill('input[id="tags"]', 'test');
            await page.press('input[id="tags"]', 'Enter');
            
            await setVisibility(page, true);

            await page.click('button:has-text("Save")');

            // Verify both posts are present
            await expect(page.locator(`text="New title"`)).toBeVisible();
            await expect(page.locator(`text=${randomTitle}`)).toBeVisible();
        });

        // Pages Management
        await test.step('Create and manage pages', async () => {
            const pageTitle = 'Test Page';
            const pageSlug = 'test-page';

            await page.click('text=Pages');
            await page.click('text=Create Your First Page');
            await page.locator('input[name="title"]').pressSequentially(pageTitle);
            const slugValue = await page.locator('input[name="slug"]').inputValue();
            expect(slugValue).toBe(pageSlug);
            await page.fill('textarea[name="content"]', 'Test page content');
            await page.click('button:has-text("Save")');

            // Verify page creation
            await expect(page.locator(`text=${pageTitle}`)).toBeVisible();

            // Edit page
            await page.click('a:has-text("Edit")');
            await page.locator('input[name="title"]').clear();
            await page.locator('input[name="title"]').pressSequentially("Edited Page title");
            const immutableSlugValue = await page.locator('input[name="slug"]').inputValue();
            expect(immutableSlugValue).toBe(pageSlug);
            await page.fill('textarea[name="content"]', 'Test content');
            await page.click('button:has-text("Save")');
        });

        // Menu Management
        await test.step('Create and verify menu items', async () => {
            await page.click('text=Menu Items');

            // Create custom URL menu item
            await page.click('text=Create your first menu item!');
            await page.fill('input[name="label"]', 'Custom Link');
            await page.fill('input[name="url"]', 'https://example.com');
            await page.click('button:has-text("Create Menu Item")');

            // Create page link menu item
            await page.click('text=Create New Menu Item');
            await page.fill('input[name="label"]', 'Test Page');
            await page.selectOption('select[name="page_id"]', { label: 'Edited Page title' });
            await page.click('button:has-text("Create Menu Item")');

            // Verify menu items
            await page.goto('/');
            await expect(page.locator('a[href="https://example.com"]')).toBeVisible();
            await expect(page.locator('a[href="/pages/test-page"]')).toBeVisible();
        });

        // Tags Management
        await test.step('Verify tags functionality', async () => {
            await page.goto('/');

            // Verify tags are listed
            await expect(page.locator('.post-tag:has-text("test")').first()).toBeVisible();
            await expect(page.locator('.post-tag:has-text("another")').first()).toBeVisible();
            await expect(page.locator('.post-tag:has-text("e2e")').first()).toBeVisible();

            // Click tag and verify posts
            await page.click('.post-tag:has-text("test")');
            await expect(page.locator('text=New title')).toBeVisible();
            await expect(page.locator(`text=${randomTitle}`)).toBeVisible();

            // Check public tag page
            await page.goto('/tags/e2e');
            await expect(page.locator('text=New title')).toBeVisible();
            await page.goto('/tags/another');
            await expect(page.locator('text=New title')).not.toBeVisible();
        });

        // Settings Management
        await test.step('Modify and verify settings', async () => {
            await page.goto('/admin/settings');

            await page.fill('input[name="title"]', 'Updated Title');
            await page.fill('input[name="subtitle"]', 'Updated Subtitle');
            await page.fill('input[name="posts-per-page"]', '1');

            await page.click('button:has-text("Save")');
            await expect(page.locator('button:has-text("Saved")')).toBeVisible();

            // Verify post time updated
            await page.goto('/admin/posts');
            const thirdPost = page.locator('tbody tr').nth(2);
            const publishedAt = await thirdPost.locator('td').nth(2).textContent();
            expect(publishedAt).toContain('October 26, 1985 at 11:00 AM');

        });

        // Verify changes on public site
        await test.step('Verify changes on public site', async () => {
            // Test navigation back to public site
            await page.goto('/');
            await expect(page).toHaveURL('/');

            await expect(page.locator('.pagination')).toBeVisible();
            await expect(page.locator('a[href="?page=2"]')).toBeVisible();
            await expect(page.locator('text=Updated Title')).toBeVisible();
            await expect(page.locator('text=Updated Subtitle')).toBeVisible();
            const title = await page.title();
            await expect(title).toContain('Updated Title');
        });
    });
});
