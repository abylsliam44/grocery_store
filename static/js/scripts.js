document.addEventListener("DOMContentLoaded", () => {
    // Подсвечиваем текущую страницу в навигации
    const links = document.querySelectorAll("nav ul li a");
    links.forEach(link => {
        if (link.href === window.location.href) {
            link.style.fontWeight = "bold";
            link.style.textDecoration = "underline";
        }
    });

    // Обработка добавления в корзину
    const addToCartForms = document.querySelectorAll("form[action='/cart']");
    addToCartForms.forEach(form => {
        form.addEventListener("submit", async (e) => {
            e.preventDefault(); // Убедитесь, что это уместно
            const formData = new FormData(form);

            try {
                const response = await fetch("/cart", {
                    method: "POST",
                    body: formData,
                });

                if (response.ok) {
                    alert("Product added to cart!");
                } else {
                    alert("Failed to add to cart.");
                }
            } catch (error) {
                console.error("Error:", error);
            }
        });
    });

    // Анимация перехода между страницами
    links.forEach(link => {
        link.addEventListener("click", (e) => {
            e.preventDefault();
            document.body.style.opacity = 0;
            setTimeout(() => {
                window.location.href = link.href;
            }, 300);
        });
    });

    // Уведомления
    function showAlert(message, isError = false) {
        const alert = document.createElement("div");
        alert.className = "alert" + (isError ? " error" : "");
        alert.textContent = message;
        document.body.appendChild(alert);

        setTimeout(() => {
            alert.classList.add("show");
        }, 10);

        setTimeout(() => {
            alert.classList.remove("show");
            alert.addEventListener("transitionend", () => alert.remove());
        }, 3000);
    }
});


function showAlert(message, isError = false) {
    const alert = document.createElement("div");
    alert.className = "alert" + (isError ? " error" : "");
    alert.textContent = message;
    document.body.appendChild(alert);

    setTimeout(() => {
        alert.classList.add("show");
    }, 10);

    setTimeout(() => {
        alert.classList.remove("show");
        alert.addEventListener("transitionend", () => alert.remove());
    }, 3000);
}

