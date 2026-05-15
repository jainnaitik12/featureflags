function setStatus(message, isError) {
  const el = document.getElementById('status');
  el.textContent = message;
  el.className = isError ? 'error' : 'muted';
}

function renderFlags(flags) {
  const tbody = document.getElementById('flagsTable');
  tbody.replaceChildren();

  for (const flag of flags) {
    const tr = document.createElement('tr');

    const nameTd = document.createElement('td');
    nameTd.textContent = flag.name;
    tr.appendChild(nameTd);

    const enTd = document.createElement('td');
    enTd.textContent = flag.enabled ? 'true' : 'false';
    tr.appendChild(enTd);

    const actionsTd = document.createElement('td');
    const toggleBtn = document.createElement('button');
    toggleBtn.type = 'button';
    toggleBtn.className = flag.enabled ? 'btn-toggle-on' : 'btn-toggle-off';
    toggleBtn.textContent = flag.enabled ? 'Disable' : 'Enable';
    toggleBtn.addEventListener('click', function () {
      toggleFlag(flag.name);
    });

    const delBtn = document.createElement('button');
    delBtn.type = 'button';
    delBtn.className = 'btn-delete';
    delBtn.textContent = 'Delete';
    delBtn.addEventListener('click', function () {
      deleteFlag(flag.name);
    });

    actionsTd.appendChild(toggleBtn);
    actionsTd.appendChild(document.createTextNode(' '));
    actionsTd.appendChild(delBtn);
    tr.appendChild(actionsTd);

    tbody.appendChild(tr);
  }
}

function renderAudit(events) {
  const tbody = document.getElementById('auditTable');
  tbody.replaceChildren();

  for (const event of events) {
    const tr = document.createElement('tr');
    const ts = event.changed_at ? new Date(event.changed_at).toLocaleString() : '';

    function td(text) {
      const cell = document.createElement('td');
      cell.textContent = text;
      return cell;
    }

    tr.appendChild(td(ts));
    tr.appendChild(td(event.flag_name));
    tr.appendChild(td(event.action));
    tr.appendChild(td(String(event.old_value)));
    tr.appendChild(td(String(event.new_value)));
    tbody.appendChild(tr);
  }
}

async function createFlag() {
  const nameInput = document.getElementById('flagName');
  const enabledInput = document.getElementById('flagEnabled');
  const name = nameInput.value.trim();
  if (!name) return setStatus('Flag name is required', true);

  const res = await createFlagApi(name, enabledInput.checked);
  if (!res.ok) return setStatus('Create failed', true);
  nameInput.value = '';
  enabledInput.checked = false;
  setStatus('Flag created');
  await refresh();
}

async function toggleFlag(name) {
  const res = await toggleFlagApi(name);
  if (!res.ok) return setStatus('Toggle failed', true);
  setStatus('Toggled ' + name);
  await refresh();
}

async function deleteFlag(name) {
  const res = await deleteFlagApi(name);
  if (!res.ok) return setStatus('Delete failed', true);
  setStatus('Deleted ' + name);
  await refresh();
}

async function refresh() {
  try {
    const [flags, audit] = await Promise.all([fetchFlags(), fetchAudit()]);
    renderFlags(flags);
    renderAudit(audit);
  } catch (err) {
    setStatus(err.message, true);
  }
}

document.addEventListener('DOMContentLoaded', function () {
  document.getElementById('btnCreate').addEventListener('click', createFlag);
  refresh();
  setInterval(refresh, 5000);
});
