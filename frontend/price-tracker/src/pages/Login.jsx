import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';

import api, { setAccessToken } from '../api';
import './Auth.css';

export default function Login() {
  const navigate = useNavigate();

  const [form, setForm] = useState({ email: '', password: '' });
  const [error, setError] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleChange = (event) => {
    const { name, value } = event.target;
    setForm((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    setError('');
    setIsSubmitting(true);

    try {
      const { data } = await api.post('/api/auth/login', form);
      setAccessToken(data.access_token);
      navigate('/dashboard', { replace: true });
    } catch (err) {
      setError(
        err.response?.data?.message || 'Не вдалося увійти. Перевірте email і пароль.'
      );
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="auth-page">
      <div className="auth-card">
        <p className="auth-eyebrow">Price Tracer</p>
        <h1 className="auth-title">Вхід</h1>

        <form onSubmit={handleSubmit} noValidate>
          <div className="field">
            <label htmlFor="email">Email</label>
            <input
              id="email"
              name="email"
              type="email"
              autoComplete="email"
              placeholder="you@example.com"
              value={form.email}
              onChange={handleChange}
              required
            />
          </div>

          <div className="field">
            <label htmlFor="password">Пароль</label>
            <input
              id="password"
              name="password"
              type="password"
              autoComplete="current-password"
              placeholder="••••••••"
              value={form.password}
              onChange={handleChange}
              required
            />
          </div>

          {error && <p className="error-text">{error}</p>}

          <button type="submit" className="btn btn-block" disabled={isSubmitting}>
            {isSubmitting ? 'Входимо…' : 'Увійти'}
          </button>
        </form>

        <p className="auth-footer">
          Немає акаунта? <Link to="/register">Зареєструватися</Link>
        </p>
      </div>
    </div>
  );
}
