// idealcore — Service Worker для PWA
const CACHE_NAME = 'idealcore-v1';
const STATIC_ASSETS = [
  '/', '/static/css/style.css', '/static/js/app.js',
  '/static/js/diary.js', '/static/js/intention.js', '/static/js/insights.js',
  '/static/icons/favicon.ico'
];

self.addEventListener('install', (event) => {
  event.waitUntil(caches.open(CACHE_NAME).then(cache => cache.addAll(STATIC_ASSETS)).then(() => self.skipWaiting()));
});

self.addEventListener('activate', (event) => {
  event.waitUntil(caches.keys().then(names => Promise.all(names.filter(name => name !== CACHE_NAME).map(name => caches.delete(name)))).then(() => self.clients.claim()));
});

self.addEventListener('fetch', (event) => {
  const { request } = event;
  const url = new URL(request.url);
  if (url.pathname.startsWith('/api')) { event.respondWith(fetch(request)); return; }
  event.respondWith(caches.match(request).then(cached => {
    const networked = fetch(request).then(response => {
      if (response.ok && response.type === 'basic') {
        const clone = response.clone();
        caches.open(CACHE_NAME).then(cache => cache.put(request, clone));
      }
      return response;
    }).catch(() => cached);
    return cached || networked;
  }));
});

self.addEventListener('sync', (event) => {
  if (event.tag === 'sync-diary') { event.waitUntil(console.log('Background sync: diary entries')); }
});
