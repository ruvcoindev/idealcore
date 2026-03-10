// idealcore — логика страницы дневника

document.addEventListener('DOMContentLoaded', () => {
  const form = document.getElementById('diary-form');
  const content = document.getElementById('diary-content');
  const loader = document.getElementById('diary-loader');
  const userId = document.getElementById('user-id')?.value || 'anonymous';
  const analyzeBtn = document.getElementById('analyze-btn');
  const analysisResult = document.getElementById('analysis-result');
  const analysisContent = document.getElementById('analysis-content');

  // Загружаем вопросы
  async function loadQuestions() {
    try {
      const res = await fetch('/diary/questions');
      const data = await res.json();
      renderDiary(data.sections, data.questions);
      loader.classList.add('hidden');
      form.classList.remove('hidden');
      loadExistingAnswers();
    } catch (err) {
      loader.textContent = 'Ошибка загрузки: ' + err.message;
      console.error('Failed to load questions:', err);
    }
  }

  // Рендерим форму
  function renderDiary(sections, questions) {
    content.innerHTML = '';
    
    sections.forEach(section => {
      const sectionEl = document.createElement('section');
      sectionEl.className = 'diary-section';
      sectionEl.dataset.section = section.id;
      
      sectionEl.innerHTML = `
        <h2>${section.title}</h2>
        <p class="section-desc">${section.description}</p>
        ${questions.filter(q => q.section === section.id)
          .sort((a,b) => a.order - b.order)
          .map(q => `
            <div class="question-block" data-question="${q.id}">
              <label for="q-${q.id}">${q.text}</label>
              <textarea 
                id="q-${q.id}" 
                name="${q.id}" 
                class="question-input"
                placeholder="${q.placeholder}"
                minlength="${q.minLength}"
                maxlength="${q.maxLength}"
              ></textarea>
              <div class="question-meta">
                <span class="char-count">0</span> символов
                <button type="button" class="btn-save-small" data-qid="${q.id}">Сохранить</button>
              </div>
            </div>
          `).join('')}
      `;
      content.appendChild(sectionEl);
    });

    // Счётчики символов
    document.querySelectorAll('.question-input').forEach(el => {
      const counter = el.closest('.question-block').querySelector('.char-count');
      const updateCount = () => {
        counter.textContent = el.value.length;
        if (el.value.length > 500) counter.style.color = 'var(--warning)';
        else if (el.value.length > 1000) counter.style.color = 'var(--error)';
        else counter.style.color = 'var(--text-muted)';
      };
      el.addEventListener('input', updateCount);
      updateCount();
    });

    // Кнопки "Сохранить" для каждого вопроса
    document.querySelectorAll('.btn-save-small').forEach(btn => {
      btn.addEventListener('click', async (e) => {
        const qid = e.target.dataset.qid;
        const answer = document.getElementById(`q-${qid}`)?.value || '';
        const section = e.target.closest('.diary-section')?.dataset.section;
        
        if (!answer.trim()) {
          UI.toast('Пустой ответ не сохраняется', 'info');
          return;
        }

        UI.setLoading(e.target, true);
        try {
          await API.post('/diary/save', {
            user_id: userId, section, answer, tags: [qid]
          });
          UI.toast('✓ Сохранено', 'success');
        } catch (err) {
          UI.toast('Ошибка: ' + err.message, 'error');
        } finally {
          UI.setLoading(e.target, false);
        }
      });
    });
  }

  // Загружаем существующие ответы
  async function loadExistingAnswers() {
    try {
      const res = await fetch(`/api/diary/entries?user=${userId}`);
      const data = await res.json();
      if (data.entries) {
        Object.entries(data.entries).forEach(([qid, answer]) => {
          const el = document.getElementById(`q-${qid}`);
          if (el) el.value = answer;
          // Обновляем счётчик
          const counter = el?.closest('.question-block')?.querySelector('.char-count');
          if (counter) counter.textContent = answer.length;
        });
      }
    } catch (err) {
      console.warn('Could not load existing answers:', err);
    }
  }

  // Сохранение всей формы
  form?.addEventListener('submit', async (e) => {
    e.preventDefault();
    const entries = {};
    
    document.querySelectorAll('.question-input').forEach(el => {
      if (el.value.trim()) {
        const section = el.closest('.diary-section')?.dataset.section;
        entries[section] = el.value;
      }
    });

    if (Object.keys(entries).length === 0) {
      UI.toast('Нет ответов для сохранения', 'info');
      return;
    }

    const submitBtn = form.querySelector('button[type="submit"]');
    UI.setLoading(submitBtn, true);

    try {
      for (const [section, answer] of Object.entries(entries)) {
        await API.post('/diary/save', {
          user_id: userId, section, answer, tags: []
        });
      }
      UI.toast('✓ Все записи сохранены', 'success');
    } catch (err) {
      UI.toast('Ошибка: ' + err.message, 'error');
    } finally {
      UI.setLoading(submitBtn, false);
    }
  });

  // Анализ записей
  analyzeBtn?.addEventListener('click', async () => {
    const allText = Array.from(document.querySelectorAll('.question-input'))
      .map(el => el.value.trim()).filter(Boolean).join('\n\n');
    
    if (!allText) {
      UI.toast('Заполни дневник для анализа', 'info');
      return;
    }

    UI.setLoading(analyzeBtn, true);
    UI.show(analysisResult);
    analysisContent.innerHTML = '<p>Анализирую...</p>';

    try {
      // Заглушка анализа (пока нет бэкенда)
      analysisContent.innerHTML = `
        <p>✅ Записей проанализировано: ${allText.split('\n').filter(l=>l.trim()).length}</p>
        <p>💡 Подсказка: перечитай ответы через неделю — паттерны станут виднее.</p>
      `;
      UI.toast('✓ Анализ завершён', 'success');
    } catch (err) {
      analysisContent.innerHTML = `<p class="error">Ошибка: ${err.message}</p>`;
      UI.toast('Ошибка анализа', 'error');
    } finally {
      UI.setLoading(analyzeBtn, false);
    }
  });

  // Старт
  loadQuestions();
});
