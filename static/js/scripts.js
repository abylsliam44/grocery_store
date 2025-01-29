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
                    showAlert("Product added to cart!");
                } else {
                    showAlert("Failed to add to cart.", true);
                }
            } catch (error) {
                console.error("Error:", error);
                showAlert("An error occurred. Please try again.", true);
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

    // Приветственное сообщение в зависимости от времени суток
    const greetingElement = document.getElementById("greeting-message");
    const currentHour = new Date().getHours();
    let greetingText = "";

    if (currentHour >= 5 && currentHour < 12) {
        greetingText = "Good Morning";
    } else if (currentHour >= 12 && currentHour < 18) {
        greetingText = "Good Afternoon";
    } else {
        greetingText = "Good Evening";
    }

    if (greetingElement) {
        greetingElement.textContent = greetingText;
    }

    // Анимация карточек продуктов при загрузке
    const productCards = document.querySelectorAll(".product-card");
    productCards.forEach((card, index) => {
        setTimeout(() => {
            card.classList.add("animate__animated", "animate__fadeInUp");
        }, 200 * index); // Добавляем анимацию с задержкой
    });

    const form = document.getElementById("city-form");
    const cityInput = document.getElementById("city-input");
    const weatherElement = document.getElementById("weather-info");

    // Обработчик формы
    form.addEventListener("submit", (e) => {
        e.preventDefault(); // Предотвратить перезагрузку страницы при отправке формы
        const city = cityInput.value.trim(); // Получаем город из input

        if (city) {
            getWeather(city); // Запрос погоды
        }
    });

    // Получение погоды
    function getWeather(city) {
        const apiKey = "bc293dcd2a2e840cdaefee93433832a5"; // Замените на ваш ключ
        fetch(`https://api.openweathermap.org/data/2.5/weather?q=${city}&appid=${apiKey}&units=metric`)
            .then(response => {
                if (!response.ok) {
                    throw new Error("City not found.");
                }
                return response.json();
            })
            .then(data => {
                const temperature = data.main.temp;
                const description = data.weather[0].description;
                weatherElement.textContent = `Weather in ${city}: ${description}, ${temperature}°C`;
            })
            .catch(error => {
                console.error("Error fetching weather data:", error);
                weatherElement.textContent = `Error: ${error.message}`;
            });
    }
});
