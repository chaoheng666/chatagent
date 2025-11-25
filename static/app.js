const sendBtn = document.getElementById('send');
const input = document.getElementById('input');
const messagesEl = document.getElementById('messages');

function appendMessage(text, clazz) {
  const wrap = document.createElement('div');
  wrap.className = 'msg ' + clazz;
  wrap.textContent = text;
  messagesEl.appendChild(wrap);
  messagesEl.scrollTop = messagesEl.scrollHeight;
}

async function send() {
  const msg = input.value.trim();
  if (!msg) return;
  appendMessage(msg, 'user');
  input.value = '';
  try {
    const res = await fetch('/api/chat', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ message: msg })
    });
    if (!res.ok) {
      const t = await res.text();
      appendMessage('错误: ' + t, 'bot');
      return;
    }
    const data = await res.json();
    appendMessage(data.reply || '(无回复)', 'bot');
  } catch (err) {
    appendMessage('请求失败: ' + err.message, 'bot');
  }
}

sendBtn.addEventListener('click', send);
input.addEventListener('keydown', (e) => { if (e.key === 'Enter' && (e.ctrlKey || e.metaKey)) send(); });
