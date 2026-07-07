import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';

import api from '../api';
import './Auth.css';

export default function Register() {
  const navigate = useNavigate();

  const [form, setForm] = useState({ name: '', email: '', password: '' });
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
      await api.post('/api/auth/register', form);
      // Реєстрація не логінить автоматично — ведемо на /login,
      // де юзер вже свідомо вводить пароль і отримує токени.
      navigate('/login', { replace: true });
    } catch (err) {
      setError(
        err.response?.data?.message || 'Не вдалося зареєструватися. Спробуйте ще раз.'
      );
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="auth-page">
      <div className="auth-card">
        <p className="auth-eyebrow">Price Tracker</p>
        <h1 className="auth-title">Реєстрація</h1>

        <form onSubmit={handleSubmit} noValidate>
          <div className="field">
            <label htmlFor="name">Ім'я</label>
            <input
              id="name"
              name="name"
              type="text"
              autoComplete="name"
              placeholder="Ваше ім'я"
              value={form.name}
              onChange={handleChange}
              required
            />
          </div>

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
              autoComplete="new-password"
              placeholder="••••••••"
              value={form.password}
              onChange={handleChange}
              required
              minLength={8}
            />
          </div>

          {error && <p className="error-text">{error}</p>}

          <button type="submit" className="btn btn-block" disabled={isSubmitting}>
            {isSubmitting ? 'Створюємо акаунт…' : 'Зареєструватися'}
          </button>
        </form>

        <p className="auth-footer">
          Вже є акаунт? <Link to="/login">Увійти</Link>
        </p>
      </div>
    </div>
  );
}
