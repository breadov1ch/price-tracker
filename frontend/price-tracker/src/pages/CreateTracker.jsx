import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';

import api from '../api';
import './CreateTracker.css';

export default function CreateTracker() {
  const navigate = useNavigate();

  const [url, setUrl] = useState('');
  const [name, setName] = useState(''); // Стейт для назви
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false); // Для блокування кнопки

  const handleSubmit = async (event) => {
    event.preventDefault();
    setError('');
    setSuccess(false);
    setIsSubmitting(true);

    try {
      // ВІДПРАВЛЯЄМО РІВНО ТІ КЛЮЧІ, ЯКІ ОЧІКУЄ ТВІЙ GIN: "name" та "url"
      await api.post('/api/products', { 
        url: url,
        name: name 
      });
      
      setSuccess(true);
      setUrl('');
      setName('');
      
      setTimeout(() => navigate('/dashboard'), 2000);
    } catch (err) {
      setError(
        err.response?.data?.message || 'Не вдалося додати товар. Перевірте посилання.'
      );
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="create-page">
      <nav className="create-nav">
        <Link to="/dashboard">← До панелі</Link>
      </nav>

      <div className="create-card">
        <p className="create-eyebrow">Новий трекер</p>
        <h1 className="create-title">Додати товар</h1>
        <p className="create-hint">
          Вкажи назву та встав посилання на товар — ми занесемо його у систему
          та почнемо моніторинг.
        </p>

        <form onSubmit={handleSubmit} noValidate>
          {/* ПОЛЕ НАЗВИ */}
          <div className="field">
            <label htmlFor="name">Назва товару</label>
            <input
              id="name"
              name="name"
              type="text"
              placeholder="Наприклад: Ноутбук / Світшот Opium"
              value={name}
              onChange={(event) => setName(event.target.value)}
              required
            />
          </div>

          {/* ПОЛЕ URL */}
          <div className="field">
            <label htmlFor="url">URL товару</label>
            <input
              id="url"
              name="url"
              type="url"
              placeholder="https://example.com/product/123"
              value={url}
              onChange={(event) => setUrl(event.target.value)}
              required
            />
          </div>

          {error && <p className="error-text">{error}</p>}
          {success && <p className="success-text">Додано. Переходимо на панель…</p>}

          <button type="submit" className="btn btn-block" disabled={isSubmitting}>
            {isSubmitting ? 'Додаємо…' : 'Почати відстежувати'}
          </button>
        </form>
      </div>
    </div>
  );
}