document.getElementById("myForm").addEventListener('submit', function(event){
    event.preventDefault()

    const name = document.getElementById('name').value;
    const check = document.getElementById('check').checked;
    const data =  {name:name, check:check}
    post('/api/stat',data)
});

document.getElementById('submitButton').addEventListener('click', function(event){
    event.preventDefault()
    post('/api/config', {})
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
        console.log(response)
        if (!response.ok) {
            throw new Error('Response not Ok.')
        }
        return response.json();
    })
    .then(data => {
        document.getElementById('response').value = JSON.stringify(data, undefined, 2);
    })
    .catch(error => {
        console.error('Request failed.', error)
    });
}