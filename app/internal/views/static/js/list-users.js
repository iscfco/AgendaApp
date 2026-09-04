
// Add event listener to the search button
document.addEventListener("DOMContentLoaded", () => {

    // Add event and perform the search
    const button = document.getElementById('list-users-search-button');
    button.addEventListener('click', () => {
        const username = document.getElementById('usernameVal').value;
        const role = document.getElementById('roleVal').value;
        const email = document.getElementById('emailVal').value;
        const status = document.getElementById('statusVal').value;
        const limit = document.getElementById('recordsPerPageVal').value;

        const filtros = {
            show: 'true',
            user_full_name: username,
            role: role,
            email: email,
            status: status,

            // Pagination fields
            page: '1' ,
            limit: limit
        };

        const parametros = new URLSearchParams(filtros).toString();
        window.location.href = `/user?${parametros}`;
    });

});

// Pagination
document.addEventListener("DOMContentLoaded", () => {
    const nav = document.getElementById("orders-pagination-nav");
    if (!nav) {
        return; // Caso cuando se abre la pagina por primera vez y no hay paginacion
    }

    // Leemos los atributos directamente del HTML real renderizado por Go
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
