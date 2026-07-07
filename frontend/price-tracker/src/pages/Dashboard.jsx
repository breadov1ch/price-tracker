import { useEffect, useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';

import api, { clearAccessToken } from '../api';
import './Dashboard.css';

/**
 * Визначає статус ціни для кольорового акценту.
 */
function getPriceStatus(product) {
  if (product.status === 'up' || product.status === 'down' || product.status === 'same') {
    return product.status;
  }
  return 'same';
}

function formatPrice(product) {
  const value = Number(product.product_last_price ?? 0);
  const currency = product.currency || 'UAH';

  try {
    return new Intl.NumberFormat('uk-UA', { style: 'currency', currency }).format(value);
  } catch {
    return `${value.toFixed(2)} ${currency}`;
  }
}

function formatDelta(product) {
  return null;
}

export default function Dashboard() {
  const navigate = useNavigate();

  const [products, setProducts] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [deletingId, setDeletingId] = useState(null);

  useEffect(() => {
    let isMounted = true;

    api
      .get('/api/products')
      .then(({ data }) => {
        if (isMounted) setProducts(data ?? []);
      })
      .catch(() => {
        if (isMounted) setError('Не вдалося завантажити товари. Спробуйте оновити сторінку.');
      })
      .finally(() => {
        if (isMounted) setIsLoading(false);
      });

    return () => {
      isMounted = false;
    };
  }, []);

  const handleDelete = async (id) => {
    setDeletingId(id);
    const previous = products;
    
    setProducts((prev) => prev.filter((p) => p.product_id !== id));

    try {
      await api.delete(`/api/products/${id}`);
    } catch {
      setProducts(previous);
      setError('Не вдалося видалити товар. Спробуйте ще раз.');
    } finally {
      setDeletingId(null);
    }
  };

  const handleLogout = () => {
    clearAccessToken();
    navigate('/login', { replace: true });
  };

  return (
    <div className="dashboard-page-dark">
      {/* Верхняя панель управления */}
      <header className="dash-header-opium">
        <div className="brand-zone">
          <span className="brand-eyebrow">MONITORING SYSTEM</span>
          <h1 className="brand-title-uppercase">TRACKED_ITEMS</h1>
        </div>

        <div className="action-zone">
          <button className="btn-opium-ghost" onClick={handleLogout}>
            LOGOUT_
          </button>
          <Link to="/create" className="btn-opium-solid">
            + ADD_NEW_ITEM
          </Link>
        </div>
      </header>

      {/* Системные сообщения */}
      {isLoading && <p className="system-status-msg">LOADING_SYSTEM_DATA...</p>}
      {!isLoading && error && <p className="system-status-msg error-neon">{error}</p>}

      {/* Пустое состояние */}
      {!isLoading && !error && products.length === 0 && (
        <div className="opium-empty-state">
          <p className="empty-text">NO_ITEMS_TRACKED_YET</p>
          <Link to="/create" className="btn-opium-solid">
            INITIALIZE_FIRST_TRACKER
          </Link>
        </div>
      )}

      {/* Сетка карточек */}
      {!isLoading && products.length > 0 && (
        <main className="opium-grid">
          {products.map((product) => {
            const status = getPriceStatus(product);
            const delta = formatDelta(product);

            return (
              <div className="opium-card" key={product.product_id}>
                {/* Медиа-блок (Изображение) */}
                <div className="card-media-vault">
                  {product.image_url ? (
                    <img 
                      src={product.image_url} 
                      alt={product.product_name} 
                      className="card-img-filtered" 
                    />
                  ) : (
                    <div className="card-img-void">
                      <span>IMAGE_NOT_FOUND</span>
                    </div>
                  )}
                </div>

                {/* Информационный блок */}
                <div className="card-meta-body">
                  <div className="meta-head-row">
                    <h3 className="item-title-bold">{product.product_name || 'UNNAMED_PRODUCT'}</h3>
                    <a
                      className="item-redirect-badge"
                      href={product.product_url}
                      target="_blank"
                      rel="noreferrer"
                    >
                      SRC_↗
                    </a>
                  </div>

                  {/* Дата-панель метрик */}
                  <div className="metrics-display-subgrid">
                    <div className="metric-box">
                      <span className="metric-label">LAST_PRICE</span>
                      <span className={`metric-value price-color-${status}`}>
                        {formatPrice(product)}
                      </span>
                    </div>
                    
                  
                  </div>

                  {/* Подвал карточки */}
                  <div className="card-meta-footer">
                    <span className="timestamp-data">
                      REFRESHED: {product.updated_at 
                        ? new Date(product.updated_at).toLocaleString('uk-UA', {
                            day: '2-digit',
                            month: '2-digit',
                            year: 'numeric',
                            hour: '2-digit',
                            minute: '2-digit'
                          }) 
                        : 'NOW'}
                    </span>
                    
                    <button
                      className="btn-destroy-item"
                      onClick={() => handleDelete(product.product_id)}
                      disabled={deletingId === product.product_id}
                    >
                      {deletingId === product.product_id ? 'WAIT...' : 'DESTROY_'}
                    </button>
                  </div>
                </div>
              </div>
            );
          })}
        </main>
      )}
    </div>
  );
}