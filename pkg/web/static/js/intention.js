// idealcore — логика страницы генерации намерений

document.addEventListener('DOMContentLoaded', () => {
  const form = document.getElementById('intention-form');
  const userId = document.getElementById('user-id')?.value || 'anonymous';
  const generateBtn = document.getElementById('generate-btn');
  const copyBtn = document.getElementById('copy-btn');
  const result = document.getElementById('result');
  const intentionText = document.getElementById('intention-text');
  const ragInfo = document.getElementById('rag-info');

  const themeSelect = document.getElementById('theme');
  if (themeSelect && themeSelect.options.length <= 1) {
    API.get('/intention/themes').then(data => {
      data.themes?.forEach(theme => {
        const opt = document.createElement('option');
        opt.value = theme; opt.textContent = theme;
        themeSelect.appendChild(opt);
      });
    }).catch(() => console.warn('Не удалось загрузить темы'));
  }

  form?.addEventListener('submit', async (e) => {
    e.preventDefault();
    const theme = document.getElementById('theme')?.value;
    if (!theme) { UI.toast('Выбери тему намерения', 'info'); return; }
    UI.setLoading(generateBtn, true);
    copyBtn.disabled = true;
    UI.hide(result);
    try {
      const payload = {
        user_id: userId, theme,
        trauma_types: FormUtils.getChecked('trauma'),
        attachment_style: document.getElementById('attachment')?.value || '',
        defense_mechanisms: FormUtils.getChecked('defense'),
        chakra_focus: FormUtils.getChecked('chakra'),
        raw_prompt: document.getElementById('raw-prompt')?.value || ''
      };
      const data = await API.post('/intention/generate', payload);
      intentionText.textContent = data.intention;
      if (data.rag_used) UI.show(ragInfo); else UI.hide(ragInfo);
      UI.show(result); copyBtn.disabled = false;
      UI.toast('✓ Намерение сгенерировано', 'success');
      result.scrollIntoView({ behavior: 'smooth', block: 'start' });
    } catch (err) { UI.toast(`Ошибка: ${err.message}`, 'error'); }
    finally { UI.setLoading(generateBtn, false); }
  });

  copyBtn?.addEventListener('click', async () => {
    const text = intentionText.textContent;
    if (!text) return;
    try {
      await UI.copyText(text);
      copyBtn.textContent = '✓ Скопировано';
      setTimeout(() => { copyBtn.textContent = 'Копировать'; }, 2000);
      UI.toast('Намерение скопировано в буфер', 'success');
    } catch { UI.toast('Не удалось скопировать', 'error'); }
  });

  document.querySelectorAll('input[name="trauma"], input[name="chakra"]').forEach(cb => {
    cb.addEventListener('change', (e) => { if (e.target.checked) UI.toast(`✓ Добавлено: ${e.target.value}`, 'info'); });
  });
});
