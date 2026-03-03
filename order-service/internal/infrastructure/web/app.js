async function loadOrder() {
    const id = document.getElementById('orderId').value;
    if (!id) return;
    try {
        const res = await fetch('/order/' + id);
        const data = await res.json();
        if (!res.ok) {
            document.getElementById('result').textContent = `Ошибка: ${data.error || 'неизвестная ошибка'}`;
            return;
        }
        document.getElementById('result').textContent = JSON.stringify(data, null, 2);
    } catch (err) {
        document.getElementById('result').textContent = `Ошибка запроса: ${err.message}`;
    }
}