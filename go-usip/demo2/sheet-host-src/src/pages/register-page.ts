import { renderAuthShell, attachFormMessage } from '../components/auth-shell'
import { register } from '../services/auth-service'

export function renderRegisterPage() {
  renderAuthShell(
    'Create Account',
    'Join and start using your workspace in seconds.',
    `<form id="register-form" class="auth-form">
      <label for="register-nickname"><b>Nickname</b></label>
      <input id="register-nickname" type="text" name="nickname" required>
      <label for="register-username"><b>Username</b></label>
      <input id="register-username" type="text" name="username" required>
      <label for="register-password"><b>Password</b></label>
      <input id="register-password" type="password" name="password" required>
      <div id="form-message" class="form-message"></div>
      <button type="submit" class="auth-submit">Register</button>
    </form>
    <footer class="auth-footer">
      <span>Already have an account?</span>
      <a href="/login">Go to login</a>
    </footer>`,
  )

  const form = document.querySelector<HTMLFormElement>('#register-form')
  if (!form)
    return

  form.addEventListener('submit', async (event) => {
    event.preventDefault()
    const formData = new FormData(form)

    try {
      await register(
        String(formData.get('nickname') ?? ''),
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
