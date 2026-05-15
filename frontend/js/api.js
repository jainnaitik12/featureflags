async function fetchFlags() {
  const res = await fetch('/api/flags');
  if (!res.ok) throw new Error('Failed to fetch flags');
  return res.json();
}

async function fetchAudit() {
  const res = await fetch('/audit');
  if (!res.ok) throw new Error('Failed to fetch audit log');
  return res.json();
}

async function createFlagApi(name, enabled) {
  const res = await fetch('/api/flags', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name, enabled })
  });
  return res;
}

async function toggleFlagApi(name) {
  return fetch(`/api/flags/${encodeURIComponent(name)}/toggle`, { method: 'PATCH' });
}

async function deleteFlagApi(name) {
  return fetch(`/api/flags/${encodeURIComponent(name)}`, { method: 'DELETE' });
}
