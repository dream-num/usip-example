import { renderAuthShell, attachFormMessage } from '../components/auth-shell'
import { login } from '../services/auth-service'

export function renderLoginPage() {
  renderAuthShell(
    'Welcome Back',
    'Sign in to continue to your workspace.',
    `<form id="login-form" class="auth-form">
      <label for="login-username"><b>Username</b></label>
      <input id="login-username" type="text" name="username" required>
      <label for="login-password"><b>Password</b></label>
      <input id="login-password" type="password" name="password" required>
      <div id="form-message" class="form-message"></div>
      <button type="submit" class="auth-submit">Login</button>
    </form>
    <footer class="auth-footer">
      <span>No account yet?</span>
      <a href="/register">Go to register</a>
    </footer>`,
  )

  const form = document.querySelector<HTMLFormElement>('#login-form')
  if (!form)
    return

  form.addEventListener('submit', async (event) => {
    event.preventDefault()
    const formData = new FormData(form)

    try {
      await login(
        String(formData.get('username') ?? ''),
        String(formData.get('password') ?? ''),
      )
      location.href = '/files'
    }
    catch (error) {
      attachFormMessage((error as Error).message)
    }
  })
}
