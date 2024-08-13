document.getElementById("myForm").addEventListener('submit', function(event){
    event.preventDefault()

    const name = document.getElementById('name').value;
    const check = document.getElementById('check').checked;
    const data =  {name:name, check:check}
    post('/api/echo',data)
});

document.getElementById('apiCommands').addEventListener('click', function(event){
    event.preventDefault()
    post('/api/commands', {})
})

document.getElementById('apiConfigs').addEventListener('click', function(event){
    event.preventDefault()
    post('/api/configs', {})
})

function post(url, data) {
    fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Response not Ok.')
        }
        return response.json();
    })
    .then(data => {
        console.log(data)
        document.getElementById('response').value = JSON.stringify(data, undefined, 2);
    })
    .catch(error => {
        console.error('Request failed.', error)
    });
}