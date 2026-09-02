
// Add event listener to the search button
document.addEventListener("DOMContentLoaded", () => {

    // Add event ato create order
    const boton = document.getElementById('create-orders-button');
    boton.addEventListener('clickX', () => {
        const client_name = document.getElementById('createOrder-ClientName').value;
        const client_address = document.getElementById('createOrder-ClientAddress').value; 
        const due_date = document.getElementById('createOrder-DueDate').value;
        const client_phone = document.getElementById('createOrder-ClientPhone').value;
        const total = document.getElementById('createOrder-Total').value;
        const down_payment = document.getElementById('createOrder-DownPayment').value;
        const description = document.getElementById('createOrder-Description').value;


        const filtros = {
            client_name: client_name,
            client_phone: client_phone,
            client_address: client_address,
            total_price: total,
            down_payment: down_payment,
            delivery_date: due_date,
            description: description
        };

        const parametros = new URLSearchParams(filtros).toString();
        window.location.href = `/?${parametros}`;
    });
});

// Add event listener to the search button
document.addEventListener("DOMContentLoaded", () => {

    // Add event ato create order
    const boton = document.getElementById('create-orders-button');
    boton.addEventListener('clickX', () => {

    });


    
    const form = document.getElementById('create-orders-form');
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
        const client_name = document.getElementById('createOrder-ClientName').value;
        const client_address = document.getElementById('createOrder-ClientAddress').value; 
        const due_date = document.getElementById('createOrder-DueDate').value;
        const client_phone = document.getElementById('createOrder-ClientPhone').value;
        const total = document.getElementById('createOrder-Total').value;
        const down_payment = document.getElementById('createOrder-DownPayment').value;
        const description = document.getElementById('createOrder-Description').value;

        const delivery_date_formatted = due_date ? `${due_date}T00:00:00Z` : "";
        const orderData = {
            client_name: client_name,
            client_phone: client_phone,
            client_address: client_address,
            total_price: parseFloat(total, 10),
            down_payment: parseFloat(down_payment, 10),
            delivery_date: delivery_date_formatted,
            description: description
        };

        // Perform Post /orders
        try {
            const response = await fetch('/order', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(orderData)
            });

            if (response.ok) {
                const result = await response.json();
                window.location.href = '/';
                alert('¡Pedido creado con éxito!');
                // form.reset(); // Limpia los campos del formulario
                // form.classList.remove('was-validated'); // Quita los colores verdes de validación
                
                // // Aquí puedes redireccionar al usuario si quieres:
                // const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));
                // await sleep(5000); 
                // window.location.href = '/';
                
            } else {
                alert('Hubo un error en el servidor al crear el pedido.');
            }

        } catch (error) {
            console.error('Error en la petición:', error);
            alert('No se pudo conectar con el servidor.');
        }
    
    }, false)
});

