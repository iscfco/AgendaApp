
// Add event listener to the search button
document.addEventListener("DOMContentLoaded", () => {
    const form = document.getElementById('update-user-form');
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
        var id = document.getElementById('userDetails-ID').value;
        var userData = {
            user_full_name: document.getElementById('userDetails-UserFullName').value,
            email: document.getElementById('userDetails-Email').value,
            phone: document.getElementById('userDetails-Phone').value,
            password: document.getElementById('userDetails-Password').value,
            role: document.getElementById('userDetails-Role').value,
            status: document.getElementById('userDetails-Status').value,
        };

        update = confirm("¿Estás seguro de que deseas actualizar este usuario?");
        if (!update) {
            return;
        } 

        // Perform PUT /orders/:id
        try {
            const response = await fetch('/user/'+id, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(userData)
            });

            if (response.ok) {
                const result = await response.json();
                window.location.href = document.referrer;
                alert('¡Pedido actualizado con éxito!');                
            } else {
                const result = await response.json();
                console.error('Bad response:', result);
                alert(`Hubo un error en el servidor al actualizar el pedido: ${result.error}`);
            }

        } catch (error) {
            console.error('Error en la petición:', error);
            alert('Error interno');
        }
    
    }, false)
});

function removeOrder(boton) {
    const orderId = boton.getAttribute('data-id');
    
    if (!confirm('¿Estás seguro de que deseas eliminar esta orden?')) {
        return;
    }

    fetch(`/order/${orderId}`, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => {
        if (response.ok) {
            alert('Orden eliminada con éxito');
            window.location.href = document.referrer;
        } else {
            alert('Error al eliminar la orden');
        }
    })
    .catch(error => console.error('Error:', error));
}