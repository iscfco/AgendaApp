// Espera a que el HTML termine de cargar
document.addEventListener("DOMContentLoaded", () => {
  const botonPedidos = document.getElementById("v-pills-pedidos-tab");

  // Escucha el evento de clic
  botonPedidos.addEventListener("click", () => {
    // Redirección programática
    window.location.href = "/?show=false";
  });
});