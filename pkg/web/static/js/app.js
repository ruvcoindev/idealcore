// idealcore — базовые утилиты для фронтенда

const API = {
  base: '/api',
  async request(endpoint, options = {}) {
    const url = this.base + endpoint;
    const config = {
      headers: { 'Content-Type': 'application/json', ...options.headers },
      ...options
    };
    const res = await fetch(url, config);
    if (!res.ok) {
      const err = await res.json().catch(() => ({}));
      throw new Error(err.error || `HTTP ${res.status}`);
    }
    return res.json();
  },
  post(endpoint, data) { return this.request(endpoint, { method: 'POST', body: JSON.stringify(data) }); },
  get(endpoint) { return this.request(endpoint, { method: 'GET' }); }
};

const FormUtils = {
  serialize(form) {
    const data = {};
    new FormData(form).forEach((value, key) => {
      if (data[key]) {
        if (!Array.isArray(data[key])) data[key] = [data[key]];
        data[key].push(value);
      } else { data[key] = value; }
    });
    return data;
  },
  getChecked(name) {
    return Array.from(document.querySelectorAll(`input[name="${name}"]:checked`)).map(el => el.value);
  },
  setCharCount(textarea, counter) {
    const count = textarea.value.length;
    counter.textContent = count;
    if (count > 500) counter.style.color = 'var(--warning)';
    else if (count > 1000) counter.style.color = 'var(--error)';
    else counter.style.color = 'var(--text-muted)';
  }
};

const UI = {
  show(el) { el.classList.remove('hidden'); },
  hide(el) { el.classList.add('hidden'); },
  setLoading(btn, loading) {
    if (loading) {
      btn.dataset.original = btn.innerHTML;
      btn.classList.add('loading');
      btn.disabled = true;
    } else {
      btn.innerHTML = btn.dataset.original || btn.innerHTML;
      btn.classList.remove('loading');
      btn.disabled = false;
    }
  },
  async copyText(text) {
    if (navigator.clipboard?.writeText) return navigator.clipboard.writeText(text);
    const ta = document.createElement('textarea');
    ta.value = text; ta.style.position = 'fixed'; ta.style.opacity = '0';
    document.body.appendChild(ta); ta.select(); document.execCommand('copy');
    document.body.removeChild(ta);
    return Promise.resolve();
  },
  toast(message, type = 'info') {
    const colors = { info: '#3b82f6', success: '#22c55e', error: '#ef4444' };
    const toast = document.createElement('div');
    toast.textContent = message;
    toast.style.cssText = `position: fixed; bottom: 2rem; right: 2rem; padding: 1rem 1.5rem; background: ${colors[type]}; color: white; border-radius: 8px; z-index: 1000; box-shadow: 0 4px 12px rgba(0,0,0,0.3); animation: fadeIn 0.3s ease-out;`;
    document.body.appendChild(toast);
    setTimeout(() => { toast.style.opacity = '0'; toast.style.transition = 'opacity 0.3s'; setTimeout(() => toast.remove(), 300); }, 3000);
  }
};

if (typeof window !== 'undefined') { window.API = API; window.FormUtils = FormUtils; window.UI = UI; }
