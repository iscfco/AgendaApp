
// Add event listener to the search button
document.addEventListener("DOMContentLoaded", () => {

    // Add event and perform the search
    const boton = document.getElementById('list-orders-search-button');
    boton.addEventListener('click', () => {
        const authorName = document.getElementById('userAuthorNameVal').value;
        const keyword = document.getElementById('keywordVal').value;
        const client_name = document.getElementById('clientNameVal').value;
        const from_date = document.getElementById('fromDateVal').value;
        const to_date = document.getElementById('toDateVal').value;
        const status = getStatusValue(document.getElementById('orderStatusButton').textContent);
        const limit = document.getElementById('recordsPerPageVal').value;

        const filtros = {
            show: 'true',
            user_creator_name: authorName,
            keyword: keyword,
            client_name: client_name,
            from_date: from_date,
            to_date: to_date,
            status: status,

            // Pagination fields
            page: '1' ,
            limit: limit
        };

        const parametros = new URLSearchParams(filtros).toString();
        window.location.href = `/?${parametros}`;
    });

    // Add event listener to allow select the Status filter
    document.querySelectorAll('.order-status-value').forEach(item => {
        item.addEventListener('click', function () {
        const orderStatusButton = document.getElementById('orderStatusButton');
        orderStatusButton.textContent = this.textContent;
        });
    });

});


// Status Enum
const StatusEnum = Object.freeze({
  pendiente: "pending",
  entregado: "delivered",
  todos: "all" // Represents all statuses
});

function getStatusValue(key) {
  const normalizedKey = key.trim().toLowerCase(); 
  return StatusEnum[normalizedKey] || "";
}

// Pagination
document.addEventListener("DOMContentLoaded", () => {
    const nav = document.getElementById("orders-pagination-nav");
    // Leemos los atributos directamente del HTML real renderizado por Go
    const currentPage = parseInt(nav.dataset.current) || 1;
    const totalPages = parseInt(nav.dataset.total) || 1;


    // Buscamos todos los enlaces de la paginación
    const pageLinks = document.querySelectorAll(".orders-pagination-link");

    pageLinks.forEach(link => {
        link.addEventListener("click", (e) => {
            e.preventDefault(); // Evitamos que el '#' recargue o mueva la pantalla

            // Validar si el botón padre está deshabilitado (Bootstrap 'disabled')
            if (link.parentElement.classList.contains("disabled")) {
                return;
            }

            // 3. Leer el atributo data-page que definimos en el HTML
            let targetPage = link.dataset.page;

            // 4. Resolver si es un número, "first" o "last"
            if (targetPage === "FIRST") {
                targetPage = 1;
            } else if (targetPage === "LAST") {
                targetPage = totalPages;
            } else {
                targetPage = parseInt(targetPage);
            }

            // 5. Modificar la URL actual respetando los demás parámetros (como el limit o filtros)
            const url = new URL(window.location.href);
            url.searchParams.set("page", targetPage);

            // 6. Redireccionar para aplicar el cambio
            window.location.href = url.toString();
        });
    });
});



     
