
// Add event listener to the search button
document.addEventListener("DOMContentLoaded", () => {
    const form = document.getElementById('update-order-form');
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
        var id = document.getElementById('orderDetails-ID').value;
        var due_date = document.getElementById('orderDetails-DueDate').value;
        var delivery_date_formatted = due_date ? `${due_date}T00:00:00Z` : "";
        var orderData = {
            status: document.getElementById('orderDetails-Status').value,
            client_name: document.getElementById('orderDetails-ClientName').value,
            client_phone: document.getElementById('orderDetails-ClientPhone').value,
            client_address: document.getElementById('orderDetails-ClientAddress').value,
            total_price: parseFloat(document.getElementById('orderDetails-Total').value, 10),
            down_payment: parseFloat(document.getElementById('orderDetails-DownPayment').value, 10),
            delivery_date: delivery_date_formatted,
            description: document.getElementById('orderDetails-Description').value
        };

        update = confirm("¿Estás seguro de que deseas actualizar este pedido?");
        if (!update) {
            return;
        } 

        // Perform PUT /orders/:id
        try {
            const response = await fetch('/order/'+id, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(orderData)
            });

            if (response.ok) {
                const result = await response.json();
                window.location.href = document.referrer;
                alert('¡Pedido actualizado con éxito!');                
            } else {
                alert('Hubo un error en el servidor al actualizar el pedido.');
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