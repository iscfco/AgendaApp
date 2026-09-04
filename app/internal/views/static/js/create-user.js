
// Add event listener to the search button
document.addEventListener("DOMContentLoaded", () => {    
    const form = document.getElementById('create-user-form');
    form.addEventListener('submit', async event => {
        // Validate Form
        if (!form.checkValidity()) {
            event.preventDefault() 
            event.stopPropagation()
            form.classList.add('was-validated')
            return;
        }
        event.preventDefault()
        form.classList.add('was-validated')
        
        // Collect form data
        const userData = {
            user_full_name: document.getElementById('createUser-FullName').value,
            email:     document.getElementById('createUser-Email').value,
            password:  document.getElementById('createUser-Password').value,
            phone:     document.getElementById('createUser-Phone').value,
            role:      document.getElementById('createUser-Role').value,
            status:    document.getElementById('createUser-Status').value
        };

        // Perform Post /orders
        try {
            const response = await fetch('/user', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(userData)
            });

            if (response.ok) {
                const result = await response.json();
                alert('Usuario creado con éxito!');
                window.location.href = document.referrer;
            } else {
                const result = await response.json();
                console.error('Bad response:', result);
                alert(`Hubo un error en el servidor al crear el usuario: ${result.error}`);
            }

        } catch (error) {
            console.error('Error en la petición:', error);
            alert('No se pudo conectar con el servidor.');
        }
    
    }, false)
});

