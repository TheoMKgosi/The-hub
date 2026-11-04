import { describe, it, expect } from 'vitest'
import { mountSuspended } from '@nuxt/test-utils/runtime'
import PrivacyPolicyPage from '~/pages/privacy-policy.vue'
import TermsOfServicePage from '~/pages/terms-of-service.vue'

describe('privacy policy page', () => {
  it('renders the title and last-updated date', async () => {
    const wrapper = await mountSuspended(PrivacyPolicyPage)

    expect(wrapper.find('h1').text()).toBe('Privacy Policy')
    expect(wrapper.text()).toContain('Last updated: December 2024')
  })

  it('renders the main policy sections', async () => {
    const wrapper = await mountSuspended(PrivacyPolicyPage)

    const headings = wrapper.findAll('h2').map(h => h.text())
    expect(headings).toEqual([
      'Information We Collect',
      'How We Use Your Information',
      'Information Sharing',
      'Data Security',
      'Contact Us'
    ])
  })

  it('includes the contact email', async () => {
    const wrapper = await mountSuspended(PrivacyPolicyPage)

    expect(wrapper.text()).toContain('privacy@thehub.com')
  })
})

describe('terms of service page', () => {
  it('renders the title', async () => {
    const wrapper = await mountSuspended(TermsOfServicePage)

    expect(wrapper.find('h1').text()).toBe('Terms of Service')
  })

  it('renders the main sections', async () => {
    const wrapper = await mountSuspended(TermsOfServicePage)

    const headings = wrapper.findAll('h2').map(h => h.text())
    expect(headings).toContain('Acceptance of Terms')
    expect(headings).toContain('Use License')
    expect(headings).toContain('User Responsibilities')
    expect(headings).toContain('Prohibited Uses')
    expect(headings).toContain('Service Availability')
    expect(headings).toContain('Contact Information')
  })
})
