const API_BASE = '/api/links';

function el(id){return document.getElementById(id)}
async function loadLinks(){
  const container = el('links-list');
  container.textContent = 'Loading…';
  try{
    const res = await fetch(API_BASE);
    if(!res.ok) throw new Error(`Server returned ${res.status}`);
    const data = await res.json();
    renderLinks(data || []);
  }catch(err){
    container.innerHTML = `<div class="msg">Failed to load links: ${err.message}</div>`;
  }
}

function renderLinks(links){
  const container = el('links-list');
  if(!links.length){
    container.innerHTML = '<div class="msg">No shortlinks yet.</div>';
    return;
  }

  const table = document.createElement('table');
  table.className = 'links-table';

  const thead = document.createElement('thead');
  const headerRow = document.createElement('tr');
  ['Short', 'Destination', 'Actions'].forEach((h, i) => {
    const th = document.createElement('th');
    th.textContent = h;
    if (h === 'Actions') th.className = 'small';
    headerRow.appendChild(th);
  });
  thead.appendChild(headerRow);
  table.appendChild(thead);

  const tbody = document.createElement('tbody');
  links.forEach(l => {
    const tr = document.createElement('tr');
    const origin = window.location.origin.replace(/:\d+$/,'');
      const shortUrl = `${origin}/${encodeURIComponent(l.shortName)}`;

    const tdShort = document.createElement('td');
    const aShort = document.createElement('a');
    aShort.setAttribute('href', shortUrl);
    aShort.setAttribute('target', '_blank');
    aShort.setAttribute('rel', 'noopener noreferrer');
    aShort.textContent = `/${l.shortName}`;
    tdShort.appendChild(aShort);

    const tdDest = document.createElement('td');
    // Render destination as a safe link only if it appears to be an http(s) URL
    try {
      const parsed = new URL(l.url);
      if (parsed.protocol === 'http:' || parsed.protocol === 'https:') {
        const aDest = document.createElement('a');
        aDest.setAttribute('href', l.url);
        aDest.setAttribute('target', '_blank');
        aDest.setAttribute('rel', 'noopener noreferrer');
        aDest.textContent = l.url;
        tdDest.appendChild(aDest);
      } else {
        tdDest.textContent = l.url;
      }
    } catch (err) {
      tdDest.textContent = l.url;
    }

    const tdActions = document.createElement('td');
    tdActions.className = 'small';
    const copy = document.createElement('button');
    copy.className = 'copy-btn';
    copy.textContent = 'Copy';
    copy.addEventListener('click',()=>{navigator.clipboard.writeText(shortUrl); copy.textContent='Copied'; setTimeout(()=>copy.textContent='Copy',1200)});
    tdActions.appendChild(copy);

    tr.appendChild(tdShort);
    tr.appendChild(tdDest);
    tr.appendChild(tdActions);
    tbody.appendChild(tr);
  });

  table.appendChild(tbody);
  container.innerHTML = '';
  container.appendChild(table);
}

function showCreateMsg(text, isError=false){
  const elMsg = el('create-msg');
  elMsg.textContent = text;
  elMsg.style.color = isError ? '#b91c1c' : 'inherit';
}

async function createLink(e){
  e.preventDefault();
  const shortName = el('shortname').value.trim();
  const url = el('url').value.trim();
  if(!shortName || !url) return showCreateMsg('Both fields are required', true);
  const btn = e.submitter || e.target.querySelector('button[type=submit]');
  btn.disabled = true;
  showCreateMsg('Creating...');
  try{
    const res = await fetch(API_BASE, {
      method: 'POST',
      headers: {'Content-Type':'application/json'},
      body: JSON.stringify({ shortName, url })
    });
    if(res.ok){
      showCreateMsg('Created successfully — refresh to see it');
      el('create-form').reset();
      loadLinks();
    }else if(res.status === 409){
      showCreateMsg('That short name is already taken', true);
    }else{
      const text = await res.text().catch(()=>res.statusText);
      showCreateMsg(`Server error: ${res.status} ${text}`, true);
    }
  }catch(err){
    showCreateMsg(`Network error: ${err.message}`, true);
  }finally{btn.disabled = false}
}

function init(){
  el('create-form').addEventListener('submit', createLink);
  el('clear').addEventListener('click', ()=>{el('create-form').reset(); el('create-msg').textContent='';});
  el('refresh').addEventListener('click', loadLinks);
  loadLinks();
}

document.addEventListener('DOMContentLoaded', init);
