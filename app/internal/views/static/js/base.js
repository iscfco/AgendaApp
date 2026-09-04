// Espera a que el HTML termine de cargar
document.addEventListener("DOMContentLoaded", () => {
  const botonPedidos = document.getElementById("v-menu-pedidos-tab");
  botonPedidos.addEventListener("click", () => {
    window.location.href = "/?show=false";
  });

  const botonUsers = document.getElementById("v-menu-usuarios-tab");
  if (botonUsers) {
    botonUsers.addEventListener("click", () => {
      window.location.href = "/user?show=false";
    });
  }

  const botonLogout = document.getElementById("v-menu-logout");
  botonLogout.addEventListener("click", () => {
    window.location.href = "/logout";
  });
});