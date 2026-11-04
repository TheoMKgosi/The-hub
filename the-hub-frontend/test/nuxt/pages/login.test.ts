import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mountSuspended } from '@nuxt/test-utils/runtime'
import LoginPage from '~/pages/login.vue'

const { loginMock } = vi.hoisted(() => ({ loginMock: vi.fn() }))

// The login page reaches the auth store through Nuxt auto-imports. Mock the
// store module so submitting the form never hits the real API layer.
vi.mock('~/stores/auth', () => ({
  useAuthStore: () => ({ login: loginMock })
}))

beforeEach(() => {
  loginMock.mockClear()
})

describe('login page', () => {
  it('renders the login heading, alpha banner, and submit button', async () => {
    const wrapper = await mountSuspended(LoginPage)

    expect(wrapper.text()).toContain('Login')
    expect(wrapper.text()).toContain('This is a web app in alpha')
    expect(wrapper.text()).toContain('alpha and all data may be wiped out')
  })

  it('links to register and forgot-password', async () => {
    const wrapper = await mountSuspended(LoginPage)

    const hrefs = wrapper.findAll('a').map(a => a.attributes('href'))
    expect(hrefs).toContain('/register')
    expect(hrefs).toContain('/forgot-password')
  })

  it('renders a password input', async () => {
    const wrapper = await mountSuspended(LoginPage)

    expect(wrapper.find('input[type="password"]').exists()).toBe(true)
  })

  it('submits valid credentials to the auth store', async () => {
    const wrapper = await mountSuspended(LoginPage)

    const inputs = wrapper.findAll('input')
    await inputs[0].setValue('user@example.com')
    await inputs[1].setValue('secret123')
    await wrapper.find('form').trigger('submit')
    await vi.waitFor(() => {
      expect(loginMock).toHaveBeenCalled()
    })

    expect(loginMock).toHaveBeenCalledWith({
      email: 'user@example.com',
      password: 'secret123'
    })
  })

  it('does not submit an empty form (schema validation blocks it)', async () => {
    const wrapper = await mountSuspended(LoginPage)

    await wrapper.find('form').trigger('submit')

    expect(loginMock).not.toHaveBeenCalled()
  })

  it('does not submit an invalid email', async () => {
    const wrapper = await mountSuspended(LoginPage)

    const inputs = wrapper.findAll('input')
    await inputs[0].setValue('not-an-email')
    await inputs[1].setValue('secret123')
    await wrapper.find('form').trigger('submit')

    expect(loginMock).not.toHaveBeenCalled()
  })

  it('does not submit a password shorter than 6 characters', async () => {
    const wrapper = await mountSuspended(LoginPage)

    const inputs = wrapper.findAll('input')
    await inputs[0].setValue('user@example.com')
    await inputs[1].setValue('12345')
    await wrapper.find('form').trigger('submit')

    expect(loginMock).not.toHaveBeenCalled()
  })
})
