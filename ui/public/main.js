document.getElementById('functionForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    const functionCode = document.getElementById('functionCode').value;
    const language = document.getElementById('language').value;
    const method = document.getElementById('method').value;
    const statusDiv = document.getElementById('status');
    const invokeForm = document.getElementById('invokeForm');

    statusDiv.style.display = 'block';
    statusDiv.innerHTML = 'Deploying function...';

    const response = await fetch('http://localhost:8080/submit', { // Update the URL here
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ code: functionCode, language: language, method: method }),
    });

    const result = await response.json();

    if (response.ok) {
        statusDiv.innerHTML = 'Function deployed successfully. You can now invoke it.';
        invokeForm.style.display = 'block';
    } else {
        statusDiv.innerHTML = 'Error deploying function: ' + result.message;
        invokeForm.style.display = 'none';
    }
});

document.getElementById('invokeForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    const response = await fetch('http://localhost:8080/invoke', { // Update the URL here
        method: 'POST',
    });

    const result = await response.text();
    alert('Function response: ' + result);
});